
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
	
	public Server() {
		try {
			serverSock = new ServerSocket(10000);
		}
		catch(IOException e) {
			
		}
	}
	
	public static TCPDeviceThread findDeviceThread() {
		return null;
	}
	
	public static void main(String args[]) {
		
		Server server = new Server();
		
		while(true) {
			try {
				// Listen for a connection to server
				deviceSock = serverSock.accept();
				// Create new thread for connection and add to list
				deviceConnections.add(new TCPDeviceThread(deviceSock, requestHandler));
			} 
			catch (IOException e) {
				System.out.println(e);
			}
		}
	}
}
