package server;
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
*   		January 23 - started researching, wrote most of the database connection code, but didn't have 
*   					 mysql and connector/j configured correctly at first. 
*   
*   		Steven Whaley:		
*   		January 25 - successfully connected to database, but code is all in main.
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
*   		January 30 - added get_account_info() and get_customer_info(). Added console_print() for debugging, 
*   					 which is used for viewing the contents of the 2d arrays that the other methods are
*   					 outputting. various database changes and code changes. String ID inputs now. 
*   
*   TO DO: interfacing with server code, pushing to github, updating submit_update()/submit_query(), and
*   	   also what other functionality is needed?
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
	*	getLaptopDeviceInfo takes a String for the laptop device id and returns a 2D array containing
	*	field names on the top row and the values associated with those fields in the second row.
	*
	*	Steven Whaley - updated January 30 - revised from getDeviceInfo()
	*/
	public String[][] getLaptopDeviceInfo(String id_in)
	{
		//First need to establish a connection to the database
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
			//note: password must be null, not ""
		    Connection con = DriverManager.getConnection(url, "root", "toor");
		    
		    //System.out.println("URL: " + url);
		    //System.out.println("Connection: " + con);
		    
		    Statement stmt = con.createStatement();
		    
		    con.setAutoCommit(false);
		    con.commit();
			
		    String query = "SELECT * FROM laptopDevice WHERE macAddress = \'" + id_in + "\'";
		    
		    // System.out.println("Here's how the query looks: " + query);
		    
		    ResultSet result = stmt.executeQuery(query);
	      
		    int column_length = 2;
		    int row_length = 4;
		    
		    String[][] laptopDeviceInfo = new String[column_length][row_length];
		     
		    laptopDeviceInfo[0][0] = "id";
		    laptopDeviceInfo[0][1] = "deviceName";
		    laptopDeviceInfo[0][2] = "customerId";
		    laptopDeviceInfo[0][3] = "macAddress";
		    
		    if (result.next()) 
		    {   //process results
					String id = result.getString("id");
					String name = result.getString("deviceName");
					String customer_id = result.getString("customerId");
					String macAddress = result.getString("macAddress");
					 	
					laptopDeviceInfo[1][0] = id;
					laptopDeviceInfo[1][1] = name;
					laptopDeviceInfo[1][2] = customer_id;
					laptopDeviceInfo[1][3] = macAddress;
		    }
		    
		    stmt.close();
		    con.close();
		    result.close();
		    
		    return laptopDeviceInfo;
	    }
	    catch (Exception e) 
	    {
	    	e.printStackTrace();
	    }
	    return null;
	}
	
	/*
	*	getGpsDeviceInfo takes a String for the gps device id and returns a 2D array containing
	*	field names on the top row and the values associated with those fields in the second row.
	*
	*	Steven Whaley - updated January 30 - revised from getDeviceInfo()
	*/
	public String[][] getGpsDeviceInfo(String id_in)
	{
		//First need to establish a connection to the database
		String url = "jdbc:mysql://localhost/trackerdb";

	    try 
	    {
	    	//initialize JDBC driver
	    	Class.forName("com.mysql.jdbc.Driver");
	    }
	    catch( Exception e ) 
	    {
	     	System.out.println("Unable to load driver.");
	    	e.printStackTrace();
	    }
	    
	    try 
	    {
			//parameters for url, username and password - defaulted to root and null
			//note: password must be null, not ""
		    Connection con = DriverManager.getConnection(url, "root", "toor");
		    
		    //System.out.println("URL: " + url);
		    //System.out.println("Connection: " + con);
		    
		    Statement stmt = con.createStatement();
		    
		    con.setAutoCommit(false);
		    con.commit();
			
		    String query = "SELECT * FROM gps_device WHERE id = " + id_in;
		    
		    // System.out.println("Here's how the query looks: " + query);
		    
		    ResultSet result = stmt.executeQuery(query);
	      
		    int column_length = 2;
		    int row_length = 5;
		    
		    String[][] gdevice_info = new String[column_length][row_length];
		     
		    
		    String field1 = "id";
		    String field2 = "name";
		    String field3 = "customer_id";
		    String field4 = "latitude";
		    String field5 = "longitude";
		    
		    gdevice_info[0][0] = field1;
		    gdevice_info[0][1] = field2;
		    gdevice_info[0][2] = field3;
		    gdevice_info[0][3] = field4;
		    gdevice_info[0][4] = field5;
		    
		    if (result.next()) 
		    {   //process results
					field1 = result.getString(field1);
					field2 = result.getString(field2);
					field3 = result.getString(field3);
					field4 = result.getString(field4);
					field5 = result.getString(field5);
					 	
					gdevice_info[1][0] = field1;
					gdevice_info[1][1] = field2;
					gdevice_info[1][2] = field3;
					gdevice_info[1][3] = field4;
					gdevice_info[1][4] = field5;
		    }
		    
		    stmt.close();
		    con.close();
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
	*	Steven Whaley - updated January 30 - revised from getAccountInfo()
	*/
	public String[][] getAccountInfo(String id_in)
	{
		//First need to establish a connection to the database
		String url = "jdbc:mysql://localhost/trackerdb";

	    try 
	    {
	    	//initialize JDBC driver
	    	Class.forName("com.mysql.jdbc.Driver");
	    }
	    catch( Exception e ) 
	    {
	     	System.out.println("Unable to load driver.");
	    	e.printStackTrace();
	    }
	    
	    try 
	    {
			//parameters for url, username and password - defaulted to root and null
			//note: password must be null, not ""
		    Connection con = DriverManager.getConnection(url, "root", "toor");
		    
		    //System.out.println("URL: " + url);
		    //System.out.println("Connection: " + con);
		    
		    Statement stmt = con.createStatement();
		    
		    con.setAutoCommit(false);
		    con.commit();
			
		    String query = "SELECT * FROM account WHERE id = " + id_in;
		    
		    // System.out.println("Here's how the query looks: " + query);
		    
		    ResultSet result = stmt.executeQuery(query);
	      
		    int column_length = 2;
		    int row_length = 3;
		    
		    String[][] account_info = new String[column_length][row_length];
		     
		    
		    String field1 = "id";
		    String field2 = "username";
		    String field3 = "password";
		    
		    account_info[0][0] = field1;
		    account_info[0][1] = field2;
		    account_info[0][2] = field3;
		    
		    if (result.next()) 
		    {   //process results
					field1 = result.getString(field1);
					field2 = result.getString(field2);
					field3 = result.getString(field3);
					
					 	
					account_info[1][0] = field1;
					account_info[1][1] = field2;
					account_info[1][2] = field3;
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
	*	Steven Whaley - updated January 30 - revised from getCustomerInfo()
	*/
	public String[][] getCustomerInfo(String id_in)
	{
		//First need to establish a connection to the database
		String url = "jdbc:mysql://localhost/trackerdb";

	    try 
	    {
	    	//initialize JDBC driver
	    	Class.forName("com.mysql.jdbc.Driver");
	    }
	    catch( Exception e ) 
	    {
	     	System.out.println("Unable to load driver.");
	    	e.printStackTrace();
	    }
	    
	    try 
	    {
			//parameters for url, username and password - defaulted to root and null
			//note: password must be null, not ""
		    Connection con = DriverManager.getConnection(url, "root", "toor");
		    
		    //System.out.println("URL: " + url);
		    //System.out.println("Connection: " + con);
		    
		    Statement stmt = con.createStatement();
		    
		    con.setAutoCommit(false);
		    con.commit();
			
		    String query = "SELECT * FROM customer WHERE id = " + id_in;
		    
		    // System.out.println("Here's how the query looks: " + query);
		    
		    ResultSet result = stmt.executeQuery(query);
	      
		    int column_length = 2;
		    int row_length = 6;
		    
		    String[][] customer_info = new String[column_length][row_length];
		     
		    String field1 = "id";
		    String field2 = "phone_number";
		    String field3 = "address";
		    String field4 = "email";
		    String field5 = "first_name";
		    String field6 = "last_name";
		    
		    customer_info[0][0] = field1;
		    customer_info[0][1] = field2;
		    customer_info[0][2] = field3;
		    customer_info[0][3] = field4;
		    customer_info[0][4] = field5;
		    customer_info[0][5] = field6;
		    
		    if (result.next()) 
		    {   //process results
					field1 = result.getString(field1);
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
		    con.close();
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
	 * 	consolePrint() used for debugging
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
	*	type = 1 - device
	*	type = 2 - customer
	*	type = 3 - account
	*
	*	Steven Whaley - created January 26, 2014
	*/
//	public void submitQuery(String query, int type)
//	{
//		//First need to establish a connection to the database
//		//this is looking on the local machine for a database named trackerdb
//		String url = "jdbc:mysql://localhost/trackerdb";
//		int temp = type;
//
//	    try 
//	    {
//	    	//initialize JDBC driver
//	    	Class.forName("com.mysql.jdbc.Driver");
//	    }
//	    catch( Exception e ) 
//	    {
//	     	System.out.println("Unable to load driver.");
//	    	e.printStackTrace();
//	    }
//	    
//	    try 
//	    {
//			//parameters for url, username and password - defaulted to root and null
//			//note: password must be null, not ""
//		    Connection con = DriverManager.getConnection(url, "root", null);
//		    
//		   //System.out.println("URL: " + url);
//		   //System.out.println("Connection: " + con);
//		    
//		    Statement stmt = con.createStatement();
//		    
//		    con.setAutoCommit(false);
//		    con.commit();
//			
//		    ResultSet result = stmt.executeQuery(query); 
//		    
//		    System.out.println("\nResults of Query:");
//			
//		    while (result.next()) 
//		    {   //process results
//		    	if (type == 1)
//		    	{
//					String dname = result.getString("name");
//					int dlocation = result.getInt("location");
//					String downer = result.getString("owner");
//					int did = result.getInt("id");
//					System.out.println("name: " + dname + "\nlocation: " + dlocation + "\nowner: " + downer + "\nid: " + did + "\n");
//		    	}
//		    	else if (type == 2)
//		    	{
//		    		int cid = result.getInt("id");
//    				String cname = result.getString("name");
//    				String cphone_number = result.getString("phone_number");
//    				String caddress = result.getString("address");
//    				String cemail = result.getString("email");
//    				System.out.println("\ncustomer name: " + cname + "\ncustomer phone_number: "
//    				+ cphone_number + "\ncustomer address" + caddress + "\ncustomer email" + cemail + "\n");
//		    	}
//		    	else if (type == 3)
//		    	{
//		    		int aid = result.getInt("id");
//		    		String ausername = result.getString("username");
//		    		String apassword = result.getString("password");
//		    		System.out.println("\naccount id: " + aid + "\naccount username: " + ausername + "\napassword: " + apassword + "\n");
//		    	}
//		    	else
//		    	{
//		    		System.out.println("\nType input invalid. (comments for submit_query tell which numbers to use)\n");
//		    	}
//			}
//		    stmt.close();
//		    con.close();
//		    result.close();
//	    }
//	    catch (Exception e) 
//	    {
//	    	e.printStackTrace();
//	    }
//	}
	
	/*
	*	submit_query takes an sql update as a String, updates the database
	*	and returns the result as a String.
	*
	*	note: same as submit_query() as of jan 27 needs to be updated
	*
	*	Steven Whaley - created January 26, 2014
	*/
//	public void submitUpdate(String query, int type)
//	{
//		//First need to establish a connection to the database
//		//this is looking on the local machine for a database named trackerdb
//		String url = "jdbc:mysql://localhost/trackerdb";
//		int temp = type;
//
//	    try 
//	    {
//	    	//initialize JDBC driver
//	    	Class.forName("com.mysql.jdbc.Driver");
//	    }
//	    catch( Exception e ) 
//	    {
//	     	System.out.println("Unable to load driver.");
//	    	e.printStackTrace();
//	    }
//	    
//	    try 
//	    {
//			//parameters for url, username and password - defaulted to root and null
//			//note: password must be null, not ""
//		    Connection con = DriverManager.getConnection(url, "root", null);
//		    
//		    //System.out.println("URL: " + url);
//		    //System.out.println("Connection: " + con);
//		    
//		    Statement stmt = con.createStatement();
//		    
//		    con.setAutoCommit(false);
//		    con.commit();
//			
//		    ResultSet result = stmt.executeQuery(query); 
//		    
//		    System.out.println("\nResults of Query:");
//			
//		    while (result.next()) 
//		    {   //process results
//		    	if (type == 1)
//		    	{
//					String dname = result.getString("name");
//					int dlocation = result.getInt("location");
//					String downer = result.getString("owner");
//					int did = result.getInt("id");
//					System.out.println("name: " + dname + "\nlocation: " + dlocation + "\nowner: " + downer + "\nid: " + did + "\n");
//		    	}
//		    	else if (type == 2)
//		    	{
//		    		int cid = result.getInt("id");
//    				String cname = result.getString("name");
//    				String cphone_number = result.getString("phone_number");
//    				String caddress = result.getString("address");
//    				String cemail = result.getString("email");
//    				System.out.println("\ncustomer name: " + cname + "\ncustomer phone_number: "
//    				+ cphone_number + "\ncustomer address" + caddress + "\ncustomer email" + cemail + "\n");
//		    	}
//		    	else if (type == 3)
//		    	{
//		    		int aid = result.getInt("id");
//		    		String ausername = result.getString("username");
//		    		String apassword = result.getString("password");
//		    		System.out.println("\naccount id: " + aid + "\naccount username: " + ausername + "\napassword: " + apassword + "\n");
//		    	}
//		    	else
//		    	{
//		    		System.out.println("\nType input invalid. (comments for submit_query tell which numbers to use)\n");
//		    	}
//			}
//		    stmt.close();
//		    con.close();
//		    result.close();
//	    }
//	    catch (Exception e) 
//	    {
//	    	e.printStackTrace();
//	    }
//	}
	
	
}
	

