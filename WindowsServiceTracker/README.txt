Windows Service installation
open command prompt in administrator mode and navigate to the directory containing WindowsServiceTracker.exe
use the command 'WindowsServiceTracker.exe --install'
Open the 'services' application in, select 'WindowsServiceTracker' and press 'start'
The service will also start automatically on windows startup

After the service is run for the first time, it will generate an ID.txt file. Use that ID to register your device on the server.

Additionally, in the folder that contains WindowsServiceTracker.exe, there is a file, WindowsServiceTracker.exe.xml. There are 3 settings in this file. one is the port number, it is by default 10015 and will most likely not need to be changed. The others are an IP address and Domain name. By default the IP address points to your computer and the domain name is invalid. Set either of these values to the server you wish to connect to. The domain name takes precedence over the IP address when connecting and is the preferred method. Any changes to this file will require a restart of the service to take effect.




Key-logger installation
Navigate to the directory (in Windows explorer) containing 'WTKL.exe'
Right click 'WTKL.exe' and select create shortcut
place the shortcut in 'C:\Documents and Settings\All Users\Start Menu\Programs\Startup'
