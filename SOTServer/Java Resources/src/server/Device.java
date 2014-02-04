package server;

/*
 * 
 */
public abstract class Device {
	//TODO: make class abstract and create child classes for different device types
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
