import java.util.ArrayList;

/*
 * 
 */
public class DeviceController {

	private ArrayList<Device> deviceList = new ArrayList<Device>();
	
	/*
	 * Class constructor
	 */
	public DeviceController() {
		
	}
	
	/*
	 * Instantiates new Device object and add to the list of devices
	 * @param id 	the unique id of the device being created.
	 */
	public void newDevice(int id) {
		deviceList.add(new Device(id));
		System.out.println("New device created and added");
		System.out.println("DeviceList: " + deviceList.toString());
	}
}
