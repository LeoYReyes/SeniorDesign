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
        //static IPAddress ipAddress = IPAddress.Parse("71.91.88.205");
        static IPEndPoint remoteEP = new IPEndPoint(ipAddress, 10015);
        static byte[] bytes = new byte[1024];

        // Create a TCP/IP  socket.
        //Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);

        /*
         * laptopDevice connects to server and sends MAC Address. It then waits for the server
         * to respond with a stolen/not stolen OP code.
         */
        static void testValidIDNotStolen()
        {
            Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            sender.Connect(remoteEP);
            Console.WriteLine("Socket connected to {0}", sender.RemoteEndPoint.ToString());

            byte[] msg = Encoding.ASCII.GetBytes("0123456789AB\n");

            // Send the data through the socket.
            int bytesSent = sender.Send(msg);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
            int stolenCode = buffer[0];
            Console.WriteLine("Received stolen code: {0}", stolenCode);
        }

        static void testValidIDStolen01()
        {
            Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            sender.Connect(remoteEP);
            Console.WriteLine("Socket connected to {0}", sender.RemoteEndPoint.ToString());

            byte[] msg = Encoding.ASCII.GetBytes("BA9876543210\n");

            // Send the data through the socket.
            int bytesSent = sender.Send(msg);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
            int stolenCode = buffer[0];
            Console.WriteLine("Received stolen code: {0}", stolenCode);

            byte[] buffer1 = new byte[1000];
            byte[] buffer2 = new byte[1000];
            sender.Receive(buffer1);
            sender.Receive(buffer2);
            byte[] msg2 = Encoding.ASCII.GetBytes("0.0.0.0~125.32.192.13~124.234.134.54~145.3.21.94~255.255.255.255\n");
            byte[] msg2WithOpcode = new byte[msg2.Length + 1];
            Array.Copy(msg2, 0, msg2WithOpcode, 1, msg2.Length);
            msg2WithOpcode[0] = 0x83;
            bytesSent = sender.Send(msg2WithOpcode);
            Console.WriteLine("Sent IP traceroute data");

            Thread.Sleep(100);
            byte[] msg3 = Encoding.ASCII.GetBytes("`1234567890-=~!@#$%^&*()_+qwertyuiop[]\\QWERTYUIOP{}|asdfghjkl;ASDFGHJKL:\"zxcvbnm,./ZXCVBNM<>?\n");
            byte[] msg3WithOpcode = new byte[msg3.Length + 1];
            Array.Copy(msg3, 0, msg3WithOpcode, 1, msg3.Length);
            msg3WithOpcode[0] = 0x82;
            bytesSent = sender.Send(msg3WithOpcode);
            Console.WriteLine("Sent keylog data");
            sender.Close();
        }

        static void testInvalidID01()
        {
            Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            sender.Connect(remoteEP);
            Console.WriteLine("Socket connected to {0}", sender.RemoteEndPoint.ToString());

            byte[] msg = Encoding.ASCII.GetBytes("1234\n");
            int bytesSent = sender.Send(msg);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
        }

        static void testInvalidID02()
        {
            Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            sender.Connect(remoteEP);
            Console.WriteLine("Socket connected to {0}", sender.RemoteEndPoint.ToString());

            byte[] msg = Encoding.ASCII.GetBytes("123474");
            int bytesSent = sender.Send(msg);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
        }

        /*static void test3()
        {
            Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            sender.Connect(remoteEP);
            Console.WriteLine("Socket connected to {0}", sender.RemoteEndPoint.ToString());

            byte[] msg = Encoding.ASCII.GetBytes("BA9876543210\n");

            // Send the data through the socket.
            int bytesSent = sender.Send(msg);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
            int stolenCode = buffer[0];
            Console.WriteLine("Received stolen code: {0}", stolenCode);

            msg = Encoding.ASCII.GetBytes("0.0.0.0~125.32.192.13~124.234.134.54~145.3.21.94~255.255.255.255\n");
            byte[] msgWithOpcode = new byte[msg.Length + 1];
            Array.Copy(msg, 0, msgWithOpcode, 1, msg.Length);
            msgWithOpcode[0] = 0x83;
            bytesSent = sender.Send(msgWithOpcode);
            Console.WriteLine("Sent IP traceroute data");
            sender.Close();
        }*/

        /*static void test4()
        {
            Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            sender.Connect(remoteEP);
            Console.WriteLine("Socket connected to {0}", sender.RemoteEndPoint.ToString());

            byte[] msg1 = Encoding.ASCII.GetBytes("112233445566\n");

            // Send the data through the socket.
            int bytesSent = sender.Send(msg1);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
            int stolenCode = buffer[0];
            Console.WriteLine("Received stolen code: {0}", stolenCode);

            byte[] msg2 = Encoding.ASCII.GetBytes("125.32.192.13~124.234.134.54~145.3.21.94\n");
            byte[] msgWithOpcode1 = new byte[msg2.Length + 1];
            Array.Copy(msg2, 0, msgWithOpcode1, 1, msg2.Length);
            msgWithOpcode1[0] = 0x83;
            bytesSent = sender.Send(msgWithOpcode1);
            Console.WriteLine("Sent IP traceroute data");

            byte[] msg3 = Encoding.ASCII.GetBytes("aabbccddeeffgghhiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz1234567890\n");
            byte[] msgWithOpcode2 = new byte[msg3.Length + 1];
            Array.Copy(msg3, 0, msgWithOpcode2, 1, msg3.Length);
            msgWithOpcode2[0] = 0x82;
            bytesSent = sender.Send(msgWithOpcode2);
            Console.WriteLine("Sent keylog data");

            byte[] msg4 = Encoding.ASCII.GetBytes("more keylog data sdlakfjporieuwqrlkjs;ldkfj\n");
            byte[] msgWithOpcode3 = new byte[msg4.Length + 1];
            Array.Copy(msg4, 0, msgWithOpcode3, 1, msg4.Length);
            msgWithOpcode3[0] = 0x82;
            bytesSent = sender.Send(msgWithOpcode3);
            Console.WriteLine("Sent keylog data");
            sender.Close();
        }*/

        static void Main(string[] args)
        {
            testValidIDNotStolen();
            testValidIDStolen01();
            testInvalidID01();
            testInvalidID02();
            //test4();
        }
    }
}
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
        //static IPAddress ipAddress = IPAddress.Parse("71.91.88.205");
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

            byte[] msg = Encoding.ASCII.GetBytes("0123456789AB\n");

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

            byte[] msg = Encoding.ASCII.GetBytes("BA9876543210\n");

            // Send the data through the socket.
            int bytesSent = sender.Send(msg);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
            int stolenCode = buffer[0];
            Console.WriteLine("Received stolen code: {0}", stolenCode);

            byte[] buffer1 = new byte[1000];
            byte[] buffer2 = new byte[1000];
            sender.Receive(buffer1);
            sender.Receive(buffer2);
            byte[] msg2 = Encoding.ASCII.GetBytes("0.0.0.0~125.32.192.13~124.234.134.54~145.3.21.94~255.255.255.255\n");
            byte[] msg2WithOpcode = new byte[msg2.Length + 1];
            Array.Copy(msg2, 0, msg2WithOpcode, 1, msg2.Length);
            msg2WithOpcode[0] = 0x83;
            bytesSent = sender.Send(msg2WithOpcode);
            Console.WriteLine("Sent IP traceroute data");

            Thread.Sleep(100);
            byte[] msg3 = Encoding.ASCII.GetBytes("`1234567890-=~!@#$%^&*()_+qwertyuiop[]\\QWERTYUIOP{}|asdfghjkl;ASDFGHJKL:\"zxcvbnm,./ZXCVBNM<>?\n");
            byte[] msg3WithOpcode = new byte[msg3.Length + 1];
            Array.Copy(msg3, 0, msg3WithOpcode, 1, msg3.Length);
            msg3WithOpcode[0] = 0x82;
            bytesSent = sender.Send(msg3WithOpcode);
            Console.WriteLine("Sent keylog data");
            sender.Close();
        }

        /*static void test3()
        {
            Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            sender.Connect(remoteEP);
            Console.WriteLine("Socket connected to {0}", sender.RemoteEndPoint.ToString());

            byte[] msg = Encoding.ASCII.GetBytes("BA9876543210\n");

            // Send the data through the socket.
            int bytesSent = sender.Send(msg);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
            int stolenCode = buffer[0];
            Console.WriteLine("Received stolen code: {0}", stolenCode);

            msg = Encoding.ASCII.GetBytes("0.0.0.0~125.32.192.13~124.234.134.54~145.3.21.94~255.255.255.255\n");
            byte[] msgWithOpcode = new byte[msg.Length + 1];
            Array.Copy(msg, 0, msgWithOpcode, 1, msg.Length);
            msgWithOpcode[0] = 0x83;
            bytesSent = sender.Send(msgWithOpcode);
            Console.WriteLine("Sent IP traceroute data");
            sender.Close();
        }*/

        /*static void test4()
        {
            Socket sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            sender.Connect(remoteEP);
            Console.WriteLine("Socket connected to {0}", sender.RemoteEndPoint.ToString());

            byte[] msg1 = Encoding.ASCII.GetBytes("112233445566\n");

            // Send the data through the socket.
            int bytesSent = sender.Send(msg1);
            byte[] buffer = new byte[1000];
            sender.Receive(buffer);
            int stolenCode = buffer[0];
            Console.WriteLine("Received stolen code: {0}", stolenCode);

            byte[] msg2 = Encoding.ASCII.GetBytes("125.32.192.13~124.234.134.54~145.3.21.94\n");
            byte[] msgWithOpcode1 = new byte[msg2.Length + 1];
            Array.Copy(msg2, 0, msgWithOpcode1, 1, msg2.Length);
            msgWithOpcode1[0] = 0x83;
            bytesSent = sender.Send(msgWithOpcode1);
            Console.WriteLine("Sent IP traceroute data");

            byte[] msg3 = Encoding.ASCII.GetBytes("aabbccddeeffgghhiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz1234567890\n");
            byte[] msgWithOpcode2 = new byte[msg3.Length + 1];
            Array.Copy(msg3, 0, msgWithOpcode2, 1, msg3.Length);
            msgWithOpcode2[0] = 0x82;
            bytesSent = sender.Send(msgWithOpcode2);
            Console.WriteLine("Sent keylog data");

            byte[] msg4 = Encoding.ASCII.GetBytes("more keylog data sdlakfjporieuwqrlkjs;ldkfj\n");
            byte[] msgWithOpcode3 = new byte[msg4.Length + 1];
            Array.Copy(msg4, 0, msgWithOpcode3, 1, msg4.Length);
            msgWithOpcode3[0] = 0x82;
            bytesSent = sender.Send(msgWithOpcode3);
            Console.WriteLine("Sent keylog data");
            sender.Close();
        }*/

        static void Main(string[] args)
        {
            test1();
            test2();
            //test3();
            //test4();
        }
    }
}
