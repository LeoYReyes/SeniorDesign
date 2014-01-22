import java.util.Observable;
import java.util.Observer;
import java.util.concurrent.BlockingQueue;


public class RequestHandler implements Observer {

	private BlockingQueue<Request> requestQueue;
	
	public RequestHandler() {
		
	}
	
	@Override
	public void update(Observable arg0, Object arg1) {
		// Add newly created request to queue
		synchronized(requestQueue) {
			requestQueue.add((Request)arg1);
		}
	}

}
