using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Drawing;
using System.Windows.Forms;

namespace WTKL
{
    class SystemTrayKeylogger : Form
    {
        private NotifyIcon trayIcon;
        private ContextMenu trayMenu;

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

        protected override void OnLoad(EventArgs e)
        {
            this.Visible = false;
            this.ShowInTaskbar = false;

            base.OnLoad(e);
        }

        private void OnExit(object sender, EventArgs e)
        {
            Application.Exit();
        }

        protected override void Dispose(bool isDisposing)
        {
            if (isDisposing)
            {
                trayIcon.Dispose();
            }
            base.Dispose(isDisposing);
        }

        [STAThread]
        public static void Main()
        {
            Application.Run(new SystemTrayKeylogger());
        }
    }
}
