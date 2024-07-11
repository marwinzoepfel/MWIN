package gui

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/chzyer/readline"

	"client/network" // Anpassen, wenn der Pfad anders ist
)

const (
	colorReset = "\033[0m"
	colorCyan  = "\033[36m"
)

var (
	messages []string
	mutex    sync.Mutex
	width    int
	height   int // Verfügbarer Platz für Nachrichten (ohne Eingabezeile)
)

func AddMessage(message string) {
	mutex.Lock()
	messages = append(messages, message)
	mutex.Unlock()
	updateScreen()
}

func updateScreen() {
	fmt.Print("\033[H\033[2J") // Bildschirm löschen

	mutex.Lock()
	// Anzahl der anzuzeigenden Nachrichten berechnen
	numMessagesToShow := height
	if len(messages) < height {
		numMessagesToShow = len(messages)
	}
	startIdx := len(messages) - numMessagesToShow

	// Nur die relevanten Nachrichten anzeigen
	for i := startIdx; i < len(messages); i++ {
		fmt.Println(messages[i])
	}
	mutex.Unlock()

	// Eingabezeile am unteren Rand anzeigen
	fmt.Printf("\r\033[%d;1H> ", height+1) // \r für Cursor-Rücklauf
}

func getTerminalSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return 0, 0
	}
	wh := strings.Split(string(out), " ")
	width, height := 80, 24 // Default-Werte, falls die Ermittlung fehlschlägt
	if len(wh) == 2 {
		fmt.Sscan(wh[1], &width)
		fmt.Sscan(wh[0], &height)
	}

	// Bildschirm löschen und Nachrichten anzeigen
	fmt.Print("\033[H\033[2J") // Bildschirm löschen
	return width, height
}

func RunChat(conn net.Conn, reader *bufio.Reader) {
	// Terminalgröße ermitteln, bevor Nachrichten gesendet oder empfangen werden
	width, height = getTerminalSize()

	// Bildschirm aktualisieren, um die korrekte Größe zu verwenden
	updateScreen()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:              "> ", // Prompt wird jetzt immer angezeigt
		ForceUseInteractive: true,
	})
	if err != nil {
		fmt.Println("Error creating readline:", err)
		return
	}
	defer rl.Close()

	// Goroutine zum Empfangen von Nachrichten
	go network.ReceiveMessages(conn, AddMessage)

	// Schleife zum Senden von Nachrichten und Aktualisieren der Anzeige
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		network.SendMessage(conn, line, AddMessage)
		// Bildschirm nach dem Senden aktualisieren, aber Prompt neu setzen
		rl.SetPrompt("> ")
		updateScreen()
	}
}
