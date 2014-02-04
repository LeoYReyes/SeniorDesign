package server;
/*
 * 
 * 		Steven Whaley - created February 1
 */

import java.util.ArrayList;

public class LaptopDevice extends Device
{
	ArrayList<IPList> list = new ArrayList<IPList>();
	private String keylog = "";
	
	public LaptopDevice(String id, String name) {
		super(id, name);
	}
	
	public void addIPList(String ips, String time) {
		System.out.println("IPs: " + ips);
		System.out.println("Timestamp: " + time);
		list.add(new IPList(ips, time));
		System.out.println(list.get(0).toString());
	}
	
	public String updateKeylog(String log) {
		keylog += log;
		return keylog;
	}
	
	public boolean keyLogNotEmpty() {
		if(!keylog.isEmpty()) {
			return true;
		}
		else {
			return false;
		}
	}
	
	public String getKeyLog() {
		return keylog;
	}
	
	public String getIpAdresses() {
		String output = "";
		
		for (int i = 0; i < list.size(); i++) {
			output += list.get(i);
		}
		return output;
	}

}
