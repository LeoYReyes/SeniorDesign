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
	
	/*
	 * Searches through ArrayList<Device> and returns a Device object with specified id if it exists
	 * 
	 * @param id	Device id used to search through the ArrayList
	 * @return		Returns the Device with matching id, otherwise return null
	 */
	public Device getDevice(int id) {
		for(int i = 0; i < deviceList.size(); i++) {
			if(deviceList.get(i).getId() == id) {
				return deviceList.get(i);
			}
		}
		return null;
	}
}
