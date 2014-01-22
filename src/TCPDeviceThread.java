import java.net.Socket;

public class TCPDeviceThread extends Thread {

	private Socket deviceSock;
	
	public TCPDeviceThread(Socket deviceSock) {
		this.deviceSock = deviceSock;
	}
	
	public void run() {
		
	}
}
