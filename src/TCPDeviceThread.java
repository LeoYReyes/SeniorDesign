import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.Socket;

/*
 * Thread for a connection over TCP. Parses message received and creates a request 
 * that will be handled by the server.
 */
public class TCPDeviceThread extends Thread {

	private Socket deviceSock = null;
	// Representation of the device connected
	private Device device = null;
	// Request
	private Request request = null;
	// RequestHandler, will be assigned server's request handler
	private RequestHandler requestHandler = null;
	
	public TCPDeviceThread(Socket deviceSock, RequestHandler requestHandler) {
		this.deviceSock = deviceSock;
		this.requestHandler = requestHandler;
	}
	
	public void run() {
		try {
			// incoming message buffer
			BufferedReader in = new BufferedReader(new InputStreamReader(deviceSock.getInputStream()));
			String inMessage = in.readLine();
			System.out.println("Received message: " + inMessage);
			// Create new request and add to RequestHandler
			//TODO: request type
			request = new Request();
			request.addObserver(requestHandler);
			
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	
	public Device getDevice() {
		return device;
	}
}
