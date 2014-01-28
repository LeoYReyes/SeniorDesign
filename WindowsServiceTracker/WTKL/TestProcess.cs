using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Diagnostics;
using System.Threading;

namespace WindowsServiceTracker
{
    class TestProcess
    {
        /***********************************************************************************
         * This is a separate program from the tracker that will be called to run by Tracker.cs.
         * Currently it doesn't get built by Visual Studio (don't know if that's possible),
         * but to build use Command Prompt and type:
         * csc /out:TestProcess.exe TestProcess.cs
         * 
         * The TestProcess.exe file must then be placed in the same directory that you install
         * the service from. Example directory:
         * /????/SeniorDesign/WindowsServiceTracker/WindowsTrackerService/bin/Debug/
         ***********************************************************************************/
        private const string errorLogName = "TrackerErrorLog";
        private const string errorLogMachine = "TrackerComputer";
        private const string errorLogSource = "WindowsServiceTrackerProcess";

        public TestProcess()
        {
            if (!EventLog.SourceExists(errorLogSource))
            {
                EventLog.CreateEventSource(errorLogSource, errorLogName);
            }
        }

        public static void Main()
        {
            EventLog.WriteEntry(errorLogSource, "Test event 1", EventLogEntryType.Error);
            Thread.Sleep(5000);
            EventLog.WriteEntry(errorLogSource, "Test event 2", EventLogEntryType.Error);
            Thread.Sleep(5000);
            EventLog.WriteEntry(errorLogSource, "Test event 3", EventLogEntryType.Error);
            Thread.Sleep(5000);
            EventLog.WriteEntry(errorLogSource, "Test event 4", EventLogEntryType.Error);
            Thread.Sleep(5000);
            EventLog.WriteEntry(errorLogSource, "Test event 5", EventLogEntryType.Error);
            Thread.Sleep(5000);
            EventLog.WriteEntry(errorLogSource, "Test event 6", EventLogEntryType.Error);
            Thread.Sleep(5000);
        }
    }
}
