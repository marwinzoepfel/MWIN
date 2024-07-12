package notificatioin

import (
	"fmt"
	"os/exec"
	"runtime"
)

func SendMessageToOS(message string) {
	switch runtime.GOOS {
	case "darwin":
		Notification_darwin_nativ(message)
	default:
		// Andere Plattformen: Keine Benachrichtigung (oder optionale Implementierung)
		fmt.Println("Benachrichtigungen werden auf dieser Plattform nicht unterstützt.")
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
