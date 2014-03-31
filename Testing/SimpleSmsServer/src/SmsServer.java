import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.ServerSocket;
import java.net.Socket;


public class SmsServer {
	
	public static void main(String args[]) throws Exception {
		ServerSocket ss = new ServerSocket(10016);
		
		while(true) {
			System.out.println("Waiting for connection...");
			Socket connection = ss.accept();
			System.out.println("Connected");
			BufferedReader fromClient = new BufferedReader(new InputStreamReader(connection.getInputStream()));
			OutputStream toClient = connection.getOutputStream();
			char[] next = new char[1];
			String msgFromSms = "";
			int i = 0;
			while (i < 300)
			{
				i++;
				while (fromClient.ready())
				{
					fromClient.read(next, 0, 1);
					msgFromSms += next[0];
					if (next[0] == '|')
					{
						System.out.println("Received: " + msgFromSms);
						toClient.write(msgFromSms.getBytes());
						msgFromSms = "";
					}
					i = 0;
				}
				Thread.sleep(100);
			}
			fromClient.close();
			toClient.close();
			connection.close();
		}
	}
}
