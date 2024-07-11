package notificatioin

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
	"github.com/godbus/dbus"
)

func SendMessageToOS(message string) {
	switch runtime.GOOS {
	case "darwin":
		Notification_darwin_nativ(message)

	case "linux":
		// Linux: D-Bus oder systray
		conn, err := dbus.SessionBus() // Hier wird conn für D-Bus benötigt
		if err != nil {
			systray.Run(OnSystrayReady, OnSystrayExit)
			return
		}
		obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
		call := obj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0),
			"", "Neue Nachricht", message, []string{}, map[string]dbus.Variant{}, int32(-1))
		if call.Err != nil {
			fmt.Println("Error sending notification:", call.Err)
		}

	case "windows":
		// Windows: Native Benachrichtigung oder systray
		// ... (Implementierung mit golang.org/x/sys/windows oder systray)

	default:
		// Andere Plattformen: systray
		systray.Run(OnSystrayReady, OnSystrayExit)
	}
}

// macOS: Native Benachrichtigung
func Notification_darwin_nativ(message string) {
	script := fmt.Sprintf(`display notification "%s" with title "MWIN Chat"`, message)
	cmd := exec.Command("osascript", "-e", script)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error sending notification:", err)
	}
}

func OnSystrayReady() {
	systray.SetTooltip("MWIN")
	mNewMessage := systray.AddMenuItem("Neue Nachricht", "")
	go func() {
		for {
			select {
			case msg := <-mNewMessage.ClickedCh:
				// Handle click on "New Message" menu item
				fmt.Println("New Message clicked:", msg)
			}
		}
	}()
}

func OnSystrayExit() {
	// Keine spezifischen Aufräumarbeiten erforderlich
}
