package notification

import (
    "fmt"
    "os/exec"
    "regexp"
    "runtime"
)

// SendMessageToOS sends a notification message to the operating system
func SendMessageToOS(message string) {
    switch runtime.GOOS {
    case "darwin":
        Notification_darwin_nativ(message)
    default:
        // Other platforms: No notification (or optional implementation)
        fmt.Println("Notifications are not supported on this platform.")
    }
}

// macOS: Native notification
func Notification_darwin_nativ(message string) {
    // Regular expression to detect ANSI color codes
    re := regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)

    // Remove color codes from the message
    cleanMessage := re.ReplaceAllString(message, "")

    script := fmt.Sprintf(`display notification "%s" with title "MWIN Chat"`, cleanMessage)
    cmd := exec.Command("osascript", "-e", script)
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error sending notification:", err)
    }
}
