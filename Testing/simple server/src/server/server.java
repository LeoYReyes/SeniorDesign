package server;

import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.Arrays;
import java.util.Random;

/**
 * simple server that will accept 1 tcp connection and print some basic info.
 */
public class server {
	
	private static ServerSocket ss;
	private static Socket connection;
	private static BufferedReader fromClient;
	private static DataOutputStream toClient;
	
	/**
	 * Opens a server socket on port 10011 and accepts the first connection
	 */
	private static void connect()
	{
		try {
			ServerSocket ss = new ServerSocket(10011);
			Socket connection = ss.accept();
			BufferedReader fromClient = new BufferedReader(new InputStreamReader(connection.getInputStream()));
			DataOutputStream toClient = new DataOutputStream(connection.getOutputStream());
			System.out.println("Connected");
		} catch (Exception e)
		{}
	}
	
	/**
	 * Closes any connections and server sockets
	 */
	private static void disconnect()
	{
		try
		{
			toClient.close();
			fromClient.close();
			connection.close();
			ss.close();
		} catch (Exception e)
		{}
	}
	
	private static void testConnect()
	{
		System.out.println("Testing connection and MAC address sending...");
		connect();
		try {
			System.out.println(fromClient.readLine());
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		disconnect();
	}
	
	public static void main(String args[]) throws Exception {
		testConnect();
		/*
		ServerSocket ss = new ServerSocket(10011);
		
		while(true) {
			byte[] buffer = new byte[256];
			
			boolean stolen = false;
			Random rand = new Random();
			rand.setSeed(System.currentTimeMillis());
			
			stolen = rand.nextBoolean();
			
			System.out.println("\n=========================\n\nWaiting for tcp connection...");
			Socket connection = ss.accept();
			BufferedReader fromClient = new BufferedReader(new InputStreamReader(connection.getInputStream()));
			DataOutputStream toClient = new DataOutputStream(connection.getOutputStream());
			System.out.println("Connection from: " + connection.getInetAddress().toString());
			System.out.println("By machine: " + fromClient.readLine());
			Thread.sleep(5000);
			if (stolen)
			{
				System.out.println("Reporting stolen");
				toClient.writeByte(5);
			}
			else
			{
				System.out.println("Reporting not stolen");
				toClient.writeByte(4);
				continue;
			}
			Thread.sleep(5000);
			System.out.println("Turning keylogger on...");
			toClient.writeByte(0);
			Thread.sleep(15000);
			System.out.println("Turning keylogger off...");
			toClient.writeByte(1);
			Thread.sleep(5000);
			System.out.println("Requesting keylog...");
			toClient.writeByte(3);
			System.out.println("Keyloggy stuff: " + fromClient.readLine());
			Thread.sleep(1000);
			while (fromClient.ready()) {
				System.out.println("Keyloggy stuff: " + fromClient.readLine());
				Thread.sleep(1000);
			}
			Thread.sleep(1000);
			System.out.println("Requesting trace route...");
			toClient.writeByte(2);
			System.out.println(fromClient.readLine());
			Thread.sleep(7000);
			connection.close();
			//udp.close();
			
			//Thread.sleep(6000);
			//connection.close();
		}*/
	}
	
}
