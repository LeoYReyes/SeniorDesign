package server;
import java.io.BufferedInputStream;
import java.net.Socket;
import java.net.SocketException;

/**
 * @author Leo Reyes
 *
 * Thread for a connection over TCP. Parses messages received and creates a request 
 * that will be handled by the RequestHandler.
 */
public class TCPDeviceThread extends Thread {
	// Number of messages exchanged over the connection
	private int messageCount;
	String deviceId = null;
	// Client socket
	private Socket deviceSock = null;
	// Representation of the device connected to this thread
	private Device device = null;
	// Request
	private Request request = null;
	// RequestHandler
	private RequestHandler requestHandler = null;
	// DeviceController
	private DeviceController deviceController = null;
	
	/*
	 * Creates a TCPDeviceThread referencing to a specified Socket and RequestHandler
	 * 
	 * @param deviceSock		Socket connection that this thread will connected to
	 * @param requestHandler	Reference to the Server's RequestHandler
	 */
	public TCPDeviceThread(Socket deviceSock, RequestHandler requestHandler, DeviceController deviceController) {
		this.deviceSock = deviceSock;
		this.requestHandler = requestHandler;
		this.deviceController = deviceController;
		messageCount = 0;
		System.out.println("new connection thread created");
	}
	
	public void run() {
		//TODO: No DeviceController in this class, create requests and return values
		//		to check for devices instead
		System.out.println("new thread running");
		try {
			// incoming message buffer
			BufferedInputStream in = new BufferedInputStream(deviceSock.getInputStream());
			byte[] inMessage = new byte[1024];
			// Number of bytes read
			int bytesRead = 0;
			while((bytesRead = in.read(inMessage)) > 0) {
				// Check to see if first message, if true setup device for connection
				if(messageCount < 1) {
					deviceId = new String(inMessage, 0, 12);
					System.out.println("DeviceMAC string: " + deviceId);
					//Check to see if Device exists in controller, if not load device from database
					if(deviceController.getDevice(deviceId) == null) {
						// If device was not found, load device from database
						// Create new LOAD_DEVICE Request
						request = new Request(Request.LOAD_DEVICE, deviceId, 1);
						request.addObserver(requestHandler);
						request.notifyObservers(currentThread().getId());
						Thread.sleep(1000);
						device = deviceController.getDevice(deviceId);
						//TODO: NULLPOINTER??//System.out.println("Thread Device: " + device.getId());
					}
					messageCount++;
				} 
				else {
					System.out.println("Received message: " + new String(inMessage));
					Byte opCode = inMessage[0];
					// Create new request and add to RequestHandler
					request = new Request(opCode.intValue(), new String(inMessage, 1, (bytesRead)), 1, deviceId);
					request.addObserver(requestHandler);
					request.notifyObservers(currentThread().getId());
					messageCount++;
				}
				for(int i = 0; i< inMessage.length; i++) {
					inMessage[i] = 0x00;
				}
			}// end of while((bytesRead = in.read(inMessage)) > 0)
			
		} // end of try
		catch (Exception e) {
			System.out.println(e);
		}
	}
	
	public Socket getSock() {
		return deviceSock;
	}
	public String getDeviceId() {
		return deviceId;
	}
	public Device getDevice() {
		return device;
	}
}
