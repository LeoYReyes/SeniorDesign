package server;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.sql.Timestamp;
import java.util.AbstractSequentialList;
import java.util.LinkedList;
import java.util.List;
import java.util.ListIterator;

public class IPList extends AbstractSequentialList<InetAddress> {

	List<InetAddress> list;
	Timestamp timestamp;
	
	public IPList() {
		list = new LinkedList<InetAddress>();
		timestamp = new Timestamp(System.currentTimeMillis());
	}
	
	public IPList(String ips, String timestamp) {
		this();
		String current = ips;
		while(current.indexOf("~") > 0) {
			add(current.substring(0, current.indexOf("~")));
			current = current.substring(current.indexOf("~") + 1, current.length());
		}
		add(current.substring(0, current.length()));
	}
	
	public void add(String ipAddress) {
		try {
			list.add(InetAddress.getByName(ipAddress));
		} 
		catch (UnknownHostException e) {
			e.printStackTrace();
		}
	}
	
	public Timestamp getTimestamp() {
		return timestamp;
	}
	
	@Override
	public ListIterator<InetAddress> listIterator(int index) {
		ListIterator<InetAddress> li = list.listIterator();
	    while(li.hasNext()){
	      //System.out.println(li.next());
	    }
	    return null;
	}

	@Override
	public int size() {
		// TODO Auto-generated method stub
		return list.size();
	}
	
	public String toString() {
		String str = "";
		str += "Timestamp: " + timestamp.toString() + "\n";
		str += "IPList:\n";
		ListIterator<InetAddress> li = list.listIterator();
		while(li.hasNext()) {
			str += "\t" + li.next().getHostAddress() + "\n";
		}
		return str;
	}

}
