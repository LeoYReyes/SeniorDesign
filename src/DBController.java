/*
*	Steven Whaley - created: January 23, 2014  - last updated: February 5, 2014
*
*
*	OVERVIEW:
*			
*			This is the Database Controller. It provides functionality for interacting with the database.
*			Getters are provided as well as the more general submitQuery and submitUpdate. See
*			the comments for each method for details.
*
*
*	note: 	consider revising so connection uses DataSource object instead of DriverManager
*	note: 	mysql-connector-java-5.1.28-bin.jar placed in referenced libraries 
*   note: 	on mac -  /usr/local/mysql
*
*   useful links:
*   		http://docs.oracle.com/javase/7/docs/api/java/sql/package-summary.html
*   		http://docs.oracle.com/javase/7/docs/api/javax/sql/package-summary.html
*   		http://www.developer.com/java/data/article.php/3417381/Using-JDBC-with-MySQL-Getting-Started.htm
*   		http://dev.mysql.com/doc/refman/5.0/en/tutorial.html
*   		http://www.vogella.com/tutorials/MySQLJava/article.html
*   		http://www.avajava.com/tutorials/lessons/how-do-i-use-jdbc-to-query-a-mysql-database.html
*   
*   changes:
*   		Steven Whaley:
*   		January 23 - Started researching jdbc and mysql, and wrote some database connection code. Installed
*   					 connector/j, mysql and sequel pro on my mac. Created database trackerdb.
*   
*   		Steven Whaley:		
*   		January 25 - Successfully connected to database, but code is all in main and isn't ready to
*   					 interface with server code.
*   		
*   		Steven Whaley:			
*   		January 26 - created separate submit_query() and submit_update() and started testing with a driver
*   					 program to create objects and call methods.
*   
*   		Steven Whaley:
*   		January 27 - submit_query() now successfully queries the database and provides an options input
*   					 for which table you want to query. submit_update() currently is the same code as
*   					 submit_query() and needs to use executeUpdate(). I also created get_device_info() which
*   					 gets info about a device and returns it in a 2D array.
*   
*   		Steven Whaley:
*   		January 30 - added get_account_info(), get_customer_info(), console_print(), openConnection(),
*   					 closeConnection(). Numerous database changes and code changes. id_in inputs are Strings
*   					 instead of ints now.
*   
*   		Steven Whaley:
*   		February 1 - updated to match database changes, various small changes, commented out submitQuery
*   					 and submitUpdate temporarily
*   
*   		Steven Whaley:
*   		February 3 - changed to camel case
*   
*   		Steven Whaleu:
*   		Feb. 4 & 5 - redesigned the structure some. strings for database field names are now instantiated when the
*   					 constructor is called. now using prepared statement for security against sql injections. The
*   					 query is passed when the prepared statement is created and the query is pre-compiled.
*   
*   TO DO: code updates whenever database structure changes, update submit_update and submit_query, test recent changes
*   		with main.
*  
*   	   -also what other functionality is needed?
*/

import java.sql.*;

//import javax.sql.*;

public class DBController
{
	private String account1;
	private String account2;
	private String account3;
	private String account4;
	private String account5;
			
	private String customer1;
	private String customer2;
	private String customer3;
	private String customer4;
	private String customer5;
	private String customer6;
	
	private String gpsDevice1;
	private String gpsDevice2;
	private String gpsDevice3;
	private String gpsDevice4;
	private String gpsDevice5;
	
	private String ipAddress1;
	private String ipAddress2;
	private String ipAddress3;
	
	private String ipList1;
	private String ipList2;
	private String ipList3;
	
	private String keyLogs1;
	private String keyLogs2;
	private String keyLogs3;
	private String keyLogs4;
	
	private String laptopDevice1;
	private String laptopDevice2;
	private String laptopDevice3;
	private String laptopDevice4;
	
