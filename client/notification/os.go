package notification

import (
	"fmt"
	"os/exec"
	"regexp"
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
	// Regulärer Ausdruck zum Erkennen von ANSI-Farbcodes
	re := regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)

	// Entferne Farbcodes aus der Nachricht
	cleanMessage := re.ReplaceAllString(message, "")

	script := fmt.Sprintf(`display notification "%s" with title "MWIN Chat"`, cleanMessage)
	cmd := exec.Command("osascript", "-e", script)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error sending notification:", err)
	}
}
