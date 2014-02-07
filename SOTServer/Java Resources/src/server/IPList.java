package server;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.sql.Timestamp;
import java.util.AbstractSequentialList;
import java.util.LinkedList;
import java.util.List;
import java.util.ListIterator;

public class IPList extends AbstractSequentialList<String> {

	List<String> list;
	Timestamp timestamp;
	
	public IPList() {
		list = new LinkedList<String>();
		timestamp = new Timestamp(System.currentTimeMillis());
	}
	
	public IPList(String ips) {
		this();
		String current = ips;
		while(current.indexOf("~") > 0) {
			list.add(current.substring(0, current.indexOf("~")));
			current = current.substring(current.indexOf("~") + 1, current.length());
		}
		list.add(current.substring(0, current.indexOf("\n")));
		//add(current.substring(0, current.length()));
	}
	
	public Timestamp getTimestamp() {
		return timestamp;
	}
	
	@Override
	public ListIterator<String> listIterator(int index) {
		ListIterator<String> li = list.listIterator();
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
		ListIterator<String> li = list.listIterator();
		while(li.hasNext()) {
			str += "\t" + li.next() + "\n";
		}
		return str;
	}

}