	public DBController()
	{
		account1 = "id";
		account2 = "customerId";
		account3 = "userName";
		account4 = "password";
		account5 = "admin";
				
		customer1 = "id";
		customer2 = "phoneNumber";
		customer3 = "address";
		customer4 = "email";
		customer5 = "firstName";
		customer6 = "lastName";
		
		gpsDevice1 = "id";
		gpsDevice2 = "name";
		gpsDevice3 = "customerId";
		gpsDevice4 = "latitude";
		gpsDevice5 = "longitude";
		
		ipAddress1 = "id";
		ipAddress2 = "listId";
		ipAddress3 = "ipAddress";
		
		ipList1 = "id";
		ipList2 = "deviceId";
		ipList3 = "timestamp";
		
		keyLogs1 = "id";
		keyLogs2 = "deviceId";
		keyLogs3 = "timestamp";
		keyLogs4 = "data";
		
		laptopDevice1 = "id";
		laptopDevice2 = "deviceName";
		laptopDevice3 = "customerId";
		laptopDevice4 = "macAddress";
	}
	
   /*
    * 		Used to form initial connection with database. Returns a connection object
    * 		that can be used for creating statements for sql queries. Initializes jdbc
    * 		driver. Forms Connection with DriverManager.
    * 
    * 		Steven Whaley - created January 30, 2014 - last updated January 30, 2014
    */
	private Connection openConnection()
	{
		String url = "jdbc:mysql://localhost/trackerdb";

	  	try 
	    {
	    	//initialize JDBC driver
	    	Class.forName("com.mysql.jdbc.Driver");
	    }
	    catch(Exception e) 
	    {
	     	System.out.println("Unable to load driver.");
	    	e.printStackTrace();
	    }
	    
	    try 
	    {
			//parameters for url, username and password - defaulted to root and null
	    	
		    Connection con = DriverManager.getConnection(url, "root", null);
		      
		    con.setAutoCommit(false);
		    con.commit();
		    return con;
	    }
	    catch(Exception e2)
	    {
	    	e2.printStackTrace();
	    }	
	    return null;
	}
	
   /*
    * 	Used to form close connection with database. Takes the Connection object to 
    * 	be closed as input.
    * 
    * 		Steven Whaley - created January 30, 2014 - last updated January 30, 2014
    */
	public void closeConnection(Connection c)
	{
		try
		{
			Connection con = c;
			con.close();
		}
		catch(SQLException e)
		{
			e.printStackTrace();
		}	
	}
	
   /*
	*	getLaptopDeviceInfo takes a String for the device id and returns a 2D array containing
	*	field names on the top row and the values associated with those fields in the second row.
	*
	*		Steven Whaley - created January 30 - revised from getDeviceInfo()
	*		Steven Whaley - updated February 1 - database changes
	*		Steven Whaley - updated February 5 - prepared statement
	*/
	public String[][] getLaptopDeviceInfo(String id_in)
	{
		try
		{
			Connection con = openConnection();
			
		    PreparedStatement stmt = con.prepareStatement("SELECT * FROM laptopDevice WHERE id = \'" + id_in +"\'");
		    
		    ResultSet result = stmt.executeQuery();
	      
		    int columnLength = 2;
		    int rowLength = 4;
		    
		    String[][] ldeviceInfo = new String[columnLength][rowLength];
		     
		    
		    String field1 = "";
		    String field2 = "";
		    String field3 = "";
		    String field4 = "";
		    
		    ldeviceInfo[0][0] = laptopDevice1;
		    ldeviceInfo[0][1] = laptopDevice2;
		    ldeviceInfo[0][2] = laptopDevice3;
		    ldeviceInfo[0][3] = laptopDevice4;
		    
		    if (result.next()) 
		    {   	
				field1 = Integer.toString(result.getInt(laptopDevice1));
				field2 = result.getString(laptopDevice2);
				field3 = Integer.toString(result.getInt(laptopDevice3));
				field4 = result.getString(laptopDevice4);
					 	
				ldeviceInfo[1][0] = field1;
				ldeviceInfo[1][1] = field2;
				ldeviceInfo[1][2] = field3;
				ldeviceInfo[1][3] = field4;
		    }
		    
		    stmt.close();
		    closeConnection(con);
		    result.close();
		    
		    return ldeviceInfo;
		}
		catch (Exception e) 
		{
			e.printStackTrace();
		}
		return null;	
	}
	
