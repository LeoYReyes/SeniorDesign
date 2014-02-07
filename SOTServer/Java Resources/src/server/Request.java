package server;
import java.util.Observable;

/**
 * @author Leo Reyes
 * 
 * This class represents a Request object. Request objects contain data and
 * operations that need to be performed on the data. Extends the Observable
 * class. 
 * @see java.util.Observable
 * 
 */
public class Request extends Observable implements Comparable<Request> {

	// Unique ID given to each request
	public static int requestId = 0;
	
	// Request Types
	public static final int UPDATE_DEVICE_IP = 2;
	public static final int UPDATE_DEVICE_KEYLOG = 3;
	public static final int UPDATE_DEVICE_GPS = 5;
	public static final int LOAD_DEVICE = 4;
	
	private int requestType;		// Type of request
	private int priority;			// Priority of the request
	private int id;					// Unique ID of the request
	private String requestMessage;	// Message of the request
	private String deviceId;		// ID of the Device associated with the request
	
	/*
	 * Creates a Request with specified Request Type and priority. Assigns 
	 * a unique request id.
	 * 
	 * @param requestType - of type of request to be created
	 * @param requestMessage - Contents of the request
	 * @param priority - Priority given to the request
	 */
	public Request(int requestType, String requestMessage, int priority) {
		this.priority = priority;
		this.requestType = requestType;
		this.requestMessage = requestMessage;
		this.id = requestId++;
		System.out.println("Request " + id + " Created. Type: " + this.requestType);
		System.out.println("\tMessage: " + this.requestMessage);
		setChanged();
	}
	
	/*
	 * Creates a Request with specified Request Type, message, priority, and a
	 * deviceId associated with the request. Calls other constructor from within
	 * 
	 * @param requestType
	 * @param reuqestMessage
	 * @param priority
	 * @param deviceId
	 */
	public Request(int requestType, String requestMessage, int priority, String deviceId) {
		this(requestType, requestMessage, priority);
		this.deviceId = deviceId;
	}
	
	/*
	 * Returns the message of the Request
	 */
	public String getMsg() {
		return requestMessage;
	}
	
	/*
	 * Returns the request type assigned to this Request 
	 */
	public int getRequestType() {
		return requestType;
	}
	
	/*
	 * Returns the deviceId associated with the Request
	 */
	public String getDeviceId() {
		return deviceId;
	}
	
	/*
	 * Returns the unique id assigned to this Request 
	 */
	public int getId() {
		return id;
	}
	
	/*
	 * (non-Javadoc)
	 * @see java.lang.Comparable#compareTo(java.lang.Object)
	 */
	@Override
	public int compareTo(Request t) {
		if(this.priority > t.priority) {
			return 1;
		}
		return 0;
	}
}
