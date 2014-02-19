package sot.seniordesign.smsgateway;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.net.SocketAddress;

import android.app.PendingIntent;
import android.content.Intent;
import android.os.AsyncTask;
import android.telephony.SmsManager;
import android.telephony.SmsMessage;

/**
 * 
 * @author Chuck
 *
 * Defines a AsyncTask (thread) that handles the TCP connection and SMS sending/receiving
 */
public class TCPAsyncTask extends AsyncTask<String, Integer, Boolean> 
{
	public static final int NO_CONNECTION = 0;
	public static final int CONNECTION = 1;
	public static final int TIMEOUT = 10000;
	
	private SMSActivity parentActivity;
	private String ip = "";
	private String port = "";
	private volatile boolean keepAlive;
	private SMSReceiver smsRcv;
	private String msgFromServer;

	/**
	 * @param parentAct Reference to activity that called it
	 */
	public TCPAsyncTask(SMSActivity parentAct)
	{
		parentActivity = parentAct;
	}
	
	/**
	 * Sets the task to expire as soon as possible
	 */
	public void endTask()
	{
		keepAlive = false;
	}

	/**
	 * Automatically called when executed.
	 * Connects to a tcp server and loops listening for messages. Forwards messages
	 * received as SMS, returns SMS received to the server over TCP.
	 */
	protected Boolean doInBackground(String... arg) 
	{
		smsRcv = new SMSReceiver();
		Socket tcp = null;
		BufferedReader fromServer = null;
		DataOutputStream toServer = null;
		
		keepAlive = true;
		while (keepAlive)
		{
			if (tcp == null || !tcp.isConnected())
			{
				publishProgress(NO_CONNECTION);
				try 
				{
					msgFromServer = "";
					ip = arg[0];
					port = arg[1];
					SocketAddress serverInfo = new InetSocketAddress(ip, Integer.parseInt(port));
					tcp = new Socket();
					tcp.connect(serverInfo, TIMEOUT);
					fromServer = new BufferedReader(new InputStreamReader(tcp.getInputStream()));
					toServer = new DataOutputStream(tcp.getOutputStream());
					publishProgress(CONNECTION);
				}
				catch (Exception e) 
				{
					//keepAlive = false;
					try {
						Thread.sleep(1000);
					} catch (InterruptedException e1) 
					{
					}
				}
			}
			
			while (keepAlive && tcp != null && tcp.isConnected())
			{
				// read SMS message and forward to server
				try 
				{
					//send received SMSs to server
					while (SMSReceiver.hasMsg())
					{
						SmsMessage msg = SMSReceiver.getNextMsg();
						String output = "[" + msg.getOriginatingAddress() + "]" + msg.getMessageBody().toString() + "|";
						toServer.write(output.getBytes("UTF-8"));
					}
				} catch (IOException e) {
					// TODO Auto-generated catch block
					e.printStackTrace();
				}
				
				// Read TCP messages and forward via SMS
				try
				{
					char next;
					while (fromServer.ready())
					{
						next = (char) fromServer.read();
						msgFromServer += next;
						if (next == '|')
						{
							sendSMS(msgFromServer);
							msgFromServer = "";
						}
					}
				} catch (Exception e)
				{
				}
				
				
				try 
				{
					Thread.sleep(1000);
				} catch (InterruptedException e) {
					// TODO Auto-generated catch block
					e.printStackTrace();
				}
			}
		}
		
		try 
		{
			fromServer.close();
			toServer.close();
			tcp.close();
		}
		catch (Exception e) 
		{
		}
		
		return Boolean.TRUE;
	}
	
	/**
	 * Send an SMS
	 * @param phoneNum Phone number as a string of digits
	 * @param msg Message to include in SMS
	 */
	private void sendSMS(String phoneNum, String msg)
	{
		PendingIntent pend = PendingIntent.getActivity(parentActivity, 0, new Intent(parentActivity, SMSActivity.class), 0);
		SmsManager sms = SmsManager.getDefault();
		sms.sendTextMessage(phoneNum, null, msg, pend, null);
	}
	
	/**
	 * Sends the message as an SMS
	 * @param serverMsg String in the format [<phone number>]<message>|
	 */
	private void sendSMS(String serverMsg)
	{
		String number;
		String message;
		
		try
		{
		number = serverMsg.substring(serverMsg.indexOf('[') + 1, serverMsg.indexOf(']'));
		message = serverMsg.substring(serverMsg.indexOf(']') + 1, serverMsg.indexOf('|'));
		sendSMS(number, message);
		} catch (Exception e)
		{
		}
	}
	
	/**
	 * Updates connection information
	 */
	@Override
	protected void onProgressUpdate(Integer... progress)
	{
		switch (progress[0])
		{
		case NO_CONNECTION:
			parentActivity.notConnected();
			break;
		case CONNECTION:
			parentActivity.connected();
			break;
		default:
			break;
		}
	}
	
	/**
	 * Called automatically whenn the thread ends
	 */
	@Override
	protected void onPostExecute(Boolean bool)
	{
		parentActivity.disconnected();
	}

}