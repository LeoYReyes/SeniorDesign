import java.util.Observable;


/*
 * 
 */
public class Request extends Observable {

	private int requestType;
	
	public Request() {
		hasChanged();
		notifyObservers(this);
	}
}
