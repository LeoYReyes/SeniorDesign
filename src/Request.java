import java.util.Observable;


/*
 * This class represents a Request object. Request objects contain data and
 * operations that need to be performed on the data.
 */
public class Request extends Observable implements Comparable<Request> {

	public static int requestId = 0;
	
	// TODO: assign unique request IDs for each request
	public static final int GEOGRAMDATA = 0;
	public static final int LAPTOPDATA = 1;
	public static final int LOADDEVICE = 2;
	public static final int DEVICEINFO = 3;
	public static final int NEWDEVICE = 4;
	
	private int requestType;
	private int priority;
	private int id;
	private String requestMessage;
	
	/*
	 * Creates a Request with specified Request Type and priority. Assigns a unique request id.
	 * 
	 * @param requestType
	 * @param requestMessage
	 * @param priority
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
	 * 
	 */
	public String getMsg() {
		return requestMessage;
	}
	
	/*
	 * Returns the request type assigned to this Request object
	 */
	public int getRequestType() {
		return requestType;
	}
	
	/*
	 * Returns the unique id assigned to this Request object
	 */
	public int getId() {
		return id;
	}

	@Override
	public int compareTo(Request arg0) {
		if(this.priority > arg0.priority) {
			return 1;
		}
		return 0;
	}
}
