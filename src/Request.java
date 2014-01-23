import java.util.Observable;


/*
 * This class represents a Request object. Request objects contain data and
 * operations that need to be performed on the data.
 */
public class Request extends Observable implements Comparable<Request> {

	// TODO: assign unique request IDs for each request
	public static final int GEOGRAMDATA = 0;
	public static final int LAPTOPDATA = 1;
	public static final int NEWDEVICE = 2;
	public static final int DEVICEINFO = 3;
	
	private int requestType;
	private int priority;
	private String requestMessage;
	
	public Request(int requestType, String requestMessage, int priority) {
		this.priority = priority;
		this.requestType = requestType;
		this.requestMessage = requestMessage;
		System.out.println("Request Created: " + this.requestType);
		System.out.println("\tMessage: " + this.requestMessage);
		setChanged();
	}
	
	public String getMsg() {
		return requestMessage;
	}

	@Override
	public int compareTo(Request arg0) {
		if(this.priority > arg0.priority) {
			return 1;
		}
		return 0;
	}
}
