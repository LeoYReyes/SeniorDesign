package sot.seniordesign.smsgateway;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.UnsupportedEncodingException;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.net.SocketAddress;
import java.util.regex.Pattern;

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
	public static final int SMS_PROCESSED = 2;
	public static final int SMS_SENT = 3;
	public static final int TIMEOUT = 10000;
	
	private SMSActivity parentActivity;
	private String ip = "";
	private String port = "";
	private volatile boolean keepAlive;
	private String msgFromServer;
	private long timeout = 30000;
	private long pingFreq = 9000;
	private String pingMsg = "|";

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
		Socket tcp = null;
		BufferedReader fromServer = null;
		OutputStream toServer = null;
		long lastRecvMsg = 0;
		long lastPing = 0;
		boolean receivedPing = true;
		
		keepAlive = true;
		while (keepAlive)
		{
			// if not connected try to get a tcp connection
			if (!receivedPing || tcp == null || !tcp.isConnected())
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
					toServer = tcp.getOutputStream();
					publishProgress(CONNECTION);
					receivedPing = true;
					lastRecvMsg = System.currentTimeMillis();
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
			
			//while connected
			while (receivedPing && keepAlive && tcp != null)
			{
				// read SMS message and forward to server
				try 
				{
					//send received SMSs to server
					while (SMSReceiver.hasMsg())
					{
						SmsMessage msg = SMSReceiver.getNextMsg();
						String messageBody = msg.getMessageBody().toString();
						messageBody = messageBody.replace("|", "");
						messageBody = messageBody.replace("[", "");
						messageBody = messageBody.replace("]", "");
						String output = "[" + msg.getOriginatingAddress() + "]" + messageBody + "|";
						toServer.write(output.getBytes("UTF-8"));
						publishProgress(SMS_PROCESSED);
					}
				} catch (IOException e) {
					// TODO Auto-generated catch block
					e.printStackTrace();
				}
				
				// Read TCP messages and forward via SMS
				try
				{
					boolean msgRecv = false;
					char[] next = new char[1];
					
					if (fromServer.ready()){
						msgRecv = true;
					}
					
					while (fromServer.ready())
					{
						fromServer.read(next, 0, 1);
						msgFromServer += next[0];
						if (next[0] == '|')
						{
							sendSMS(msgFromServer);
							msgFromServer = "";
						}
					}
					
					if (msgRecv)
					{
						lastRecvMsg = System.currentTimeMillis();
					}
					
				} catch (Exception e)
				{
				}
				
				// check heartbeat (ping)
				// sends a ping to the server and expects a response to check
				// if connection is still alive
				try 
				{
					long time = System.currentTimeMillis();
					if (time - lastRecvMsg > timeout)
					{
						receivedPing = false;
					}
					
					if (time - lastPing > pingFreq)
					{
						try
						{
							lastPing = time;
							toServer.write(pingMsg.getBytes("UTF-8"));
						} catch (NullPointerException e)
						{
							
						}
					}
					Thread.sleep(1000);
				} catch (InterruptedException e) {
					// TODO Auto-generated catch block
					e.printStackTrace();
				} catch (UnsupportedEncodingException e) {
					// TODO Auto-generated catch block
					e.printStackTrace();
				} catch (IOException e) {
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
		//PendingIntent pend = PendingIntent.getActivity(parentActivity, 0, new Intent(parentActivity, SMSActivity.class), 0);
		SmsManager sms = SmsManager.getDefault();
		sms.sendTextMessage(phoneNum, null, msg, null/*pend*/, null);
	}
	
	/**
	 * Sends the message as an SMS
	 * @param serverMsg String in the format [<phone number>]<message>|
	 */
	private void sendSMS(String serverMsg)
	{
		String number;
		String message;
		
		if (serverMsg.length() < 10)
		{
			return;
		}
		
		try
		{
		number = serverMsg.substring(serverMsg.indexOf('[') + 1, serverMsg.indexOf(']'));
		message = serverMsg.substring(serverMsg.indexOf(']') + 1, serverMsg.indexOf('|'));
		if (!number.matches("[0-9]+")){
		     return; 
		}
		sendSMS(number, message);
		publishProgress(SMS_SENT);
		} catch (Exception e)
		{
		}
	}
	
	/**
	 * Updates connection information on the GUI
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
		case SMS_PROCESSED:
			parentActivity.incSMSProcessed();
			break;
		case SMS_SENT:
			parentActivity.incSMSSent();
			break;
		default:
			break;
		}
	}
	
	/**
	 * Called automatically when the thread ends
	 */
	@Override
	protected void onPostExecute(Boolean bool)
	{
		parentActivity.disconnected();
	}

}