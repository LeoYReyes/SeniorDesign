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
				System.out.println("Processing: " + requestQueue.take().getMsg());
			} catch (InterruptedException e) {
				System.out.println(e);
			}
		}
	}

}
