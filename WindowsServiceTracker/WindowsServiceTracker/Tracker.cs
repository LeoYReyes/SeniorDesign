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

namespace WindowsServiceTracker
{
    /***********************************************************************
     * This class is where the majority of the work of the service is done.
     * We may possibly want to split up our different functions into different
     * classes.
     ***********************************************************************/
    public partial class Tracker : ServiceBase
    {
        //Constants
        //127.0.0.1 = 0x0100007F because of network byte order
        private const long IP_ADDRESS = 0x0100007F; //127.0.0.1 as placeholder
        private const int PORT = 10000;
        private const string errorLogName = "TrackerErrorLog";
        private const string errorLogMachine = "TrackerComputer";
        private const string errorLogSource = "WindowsServiceTracker";

        private IPEndPoint IPPort = new IPEndPoint(IP_ADDRESS, PORT);
        private TcpClient tcp;

        public Tracker()
        {
            InitializeComponent();
            //Creates the error log source if it doesn't already exist
            if(!EventLog.SourceExists(errorLogSource))
            {
                EventLog.CreateEventSource(errorLogSource, errorLogName);
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

            //Write to the Windows Event Logs, shows up under Windows Logs --> Application
            EventLog.WriteEntry(errorLogSource, "Test event", EventLogEntryType.Information);

            //Some test tcp connection stuff
            try
            {
                tcp = new TcpClient();
                tcp.Connect(IPPort);
            }
            finally
            {
                //Another write to the Windows Event Logs, shows up in the same place as before but as type "Error" instead of type "Information"
                string error = "Service failed to connect to IP address " + IP_ADDRESS + " on port " + PORT;
                EventLog.WriteEntry(errorLogSource, error, EventLogEntryType.Error);
            }
            NetworkStream tcpStream = tcp.GetStream();
            if (tcpStream.CanWrite)
            {
                Byte[] machineName = Encoding.UTF8.GetBytes(Environment.MachineName + "\n");
                tcpStream.Write(machineName, 0, machineName.Length);
            }

            tcp.Close();
            tcpStream.Close();
        }

        protected override void OnStop()
        {
        }
    }
}
