import java.net.Socket;

public class TCPDeviceThread extends Thread {

	private Socket deviceSock = null;
	
	public TCPDeviceThread(Socket deviceSock) {
		this.deviceSock = deviceSock;
	}
	
	public void run() {
		
	}
}