   /*
	*	getGpsDeviceInfo takes a String for the device id and returns a 2D array containing
	*	field names on the top row and the values associated with those fields in the second row.
	*
	*		Steven Whaley - updated January 30 - revised from getDeviceInfo()
	*		Steven Whaley - updated February 5 - prepared statement
	*/
	public String[][] getGpsDeviceInfo(String idIn)
	{
		try
		{
		    Connection con = openConnection();
		    PreparedStatement stmt = con.prepareStatement("SELECT * FROM gpsDevice WHERE id = \'" + idIn + "\'");
		       
		    ResultSet result = stmt.executeQuery();
	      
		    int columnLength = 2;
		    int rowLength = 5;
		    
		    String[][] gdeviceInfo = new String[columnLength][rowLength];
		     
		    String field1 = "";
		    String field2 = "";
		    String field3 = "";
		    String field4 = "";
		    String field5 = "";
		    
		    gdeviceInfo[0][0] = gpsDevice1;
		    gdeviceInfo[0][1] = gpsDevice2;
		    gdeviceInfo[0][2] = gpsDevice3;
		    gdeviceInfo[0][3] = gpsDevice4;
		    gdeviceInfo[0][4] = gpsDevice5;
		    
		    if (result.next()) 
		    {  
				field1 = Integer.toString(result.getInt(gpsDevice1));
				field2 = result.getString(gpsDevice2);
				field3 = Integer.toString(result.getInt(gpsDevice3));
				field4 = result.getString(gpsDevice4);
				field5 = result.getString(gpsDevice5);
					 	
				gdeviceInfo[1][0] = field1;
				gdeviceInfo[1][1] = field2;
				gdeviceInfo[1][2] = field3;
				gdeviceInfo[1][3] = field4;
				gdeviceInfo[1][4] = field5;
		    }
	  
		    stmt.close();
		    closeConnection(con);
		    result.close();
		    
		    return gdeviceInfo;
	    }
	    catch (Exception e) 
	    {
	    	e.printStackTrace();
	    }
	    return null;
	}
	
   /*
	*	getAccountInfo takes a String for the account id and returns a 2D array containing
	*	field names on the top row and the values associated with those fields in the second row.
	*
	*		Steven Whaley - created January 30 - revised from getAccountInfo()
	*		Steven Whaley - updated February 1 - database changes
	*		Steven Whaley - updated February 5 - prepared statement
	*/
	public String[][] getAccountInfo(String idIn)
	{
		try
		{
			Connection con = openConnection();
			
		    PreparedStatement stmt = con.prepareStatement("SELECT * FROM account WHERE id = \'" + idIn + "\'");
		    
		    ResultSet result = stmt.executeQuery();
	      
		    int columnLength = 2;
		    int rowLength = 5;
		    
		    String[][] accountInfo = new String[columnLength][rowLength];
		     
		    String field1 = "";
		    String field2 = "";
		    String field3 = "";
		    String field4 = "";
		    String field5 = "";
		    
		    accountInfo[0][0] = account1;
		    accountInfo[0][1] = account2;
		    accountInfo[0][2] = account3;
		    accountInfo[0][3] = account4;
		    accountInfo[0][4] = account5;
		    
		    
		    if (result.next()) 
		    {   
				field1 = Integer.toString(result.getInt(account1));
				field2 = Integer.toString(result.getInt(account2));
				field3 = result.getString(account3);
				field4 = result.getString(account4);
				field5 = Short.toString(result.getShort(account5));
						
				accountInfo[1][0] = field1;
				accountInfo[1][1] = field2;
				accountInfo[1][2] = field3;
				accountInfo[1][3] = field4;
				accountInfo[1][4] = field5;
		    }
		    		    
		    stmt.close();
		    con.close();
		    result.close();
		    
		    return accountInfo;
	    }
	    catch (Exception e) 
	    {
	    	e.printStackTrace();
	    }
	    return null;
	}
	
