import java.util.Observable;
import java.util.Observer;
import java.util.concurrent.PriorityBlockingQueue;

/*
 * A RequestHandler processes Request objects and executes specified operations.
 * This class implements the Observer interface, and is notified when a Request object
 * has been created or changed. 
 */
public class RequestHandler implements Observer, Runnable {

	private PriorityBlockingQueue<Request> requestQueue = new PriorityBlockingQueue<Request>();
	private DeviceController deviceController = null;
	
	public RequestHandler() {
		System.out.println("Request Handler Initiated");
	}
	
	@Override
	public void update(Observable arg0, Object arg1) {
		System.out.println("Request created by thread: " + arg1);
		// Add newly created request to queue
		synchronized(requestQueue) {
			requestQueue.add((Request)arg0);
			System.out.println("Request added to queue");
		}
	}

	@Override
	public void run() {
		while(true) {
			try {
				// Take and store Request in the head of queue
				Request currRequest = requestQueue.take();
				System.out.println("Processing: " + currRequest.getMsg());
				// Execute operation based on request type
				switch(currRequest.getRequestType()) {
					case Request.GEOGRAMDATA:
						break;
					case Request.LAPTOPDATA:
						break;
					case Request.LOADDEVICE:
						break;
					case Request.DEVICEINFO:
						break;
					case Request.NEWDEVICE:
						deviceController.newDevice(currRequest.getId());
						break;
					default:
						System.out.println("Unrecognized Request Type");
						break;
				}
			} catch (InterruptedException e) {
				System.out.println(e);
			}
		}
	}
	
	/*
	 * Adds a DeviceController for use by RequestHandler.
	 * 
	 * @param deviceController	DeviceController to attach
	 */
	public void addDeviceController(DeviceController deviceController) {
		this.deviceController = deviceController;
	}

}
