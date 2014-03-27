using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Net;
using System.Net.Sockets;
using System.Threading;

namespace SimpleTCPClient
{
    class TCPClient
    {
        static IPHostEntry ipHostInfo = Dns.GetHostEntry(Dns.GetHostName());
        static IPAddress ipAddress = new IPAddress(0x0100007F);
        static IPEndPoint remoteEP = new IPEndPoint(ipAddress, 10015);
        static byte[] bytes = new byte[1024];

        // Create a TCP/IP  socket.
        //Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);

        /*
         * laptopDevice connects to server and sends MAC Address. It then waits for the server
         * to respond with a stolen/not stolen OP code.
         */
        static void test1()
        {
            Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            sender.Connect(remoteEP);
            Console.WriteLine("Socket connected to {0}", sender.RemoteEndPoint.ToString());

            byte[] msg = Encoding.ASCII.GetBytes("0123456789AB");

            // Send the data through the socket.
            int bytesSent = sender.Send(msg);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
            int stolenCode = buffer[0];
            Console.WriteLine("Received stolen code: {0}", stolenCode);
        }

        static void test2()
        {
            Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            sender.Connect(remoteEP);
            Console.WriteLine("Socket connected to {0}", sender.RemoteEndPoint.ToString());

            byte[] msg = Encoding.ASCII.GetBytes("BA9876543210");

            // Send the data through the socket.
            int bytesSent = sender.Send(msg);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
            int stolenCode = buffer[0];
            Console.WriteLine("Received stolen code: {0}", stolenCode);

            msg = Encoding.ASCII.GetBytes("this is new keylog data");
            byte[] msgWithOpcode = new byte[msg.Length + 1];
            Array.Copy(msg, 0, msgWithOpcode, 1, msg.Length);
            msgWithOpcode[0] = 0x82;
            bytesSent = sender.Send(msgWithOpcode);
            Console.WriteLine("Sent keylog data");
        }

        static void test3()
        {
            Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            sender.Connect(remoteEP);
            Console.WriteLine("Socket connected to {0}", sender.RemoteEndPoint.ToString());

            byte[] msg = Encoding.ASCII.GetBytes("BA9876543210");

            // Send the data through the socket.
            int bytesSent = sender.Send(msg);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
            int stolenCode = buffer[0];
            Console.WriteLine("Received stolen code: {0}", stolenCode);

            msg = Encoding.ASCII.GetBytes("125.32.192.13~124.234.134.54~145.3.21.94");
            byte[] msgWithOpcode = new byte[msg.Length + 1];
            Array.Copy(msg, 0, msgWithOpcode, 1, msg.Length);
            msgWithOpcode[0] = 0x83;
            bytesSent = sender.Send(msgWithOpcode);
            Console.WriteLine("Sent IP traceroute data");
        }

        static void Main(string[] args)
        {
            test1();
            //test2();
            test3();
            //byte[] msg2 = Encoding.ASCII.GetBytes("2192.168.0.1~127.0.0.1~72.54.10.100\n");
            //byte[] msg2 = Encoding.ASCII.GetBytes("3test keylogger data\n");
            //int bytesSent2 = sender.Send(msg2);

            // Receive the response from the remote device.
            //int bytesRec = sender.Receive(bytes);
            //Console.WriteLine("Echoed test = {0}", Encoding.ASCII.GetString(bytes, 0, bytesRec));

            // Release the socket.
            //sender.Shutdown(SocketShutdown.Both);
            //sender.Close();
        }
    }
}