   /*
	*	getCustomerInfo takes a String for the customer id and returns a 2D array containing
	*	field names on the top row and the values associated with those fields in the second row.
	*
	*		Steven Whaley - updated January 30 - revised from getCustomerInfo()
	*		Steven Whaley - updated February 1 - database changes
	*		Steven Whaley - updated February 5 - prepared statement
	*/
	public String[][] getCustomerInfo(String idIn)
	{
		 try
		 {
			Connection con = openConnection();
		    PreparedStatement stmt = con.prepareStatement("SELECT * FROM customer WHERE id = \'" + idIn + "\'");
		    	
		    ResultSet result = stmt.executeQuery();
	      
		    int columnLength = 2;
		    int rowLength = 6;
		    
		    String[][] customerInfo = new String[columnLength][rowLength];
		     
		    String field1 = "";
		    String field2 = "";
		    String field3 = "";
		    String field4 = "";
		    String field5 = "";
		    String field6 = "";
		    
		    customerInfo[0][0] = customer1;
		    customerInfo[0][1] = customer2;
		    customerInfo[0][2] = customer3;
		    customerInfo[0][3] = customer4;
		    customerInfo[0][4] = customer5;
		    customerInfo[0][5] = customer6;
		    
		    if (result.next()) 
		    {  
				field1 = Integer.toString(result.getInt(customer1));
				field2 = result.getString(customer2);
				field3 = result.getString(customer3);
				field4 = result.getString(customer4);
				field5 = result.getString(customer5);
				field6 = result.getString(customer6);
					 	
				customerInfo[1][0] = field1;
				customerInfo[1][1] = field2;
				customerInfo[1][2] = field3;
				customerInfo[1][3] = field4;
				customerInfo[1][4] = field5;
				customerInfo[1][5] = field6;
		    }
		    		    
		    stmt.close();
		    closeConnection(con);
		    result.close();
		    
		    return customerInfo;
	    }
	    catch (Exception e) 
	    {
	    	e.printStackTrace();
	    }
	    return null;
	}

	 /*
		*	getIpAddressInfo takes a String for the  list id and returns a 2D array containing
		*	field names on the top row and the values associated with those fields in the second row.
		*
		*		Steven Whaley - created February 4
		*		Steven Whaley - updated February 5 - prepared statement
		*/
		public String[][] getIpAddressInfo(String id_in)
		{
			try
			{
				Connection con = openConnection();
				
			    PreparedStatement stmt = con.prepareStatement("SELECT * FROM ipAddress WHERE id = \'" + id_in +"\'");
			    
			    ResultSet result = stmt.executeQuery();
		      
			    int columnLength = 2;
			    int rowLength = 3;
			    
			    String[][] ipInfo = new String[columnLength][rowLength];
			     
			    
			    String field1 = "";
			    String field2 = "";
			    String field3 = "";
			  
			    
			    ipInfo[0][0] = ipAddress1;
			    ipInfo[0][1] = ipAddress2;
			    ipInfo[0][2] = ipAddress3;
			  
			    
			    if (result.next()) 
			    {   	
					field1 = Integer.toString(result.getInt(ipAddress1));
					field2 = Integer.toString(result.getInt(ipAddress2));
					field3 = result.getString(ipAddress3);
						 	
					ipInfo[1][0] = field1;
					ipInfo[1][1] = field2;
					ipInfo[1][2] = field3;
			    }
			    
			    stmt.close();
			    closeConnection(con);
			    result.close();
			    
			    return ipInfo;
			}
			catch (Exception e) 
			{
				e.printStackTrace();
			}
			return null;	
		}
		
