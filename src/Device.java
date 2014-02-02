
/*
 * 	Steven Whaley - February 1
 */
public abstract class Device 
{
	private int id;
	private String name;
	
	public Device(int id, String name) 
	{
		this.id = id;
		this.name = name;
	}
	
	public int getId() 
	{
		return id;
	}
	
	public String getName()
	{
		return name;
	}
	
	public String toString() 
	{
		return "\nDevice ID: " + id + "\nDevice Name: " + name;
	}
}
