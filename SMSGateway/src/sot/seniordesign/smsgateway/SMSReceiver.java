package sot.seniordesign.smsgateway;

import java.util.LinkedList;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.telephony.SmsMessage;

/**
 * 
 * @author Chuck
 * Receives SMS and stores them to be accessed later.
 */
public class SMSReceiver extends BroadcastReceiver {
	
	private static volatile LinkedList<SmsMessage> msgList = null;
	private static final String SMS_RECEIVED = "android.provider.Telephony.SMS_RECEIVED";

	/**
	 * Receives SMS and places them in a linked list
	 */
	@Override
	public void onReceive(Context arg0, Intent arg1) {
		if (arg1.getAction().equals(SMS_RECEIVED))
		{
			if (msgList == null)
			{
				msgList = new LinkedList<SmsMessage>();
			}
			Bundle bundle = arg1.getExtras();
			
			if (bundle != null)
			{
				Object[] pdus = (Object[]) bundle.get("pdus");
				for (int i = 0; i < pdus.length; i++)
				{
					msgList.addLast(SmsMessage.createFromPdu((byte[])pdus[i]));
				}
			}
		}
	}
	
	/**
	 * 
	 * @return True if there are any SmsMessage available
	 */
	public static boolean hasMsg()
	{
		if (msgList == null)
		{
			return false;
		}
		return (msgList.size() > 0);
	}
	
	/**
	 * 
	 * @return Oldest SmsMessage available or null if there are none
	 */
	public static SmsMessage getNextMsg()
	{
		if (hasMsg())
		{
			return msgList.removeFirst();
		}
		else
		{
			return null;
		}
	}

}
