
import java.io.IOException;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.ArrayList;

public class Server {

	private static ServerSocket serverSock;
	private static Socket deviceSock;
	private static ArrayList<Socket> deviceConnections = new ArrayList<Socket>();
	
	public Server() {
		
	}
	
	public static void main(String args[]) {
		int portNumber = 10000;
		try {
			serverSock = new ServerSocket(portNumber);
		}
		catch (IOException e) {
			System.out.println(e);
		}
		
		while(true) {
			try {
				deviceSock = serverSock.accept();
				deviceConnections.add(deviceSock);
			} 
			catch (IOException e) {
				System.out.println(e);
			}
		}
	}
}
