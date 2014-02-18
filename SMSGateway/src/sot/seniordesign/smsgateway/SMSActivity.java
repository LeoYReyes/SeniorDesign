package sot.seniordesign.smsgateway;

import sot.seniordesign.smsgateway.R;

import android.os.Bundle;
import android.app.Activity;
import android.view.Menu;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

public class SMSActivity extends Activity {
	
	// true when the application is connecting or connected to a server
	private boolean connection = false;
	private TCPAsyncTask tcp;

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.activity_sms);
	}

	@Override
	public boolean onCreateOptionsMenu(Menu menu) {
		// Inflate the menu; this adds items to the action bar if it is present.
		getMenuInflater().inflate(R.menu.sm, menu);
		return true;
	}
	
	public void connectFunction(View view)
	{
		if (!connection)
		{	
			connect();
		}
		else
		{
			tcp.endTask();
		}
	}
	
	public synchronized void connect()
	{
		EditText ipText = (EditText) findViewById(R.id.hostIPText);
		EditText portText = (EditText) findViewById(R.id.hostPortText);
		Button button = (Button) findViewById(R.id.connectionButton);
		
		connection = true;
		ipText.setEnabled(false);
		portText.setEnabled(false);
		button.setEnabled(false);
		
		tcp = new TCPAsyncTask(this);
		tcp.execute(ipText.getText().toString(), portText.getText().toString());
	}
	
	public synchronized void disconnect()
	{
		EditText ipText = (EditText) findViewById(R.id.hostIPText);
		EditText portText = (EditText) findViewById(R.id.hostPortText);
		Button button = (Button) findViewById(R.id.connectionButton);
		
		connection = false;
		button.setText(getResources().getString(R.string.button_connect));
		ipText.setEnabled(true);
		portText.setEnabled(true);
		button.setEnabled(true);
	}
	
	public synchronized void connected()
	{
		Button button = (Button) findViewById(R.id.connectionButton);
		button.setText(getResources().getString(R.string.button_disconnect));
		button.setEnabled(true);
	}

}