package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	"github.com/godbus/dbus"
)

const (
	colorReset = "\033[0m"
	colorCyan  = "\033[36m"
)

// ReceiveMessages continuously reads and prints messages from the server connection.
func ReceiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection to server interrupted:", err)
			return // Exit the goroutine if there's an error
		}

		// Check if the message is empty (e.g., due to a disconnection)
		if len(message) == 0 {
			continue
		}

		// Handle potential errors when printing the message (e.g., if the terminal is closed)
		if _, err := fmt.Print(message); err != nil {
			fmt.Println("Error printing message:", err)
			return // Exit the goroutine if there's an error
		}
		fmt.Println("Nachricht für OS: " + message)
		sendMessageToOS(message)
		// notify.Notify("Benachrichtigung", message, "", "")
	}
}

func sendMessageToOS(message string) {
	switch runtime.GOOS {
	case "darwin":
		// Embed the icon

		var iconData []byte

		// macOS: Native Benachrichtigung
		// Write embedded icon to a temporary file
		tmpFile, err := os.CreateTemp(os.TempDir(), "mwin_icon_*.png")
		if err != nil {
			fmt.Println("Error creating temporary icon file:", err)
			return
		}
		defer os.Remove(tmpFile.Name()) // Remove the file when done

		_, err = tmpFile.Write(iconData)
		if err != nil {
			fmt.Println("Error writing icon to temporary file:", err)
			return
		}
		tmpFile.Close()

		// Construct AppleScript command (Alternative Approach)
		script := fmt.Sprintf(`
			 on run argv
				  display notification (item 1 of argv) with title "MWIN Chat" sound name "default"
				  tell application "System Events"
						set theIcon to POSIX file "%s"
						set the image of (first item of (get notifications)) to theIcon
				  end tell
			 end run
		`, tmpFile.Name())

		// Execute AppleScript (with message as argument)
		cmd := exec.Command("osascript", "-e", script, message)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error sending notification:", err)
			fmt.Println("AppleScript output:", string(output))
			return
		}

	case "linux":
		// Linux: D-Bus oder systray
		conn, err := dbus.SessionBus() // Hier wird conn für D-Bus benötigt
		if err != nil {
			systray.Run(onSystrayReady, onSystrayExit)
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
		systray.Run(onSystrayReady, onSystrayExit)
	}
}

func onSystrayReady() {
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

func onSystrayExit() {
	// Keine spezifischen Aufräumarbeiten erforderlich
}

// SendMessages continuously reads messages from the user's input and sends them to the server.
func SendMessages(conn net.Conn, reader *bufio.Reader) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		currentTime := time.Now().Format("15:04:05")

		// Print the formatted message with a cyan color for the timestamp and "You:"
		_, err := fmt.Printf(colorCyan+"[%s] You: %s\n"+colorReset, currentTime, message)
		if err != nil {
			fmt.Println("Error printing message:", err)
			return // Exit the function if there's an error
		}

		// Send only the message content, without the timestamp and "You:"
		_, err = fmt.Fprintln(conn, message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return // Exit the function if there's an error
		}
	}

	// Handle potential errors from the scanner
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from console:", err)
	}
}
