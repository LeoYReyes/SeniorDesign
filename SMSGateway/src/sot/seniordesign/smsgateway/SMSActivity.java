package sot.seniordesign.smsgateway;

import sot.seniordesign.smsgateway.R;

import android.os.Bundle;
import android.app.Activity;
import android.view.Menu;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

/**
 * 
 * @author Chuck
 *
 * Activity that defines the main (only) screen. Contains text fields
 * for the user to enter server info and a button to connect. Connects
 * to the server by TCP and sends received TCP messages out as SMS and
 * sends inbound SMS over the TCP.
 * 
 * The format for both inbound and outbound messages is:
 * [<phonenumber>]<message>|
 * Example: [1234567890]Hello World|
 * 
 * note: do not include '[', ']', or '|' aside from delimiting phone numbers
 */
public class SMSActivity extends Activity 
{
	
	// true when the application is connecting or connected to a server
	private boolean connection = false;
	private TCPAsyncTask tcp;

	/**
	 * auto-generated
	 */
	@Override
	protected void onCreate(Bundle savedInstanceState) 
	{
		super.onCreate(savedInstanceState);
		setContentView(R.layout.activity_sms);
	}

	/**
	 * auto-generated
	 */
	@Override
	public boolean onCreateOptionsMenu(Menu menu) 
	{
		// Inflate the menu; this adds items to the action bar if it is present.
		getMenuInflater().inflate(R.menu.sm, menu);
		return true;
	}
	
	/**
	 * Either starts a TCP connection and SMS sending/receiving, or if already
	 * connected, disconnects from TCP and stops SMS sending/receiving
	 * @param view
	 */
	public void connectFunction(View view)
	{
		if (!connection)
		{	
			connect();
		}
		else
		{
			disconnect();
		}
	}
	
	/**
	 * Makes the text fields uneditable and changes the button text to a disconnect
	 * prompt. Starts an AsyncTask (thread) which handles TCP and SMS messaging.
	 */
	public void connect()
	{
		EditText ipText = (EditText) findViewById(R.id.hostIPText);
		EditText portText = (EditText) findViewById(R.id.hostPortText);
		Button button = (Button) findViewById(R.id.connectionButton);
		TextView connectionText = (TextView) findViewById(R.id.connectionStatus);
		
		connection = true;
		ipText.setEnabled(false);
		portText.setEnabled(false);
		button.setEnabled(true);
		button.setText(getResources().getString(R.string.button_disconnect));
		connectionText.setText(getResources().getString(R.string.connection_connecting));
		
		tcp = new TCPAsyncTask(this);
		tcp.execute(ipText.getText().toString(), portText.getText().toString());
	}
	
	/**
	 * Sets the AsyncTask handling TCP and SMS communications to expire
	 */
	public void disconnect()
	{
		tcp.endTask();
		
		EditText ipText = (EditText) findViewById(R.id.hostIPText);
		EditText portText = (EditText) findViewById(R.id.hostPortText);
		Button button = (Button) findViewById(R.id.connectionButton);
		TextView connectionText = (TextView) findViewById(R.id.connectionStatus);
		
		connection = false;
		connectionText.setText(getResources().getString(R.string.connection_disconnecting));
		button.setText(getResources().getString(R.string.button_disconnect));
		ipText.setEnabled(false);
		portText.setEnabled(false);
		button.setEnabled(false);
	}
	
	/**
	 * Makes the text fields editable and re-enables the button (if it
	 * is not already). Call when the AsyncTask handling TCP and SMS ends.
	 */
	public void disconnected()
	{
		EditText ipText = (EditText) findViewById(R.id.hostIPText);
		EditText portText = (EditText) findViewById(R.id.hostPortText);
		Button button = (Button) findViewById(R.id.connectionButton);
		TextView connectionText = (TextView) findViewById(R.id.connectionStatus);
		
		connection = false;
		connectionText.setText(getResources().getString(R.string.connection_hide));
		button.setText(getResources().getString(R.string.button_connect));
		ipText.setEnabled(true);
		portText.setEnabled(true);
		button.setEnabled(true);
	}
	
	/**
	 * Sets text saying the app is not connected to server
	 */
	public void notConnected()
	{
		TextView connectionText = (TextView) findViewById(R.id.connectionStatus);
		connectionText.setText(getResources().getString(R.string.connection_connecting));
	}
	
	/**
	 * Sets text saying the app is connected to server
	 */
	public void connected()
	{
		TextView connectionText = (TextView) findViewById(R.id.connectionStatus);
		connectionText.setText(getResources().getString(R.string.connection_connected));
	}

}