package server;
/*
 * @author: Rizwan Pirani
 * 		Steven Whaley - Feb 1
 * Child of Device sets up a gps device and takes in its latitude and longitude.
 */
public class GpsDevice extends Device
{
	private float latitude;
	private float longitude;
	
	public GpsDevice(String id, String name, float latitude, float longitude)
	{
		super(id, name);
		this.latitude = latitude;
		this.longitude = longitude;
	}
	
	public float getLatitude()
	{
		return latitude;
	}
	
	public float getLongitude()
	{
		return longitude;
	}
}
