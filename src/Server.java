
import java.io.IOException;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.ArrayList;

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
			serverSock = new ServerSocket(10000);
			requestHandler.addDeviceController(deviceController);
			new Thread(requestHandler).start();
		}
		catch(IOException e) {
			System.out.println(e);
		}
	}
	
	public static TCPDeviceThread findDeviceThread() {
		return null;
	}
	
	public static void main(String args[]) {
		
		Server server = new Server();
		
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
			} catch (Exception e) {
				System.out.println(e);
			}
		}
	}
}
