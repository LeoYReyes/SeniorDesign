Stolen Object Tracker Server Version 1.0 04/22/2014

---------------------------------------------------------------------------------
Requirements:
     *MySQL
     *Go compiler (GoLang)

---------------------------------------------------------------------------------

---------------------------------------------------------------------------------
Install Instructions:

     Database Setup:
	
	Login Credentials:
		User = root
		Password = toor
		* User and Password can be changed, but make sure to also update
		  the credentials used in the database.go source file.

	1. Create a new MySQL database named "trackerdb" (database name can be
	changed, but must also be updated in the database.go source code).
	2. Import SQL dump (if you know how to do it on your own skip this step)
	    Command Line:
		* Change to directory containing trackerdb.sql
		* mysql --user=name database_name < trackerdb.sql
	3. Database setup complete

     Server Build:
	1. Change to the /gosotserver directory
	2. run "go build server"
	3. Server executable will compile to the current directory you are in.

---------------------------------------------------------------------------------

---------------------------------------------------------------------------------
Server Setup Instructions:

	NOTE: The server executable can be in any directory, but the directories
	      in step 2 must also be in the same directory as the executable.

	1. Make sure that the database is running before running the server
	2. Copy the following directories to the same directory that the
	   server executable is located.
		* /css
		* /images
		* /js
		* /templates
	3. Run the SMS gateway application and connect to the server's IP 
	   address using port 10016.
	4. Once the database and the SMS gateway are running, you can now execute
	   the server executable.
	
---------------------------------------------------------------------------------
Accessing Website:

	Once the server is running, the website can be accessed through any
	browser with the URL: 
		
		hostname:8080

---------------------------------------------------------------------------------

---------------------------------------------------------------------------------
Known bugs/issues:

	* Device registration request:
		- The server is currently checking the device ID using only length
		  to determine what type of device the request was sent as. Anything 
		  less than 12 characters is recognized as a GPS device and anything
		  greater is a laptop device. Any IDs that aren’t of length 10(GPS)
		  or 12(laptop) could cause a problem.

---------------------------------------------------------------------------------