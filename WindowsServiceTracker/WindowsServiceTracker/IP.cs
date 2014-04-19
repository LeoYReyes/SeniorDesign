using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Net.NetworkInformation;
using System.Net;

namespace WindowsServiceTracker
{

    /* Contains the function to do a trace route to a given address.
     */
    class IP
    {

        private const string Data = "Ping trace route";

        /* Returns a list of all nodes, by IP, packets travel through between this machine and
         * a target destination. List of addresses is ordered in the same order that a packet
         * would travel through them from this machine to the target.
         */
        public static IEnumerable<IPAddress> getTraceRoute(string hostNameOrAddress)
        {
            return getTraceRoute(hostNameOrAddress, 1, 3);
        }

        /* Workhorse of the getTraceRoute function. Recursively pings the target machine with an
         * increasing time to live until it is reached and returns the list of IPs of all nodes
         * traversed.
         */
        private static IEnumerable<IPAddress> getTraceRoute(string hostNameOrAddress, int ttl, int timeouts)
        {
            Ping pinger = new Ping();
            PingOptions pingerOptions = new PingOptions(ttl, true);
            int timeout = 10000;
            byte[] buffer = Encoding.ASCII.GetBytes(Data);
            PingReply reply = default(PingReply);
            reply = pinger.Send(hostNameOrAddress, timeout, buffer, pingerOptions);
            List<IPAddress> result = new List<IPAddress>();
            
            if (reply.Status == IPStatus.Success)
            {
                result.Add(reply.Address);
            }
            else if (reply.Status == IPStatus.TtlExpired)
            {
                result.Add(reply.Address);
                IEnumerable<IPAddress> tempResult = default(IEnumerable<IPAddress>);
                tempResult = getTraceRoute(hostNameOrAddress, ttl + 1, timeouts);
                result.AddRange(tempResult);
            }
            else
            {
                if (timeouts > 0)
                {
                    IEnumerable<IPAddress> tempResult = default(IEnumerable<IPAddress>);
                    tempResult = getTraceRoute(hostNameOrAddress, ttl + 1, timeouts - 1);
                    result.AddRange(tempResult);
                }
                //Console.WriteLine("Failed");
            }

            return result;
        }
    }
}