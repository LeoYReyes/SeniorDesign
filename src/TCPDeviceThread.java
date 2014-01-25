import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.Socket;

/*
 * Thread for a connection over TCP. Parses messages received and creates a request 
 * that will be handled by the RequestHandler.
 */
public class TCPDeviceThread extends Thread {

	private Socket deviceSock = null;
	// Representation of the device connected
	private Device device = null;
	// Request
	private Request request = null;
	// RequestHandler
	private RequestHandler requestHandler = null;
	
	/*
	 * Creates a TCPDeviceThread referencing to a specified Socket and RequestHandler
	 * 
	 * @param deviceSock		Socket connection that this thread will connected to
	 * @param requestHandler	Reference to the Server's RequestHandler
	 */
	public TCPDeviceThread(Socket deviceSock, RequestHandler requestHandler) {
		this.deviceSock = deviceSock;
		this.requestHandler = requestHandler;
		System.out.println("new connection thread created");
	}
	
	public void run() {
		try {
			System.out.println("new thread running");
			// incoming message buffer
			BufferedReader in = new BufferedReader(new InputStreamReader(deviceSock.getInputStream()));
			String inMessage = in.readLine();
			System.out.println("Received message: " + inMessage);
			// Create new request and add to RequestHandler
			//TODO: request type, parse message to check who sent the message
			//TODO: instead of created new Requests here, make RequestHandler create new Requests
			request = new Request(Request.NEWDEVICE, inMessage, 1);
			request.addObserver(requestHandler);
			request.notifyObservers(currentThread().getId());
			
		} catch (IOException e) {
			System.out.println(e);
		}
	}
	
	public Device getDevice() {
		return device;
	}
}
