using System;
using System.Collections.Generic;
using System.Linq;
using System.ServiceProcess;
using System.Text;
using System.Threading.Tasks;
using System.Configuration.Install; //used for ManagedInstallerClass
using System.Reflection; //used for Assembly class

namespace WindowsServiceTracker
{
    static class Program
    {
        /// <summary>
        /// The main entry point for the application.
        /// </summary>
        static void Main(string[] args)
        {
            /**********************************************
             * Since we have to install the service every time we test it I found
             * this code online to help with that. Run the WindowsServiceTracker.exe
             * from command line and either use the argument "--install" to install
             * the service or use "--uninstall" to uninstall the service. Don't use
             * any arguments if you just want to run the service normally once it's
             * been installed.
             *********************************************/
            if (Environment.UserInteractive)
            {
                string parameter = string.Concat(args);
                switch (parameter)
                {
                    case "--install":
                        ManagedInstallerClass.InstallHelper(new[] { Assembly.GetExecutingAssembly().Location });
                        break;
                    case "--uninstall":
                        ManagedInstallerClass.InstallHelper(new[] { "/u", Assembly.GetExecutingAssembly().Location });
                        break;
                }
            }
            else
            {
                ServiceBase[] ServicesToRun;
                ServicesToRun = new ServiceBase[] 
                { 
                    new Tracker() 
                };
                ServiceBase.Run(ServicesToRun);
            }
        }
    }
}
