/*
 * 
 * 		Steven Whaley - created February 1
 */

import java.util.ArrayList;

public class LaptopDevice extends Device
{
	ArrayList<String> list = new ArrayList<String>();
	private String keylog = "";
	
	public LaptopDevice(int id, String name)
	{
		super(id, name);
	}
	
	public ArrayList<String> addIpAddress(String ip)
	{
		list.add(ip);
		
		return list;
	}
	
	public String updateKeylog(String log)
	{
		keylog += log;
		return keylog;
	}
	
	public boolean keyLogNotEmpty()
	{
		if(!keylog.isEmpty())
		{
			return true;
		}
		else
		{
			return false;
		}
	}
	
	public String getKeyLog()
	{
		return keylog;
	}
	
	public String getIpAdresses()
	{
		String output = "";
		
		for (int i = 0; i < list.size(); i++)
		{
			output += list.get(i);
		}
		return output;
	}

}
