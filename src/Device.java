
/*
 * 
 */
public class Device {
	//TODO: make class abstract and create child classes for different device types
	private int id;
	
	public Device(int id) {
		this.id = id;
	}
	
	public int getId() {
		return id;
	}
	
	public String toString() {
		return "Device" + id;
	}
}
