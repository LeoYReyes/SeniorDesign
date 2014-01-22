
import java.io.IOException;
import java.net.ServerSocket;
import java.net.Socket;

public class Server {

	private static ServerSocket serverSock;
	private static Socket deviceSock;
	
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
			} catch (IOException e) {
				System.out.println(e);
			}
		}
	}
}
