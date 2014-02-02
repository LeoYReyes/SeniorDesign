/*
*
*	Steven Whaley - created: January 23, 2014  - last updated: January 30, 2014
*
*	note: 	consider revising so connection uses DataSource object instead of DriverManager
*	note: 	mysql-connector-java-5.1.28-bin.jar placed in referenced libraries 
*   note: 	on mac -  /usr/local/mysql
*   note: 	consider revising to use preparedstatement instead of statement for security against sql injections
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
*   		Steven Whaley:
*   		February 1 - updated to match database changes
*   
*   TO DO: database design and code updates
*  
*   	   -also what other functionality is needed?
*   
*/

import java.sql.*;
//import javax.sql.*;

public class DBController
{
	public DBController()
	{
		
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
	*/
	public String[][] getLaptopDeviceInfo(String id_in)
	{
		try
		{
			Connection con = openConnection();
			
		    Statement stmt = con.createStatement();
		    
		    String query = "SELECT * FROM laptop_device WHERE id = \'" + id_in +"\'";
		    
		    ResultSet result = stmt.executeQuery(query);
	      
		    int column_length = 2;
		    int row_length = 4;
		    
		    String[][] ldevice_info = new String[column_length][row_length];
		     
		    
		    String field1 = "id";
		    String field2 = "deviceName";
		    String field3 = "customerId";
		    String field4 = "macAddress";
		    
		    ldevice_info[0][0] = field1;
		    ldevice_info[0][1] = field2;
		    ldevice_info[0][2] = field3;
		    ldevice_info[0][3] = field4;
		    
		    if (result.next()) 
		    {   	
				field1 = Integer.toString(result.getInt(field1));
				field2 = result.getString(field2);
				field3 = Integer.toString(result.getInt(field3));
				field4 = result.getString(field4);
					 	
				ldevice_info[1][0] = field1;
				ldevice_info[1][1] = field2;
				ldevice_info[1][2] = field3;
				ldevice_info[1][3] = field4;
		    }
		    
		    stmt.close();
		    closeConnection(con);
		    result.close();
		    
		    return ldevice_info;
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
	*/
	public String[][] getGpsDeviceInfo(String id_in)
	{
		try
		{
		    Connection con = openConnection();
		    Statement stmt = con.createStatement();
		    
		    String query = "SELECT * FROM gps_device WHERE id = \'" + id_in + "\'";
		       
		    ResultSet result = stmt.executeQuery(query);
	      
		    int column_length = 2;
		    int row_length = 5;
		    
		    String[][] gdevice_info = new String[column_length][row_length];
		     
		    String field1 = "id";
		    String field2 = "name";
		    String field3 = "customerId";
		    String field4 = "latitude";
		    String field5 = "longitude";
		    
		    gdevice_info[0][0] = field1;
		    gdevice_info[0][1] = field2;
		    gdevice_info[0][2] = field3;
		    gdevice_info[0][3] = field4;
		    gdevice_info[0][4] = field5;
		    
		    if (result.next()) 
		    {  
				field1 = Integer.toString(result.getInt(field1));
				field2 = result.getString(field2);
				field3 = Integer.toString(result.getInt(field3));
				field4 = result.getString(field4);
				field5 = result.getString(field5);
					 	
				gdevice_info[1][0] = field1;
				gdevice_info[1][1] = field2;
				gdevice_info[1][2] = field3;
				gdevice_info[1][3] = field4;
				gdevice_info[1][4] = field5;
		    }
	  
		    stmt.close();
		    closeConnection(con);
		    result.close();
		    
		    return gdevice_info;
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
	*/
	public String[][] getAccountInfo(String id_in)
	{
		try
		{
			Connection con = openConnection();
			
		    Statement stmt = con.createStatement();
		    
		    String query = "SELECT * FROM account WHERE id = \'" + id_in + "\'";
		    
		    ResultSet result = stmt.executeQuery(query);
	      
		    int column_length = 2;
		    int row_length = 4;
		    
		    String[][] account_info = new String[column_length][row_length];
		     
		    String field1 = "id";
		    String field2 = "customerId";
		    String field3 = "username";
		    String field4 = "password";
		    
		    account_info[0][0] = field1;
		    account_info[0][1] = field2;
		    account_info[0][2] = field3;
		    account_info[0][3] = field4;
		    
		    if (result.next()) 
		    {   
				field1 = Integer.toString(result.getInt(field1));
				field2 = Integer.toString(result.getInt(field2));
				field3 = result.getString(field3);
				field4 = result.getString(field4);
						
				account_info[1][0] = field1;
				account_info[1][1] = field2;
				account_info[1][2] = field3;
				account_info[1][3] = field4;
		    }
		    		    
		    stmt.close();
		    con.close();
		    result.close();
		    
		    return account_info;
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
	*/
	public String[][] getCustomerInfo(String id_in)
	{
		 try
		 {
			Connection con = openConnection();
		    Statement stmt = con.createStatement();
		    	
		    String query = "SELECT * FROM customer WHERE id = \'" + id_in + "\'";
		   	    
		    ResultSet result = stmt.executeQuery(query);
	      
		    int column_length = 2;
		    int row_length = 6;
		    
		    String[][] customer_info = new String[column_length][row_length];
		     
		    String field1 = "id";
		    String field2 = "phoneNumber";
		    String field3 = "address";
		    String field4 = "email";
		    String field5 = "firstName";
		    String field6 = "lastName";
		    
		    customer_info[0][0] = field1;
		    customer_info[0][1] = field2;
		    customer_info[0][2] = field3;
		    customer_info[0][3] = field4;
		    customer_info[0][4] = field5;
		    customer_info[0][5] = field6;
		    
		    if (result.next()) 
		    {  
				field1 = Integer.toString(result.getInt(field1));
				field2 = result.getString(field2);
				field3 = result.getString(field3);
				field4 = result.getString(field4);
				field5 = result.getString(field5);
				field6 = result.getString(field6);
					 	
				customer_info[1][0] = field1;
				customer_info[1][1] = field2;
				customer_info[1][2] = field3;
				customer_info[1][3] = field4;
				customer_info[1][4] = field5;
				customer_info[1][5] = field6;
		    }
		    		    
		    stmt.close();
		    closeConnection(con);
		    result.close();
		    
		    return customer_info;
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
	public void consolePrint(String[][] arr, int row_length)
	{
		int column = 2;
		int row = row_length;
		 
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
	*	submit_query takes an sql query as a String and returns the result as a String.
	*	note: currently prints to console instead of returning
	*
	*	type specifies which table to read from.
	*	type = 1 - laptop_device
	*	type = 2 - gps_device
	*	type = 3 - customer
	*	type = 4 - account
	*
	*	Steven Whaley - created January 26, 2014 - updated January 30, 2014
	*/
	public void submitQuery(String query, int type)
	{
		try
		{
			Connection con = openConnection();
		    Statement stmt = con.createStatement();
		   
		    ResultSet result = stmt.executeQuery(query); 
		    
		    System.out.println("\nResults of Query:");
			
		    while (result.next()) 
		    {   
		    	if (type == 1)
		    	{
		    		String field1 = result.getString("id");
		    		String field2 = result.getString("name");
		    		String field3 = result.getString("customer_id");
		    		
					System.out.println("id: " + field1 + "\nname: " + field2 + "\ncustomer_id: " + field3);
		    	}
		    	else if (type == 2)
		    	{
		    		String field1 = result.getString("id");
		    		String field2 = result.getString("name");
		    		String field3 = result.getString("customer_id");
		    		String field4 = result.getString("latitude");
		    		String field5 = result.getString("longitude");
		    		
		    		System.out.println("id: " + field1 + "\nname: " + field2 + "\ncustomer_id: " + field3
		    							+ "\nlatitude: " + field4 + "\nlongitude: " + field5);
		    	}
		    	else if (type == 3)
		    	{
		    		String field1 = result.getString("id");
		    		String field2 = result.getString("phone_number");
		    		String field3 = result.getString("address");
		    		String field4 = result.getString("email");
		    		String field5 = result.getString("first_name");
		    		String field6 = result.getString("last_name");
		    		
    				System.out.println("\nid: " + field1 + "\nphone_number: "
    				+ field2 + "\naddress: " + field3 + "\nemail: " + field4 + 
    				"\nfirst_name: " + field5 + "\nlast_name: " + field6 + "\n");
		    	}
		    	else if (type == 4)
		    	{
		    		String field1 = result.getString("id");
		    		String field2 = result.getString("username");
		    		String field3 = result.getString("password");
		    		
		    		System.out.println("\nid: " + field1 + "\nusername: " + field2 + "\npassword: " + field3 + "\n");
		    	}
		    	else
		    	{
		    		System.out.println("\nType input invalid. (comments for submit_query tell which numbers to use)\n");
		    	}
			}
		    stmt.close();
		    closeConnection(con);
		    result.close();
	    }
	    catch (Exception e) 
	    {
	    	e.printStackTrace();
	    }
	}
	
	/*
	*	submitUpdate takes an sql update as a String, updates the database
	*	and returns the result as a String.
	*
	*	note: same as submit_query() as of jan 27 needs to be updated
	*
	*	Steven Whaley - created January 26, 2014 - updated January 30, 2014
	*/
	public void submitUpdate(String query)
	{
		try
		{
			Connection con = openConnection();
		    Statement stmt = con.createStatement();
		   
		    stmt.executeUpdate(query); 
		   
		    stmt.close();
		    closeConnection(con);
	    }
	    catch (SQLException e) 
	    {
	    	e.printStackTrace();
	    }
	}
	
	
}
	