	    /*
		*	getIpAddressInfo takes a String for the  list id and returns a 2D array containing
		*	field names on the top row and the values associated with those fields in the second row.
		*
		*		Steven Whaley - created February 4
		*		Steven Whaley - updated February 5 - prepared statement
		*/
		public String[][] getIpListInfo(String id_in)
		{
			try
			{
				Connection con = openConnection();
				
			    PreparedStatement stmt = con.prepareStatement("SELECT * FROM ipList WHERE id = \'" + id_in +"\'");
			       
			    ResultSet result = stmt.executeQuery();
		      
			    int columnLength = 2;
			    int rowLength = 3;
			    
			    String[][] ipInfo = new String[columnLength][rowLength];
			     
			    
			    String field1 = "";
			    String field2 = "";
			    String field3 = "";
			  
			    
			    ipInfo[0][0] = ipList1;
			    ipInfo[0][1] = ipList2;
			    ipInfo[0][2] = ipList3;
			  
			    
			    if (result.next()) 
			    {   	
					field1 = Integer.toString(result.getInt(ipList1));
					field2 = Integer.toString(result.getInt(ipList2));
					field3 = result.getString(ipList3);
						 	
					ipInfo[1][0] = field1;
					ipInfo[1][1] = field2;
					ipInfo[1][2] = field3;
			    }
			    
			    stmt.close();
			    closeConnection(con);
			    result.close();
			    
			    return ipInfo;
			}
			catch (Exception e) 
			{
				e.printStackTrace();
			}
			return null;	
		}

		 /*
		*	getKeyLogs takes a String for the id and returns a 2D array containing
		*	field names on the top row and the values associated with those fields in the second row.
		*
		*		Steven Whaley - created February 4
		*		Steven Whaley - updated February 5 - prepared statement
		*/
		public String[][] getKeyLogsInfo(String id_in)
		{
			try
			{
				Connection con = openConnection();
				
			    PreparedStatement stmt = con.prepareStatement("SELECT * FROM keyLogs WHERE id = \'" + id_in +"\'");
			    
			    ResultSet result = stmt.executeQuery();
		      
			    int columnLength = 2;
			    int rowLength = 4;
			    
			    String[][] klInfo = new String[columnLength][rowLength];
			     
			    
			    String field1 = "";
			    String field2 = "";
			    String field3 = "";
			    String field4 = "";
			  
			    
			    klInfo[0][0] = keyLogs1;
			    klInfo[0][1] = keyLogs2;
			    klInfo[0][2] = keyLogs3;
			    klInfo[0][3] = keyLogs4;
			  
			    
			    if (result.next()) 
			    {   	
					field1 = Integer.toString(result.getInt(keyLogs1));
					field2 = Integer.toString(result.getInt(keyLogs2));
					field3 = result.getString(keyLogs3);
					field4 = result.getString(keyLogs4);
						 	
					klInfo[1][0] = field1;
					klInfo[1][1] = field2;
					klInfo[1][2] = field3;
					klInfo[1][3] = field4;
			    }
			    
			    stmt.close();
			    closeConnection(con);
			    result.close();
			    
			    return klInfo;
			}
			catch (Exception e) 
			{
				e.printStackTrace();
			}
			return null;	
		}

   /*
	* 	consolePrint() is used for debugging. It takes in a 2d array and the length
	* 	of the rows (ie # of columns) in the 2d array as input and outputs the contents
	*   to the console.
	* 
	* 		Steven Whaley - created January 30, 2014 - updated January 30, 2014
	*/
	public void consolePrint(String[][] arr, int rowLength)
	{
		int column = 2;
		int row = rowLength;
		 
		for (int i=0; i < column; i++)
		{
			for (int j=0; j < row; j++)
			{
				if(i == 0 && j == 0)
			    {
					System.out.print("\t\t");
			    }
			    if(i == 1 && j == 0)
			    {
			    	System.out.print("\n\t\t" + arr[i][j] + " ");
			    }
			    else
			    {
			    	System.out.print(arr[i][j] + " ");
			    }
			 }
		}
	}
	
	
	
