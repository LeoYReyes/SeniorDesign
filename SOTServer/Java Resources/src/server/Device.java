package server;

/*
 * @author Rizwan Pirani
 * Creates a new device and is the parent to LaptopDevice and GpsDevice
 */
public abstract class Device {
	private String id;
	private String name;
	
	public Device(String id, String name) {
		this.id = id;
		this.name = name;
	}
	
	public String getId() {
		return id;
	}
	
	public String toString() {
		return "Device" + id;
	}
}
