package sot.seniordesign.smsgateway;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.InputStreamReader;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.net.SocketAddress;

import android.os.AsyncTask;

public class TCPAsyncTask extends AsyncTask<String, Integer, Boolean> {
	public static final int NO_CONNECTION = 0;
	public static final int CONNECTION = 1;
	public static final int TIMEOUT = 10000;
	
	private SMSActivity parentActivity;
	private String ip = "";
	private String port = "";
	private volatile boolean keepAlive;

	public TCPAsyncTask(SMSActivity parentAct)
	{
		parentActivity = parentAct;
	}
	
	public void endTask()
	{
		keepAlive = false;
	}

	protected Boolean doInBackground(String... arg) {
		Socket tcp = null;
		BufferedReader fromServer = null;
		DataOutputStream toServer = null;
		
		keepAlive = true;
		try {
			ip = arg[0];
			port = arg[1];
			SocketAddress serverInfo = new InetSocketAddress(ip, Integer.parseInt(port));
			tcp = new Socket();
			tcp.connect(serverInfo, TIMEOUT);
			fromServer = new BufferedReader(new InputStreamReader(tcp.getInputStream()));
			toServer = new DataOutputStream(tcp.getOutputStream());
			publishProgress(CONNECTION);
		}
		catch (Exception e) {
			keepAlive = false;
		}
		
		while (keepAlive)
		{
			try {
				Thread.sleep(1000);
			} catch (InterruptedException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
			}
		}
		
		try {
			fromServer.close();
			toServer.close();
			tcp.close();
		}
		catch (Exception e) {
		}
		
		return Boolean.TRUE;
	}
	
	@Override
	protected void onProgressUpdate(Integer... progress)
	{
		switch (progress[0])
		{
		case NO_CONNECTION:
			
			break;
		case CONNECTION:
			parentActivity.connected();
			break;
		default:
			break;
		}
	}
	
	@Override
	protected void onPostExecute(Boolean bool)
	{
		parentActivity.disconnect();
	}

}