   /*
	*	submitQuery takes an sql query as a String and returns the result as a String.
	*	note: currently prints to console instead of returning
	*
	*	type specifies which table to read from.
	*	type = 1 - laptopDevice
	*	type = 2 - gpsDevice
	*	type = 3 - customer
	*	type = 4 - account
	*
	*	Steven Whaley - created January 26, 2014 - updated January 30, 2014
	*/
//	public void submitQuery(String query, int type)
//	{
//		try
//		{
//			Connection con = openConnection();
//		    PreparedStatement stmt = con.prepareStatement(query);
//		   
//		    ResultSet result = stmt.executeQuery(); 
//		    
//		    System.out.println("\nResults of Query:");
//			
//		    while (result.next()) 
//		    {   
//		    	if (type == 1)
//		    	{
//		    		String field1 = Integer.toString(result.getInt("id"));
//		    		String field2 = result.getString("deviceName");
//		    		String field3 = result.getString("customerId");
//		    		String field4 = result.getString("macAddress");
//		    		
//					System.out.println("id: " + field1 + "\nname: " + field2 + "\ncustomerId: " + field3
//					 + "\nmacAddress: " + field4);
//		    	}
//		    	else if (type == 2)
//		    	{
//		    		String field1 = result.getString("id");
//		    		String field2 = result.getString("name");
//		    		String field3 = result.getString("customerId");
//		    		String field4 = result.getString("latitude");
//		    		String field5 = result.getString("longitude");
//		    		
//		    		System.out.println("id: " + field1 + "\nname: " + field2 + "\ncustomerId: " + field3
//		    							+ "\nlatitude: " + field4 + "\nlongitude: " + field5);
//		    	}
//		    	else if (type == 3)
//		    	{
//		    		String field1 = result.getString("id");
//		    		String field2 = result.getString("phoneNumber");
//		    		String field3 = result.getString("address");
//		    		String field4 = result.getString("email");
//		    		String field5 = result.getString("firstName");
//		    		String field6 = result.getString("lastName");
//		    		
//    				System.out.println("\nid: " + field1 + "\nphoneNumber: "
//    				+ field2 + "\naddress: " + field3 + "\nemail: " + field4 + 
//    				"\nfirstName: " + field5 + "\nlastName: " + field6 + "\n");
//		    	}
//		    	else if (type == 4)
//		    	{
//		    		String field1 = result.getString("id");
//		    		String field2 = result.getString("username");
//		    		String field3 = result.getString("password");
//		    		
//		    		System.out.println("\nid: " + field1 + "\nusername: " + field2 + "\npassword: " + field3 + "\n");
//		    	}
//		    	else
//		    	{
//		    		System.out.println("\nType input invalid. (comments for submitQuery tell which numbers to use)\n");
//		    	}
//			}
//		    stmt.close();
//		    closeConnection(con);
//		    result.close();
//	    }
//	    catch (Exception e) 
//	    {
//	    	e.printStackTrace();
//	    }
//	}
	
	/*
	*	submitUpdate takes an sql update as a String, updates the database
	*	and returns the result as a String.
	*
	*
	*	Steven Whaley - created January 26, 2014 - updated January 30, 2014
	*	Steven Whaley - updated February 5 - prepared statement
	*/
	public void submitUpdate(String query)
	{
		try
		{
			Connection con = openConnection();
		    PreparedStatement stmt = con.prepareStatement(query);
		   
		    stmt.executeUpdate(); 
		   
		    stmt.close();
		    closeConnection(con);
	    }
	    catch (SQLException e) 
	    {
	    	e.printStackTrace();
	    }
	}
	
	
}
	
