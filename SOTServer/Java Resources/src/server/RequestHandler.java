package server;
import java.util.Observable;
import java.util.Observer;
import java.util.concurrent.PriorityBlockingQueue;

/**
 * @author Leo Reyes
 *
 * A RequestHandler processes Request objects and executes specified operations.
 * This class implements the Observer interface, and is notified when a Request object
 * has been created or changed. It also implements the Runnable interface. The 
 * RequestHandler runs in its own thread and is used to keep track of created
 * Requests and process the Request.
 * 
 * @see java.util.Observer
 */
public class RequestHandler implements Observer, Runnable {

	// PriorityQueue for the storage of Requests, it is a shared resource for later use
	private PriorityBlockingQueue<Request> requestQueue = new PriorityBlockingQueue<Request>();
	// Reference to a DeviceController
	private DeviceController deviceController = null;
	// Reference to a DBController
	private DBController dbController = null;
	
	/*
	 * Constructor
	 */
	public RequestHandler() {
		System.out.println("Request Handler Initiated");
	}
	
	/*
	 * (non-Javadoc)
	 * @see java.util.Observer#update(java.util.Observable, java.lang.Object)
	 * Updates the requestQueue when a new Request is created.
	 */
	@Override
	public void update(Observable observable, Object obj) {
		System.out.println("Request created by thread: " + obj);
		// Add newly created request to queue
		synchronized(requestQueue) {
			requestQueue.add((Request)observable);
			System.out.println("Request added to queue");
		}
	}

	/*
	 * (non-Javadoc)
	 * @see java.lang.Runnable#run()
	 * Processing of the Requests
	 */
	@Override
	public void run() {
		while(true) {
			try {
				// Take the Request in the front of queue
				Request currRequest = requestQueue.take();
				System.out.println("Processing: Request" + currRequest.getId());
				// Execute operation based on request type
				switch(currRequest.getRequestType()) {
					case Request.UPDATE_DEVICE_IP:
						// Updates the IP Address of a Device
						System.out.println("Request: UPDATE IP");
						deviceController.updateIp(currRequest.getDeviceId(), currRequest.getMsg());
						break;
					case Request.UPDATE_DEVICE_KEYLOG:
						// Updates the Keylog of a Device
						System.out.println("Request: UPDATE KEY LOGS");
						deviceController.updateKeyLog(currRequest.getDeviceId(), currRequest.getMsg());
						break;
					case Request.UPDATE_DEVICE_GPS:
						//Updates the GPS coordinates of a Device
						//TODO: create device controller method to update device's coordinates
						System.out.println("Request: UPDATE GPS");
						break;
					case Request.LOAD_DEVICE:
						System.out.println("Request: LOAD DEVICE");
						// Fetch device data from DB and create new Device with controller
						//TODO: Distinguish between devices and create appropriate
						//		device type. i.e. LaptopDevice or GpsDevice
						System.out.println("Device ID: " + currRequest.getMsg());
						// Retrieves LaptopDevice info. currRequest.getMsg() is expected
						// to be an id (MAC Address)
						String[][] deviceInfo = dbController.getLaptopDeviceInfo(currRequest.getMsg());
						// Look at structure of what the DBController returns to know what
						// the array represents in deviceInfo
						//TODO: Maybe create some static constants to represent columns
						//		and rows of the array. Would be done in DBController.
						deviceController.newLaptopDevice(deviceInfo[1][3], deviceInfo[1][1]);
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
	
	/*
	 * Adds a DBController for use by RequestHandler.
	 * 
	 * @param dbController	DBController to attach
	 */
	public void addDbController(DBController dbController) {
		this.dbController = dbController;
	}

}
