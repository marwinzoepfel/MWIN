package gui

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/chzyer/readline"

	"client/network" // Anpassen, wenn der Pfad anders ist
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
	// Terminalgröße neu ermitteln
	width, height = getTerminalSize()

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

	fmt.Printf("\r\033[%d;1H> ", height+1) // \r für Cursor-Rücklauf
}

func getTerminalSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return 80, 24 // Default-Werte
	}
	wh := strings.Split(string(out), " ")
	width, height := 80, 24 // Default-Werte
	if len(wh) == 2 {
		fmt.Sscan(wh[1], &width)
		fmt.Sscan(wh[0], &height)
	}
	return width, height
}

func RunChat(conn net.Conn, reader *bufio.Reader) {
	width, height = getTerminalSize()
	updateScreen()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:              "> ",
		ForceUseInteractive: true,
	})
	if err != nil {
		fmt.Println("Error creating readline:", err)
		return
	}
	defer rl.Close()

	// Signal-Handler für SIGWINCH (Fenstergrößenänderung)
	sigwinchCh := make(chan os.Signal, 1)
	signal.Notify(sigwinchCh, syscall.SIGWINCH)

	go func() {
		for range sigwinchCh {
			updateScreen() // Bildschirm bei Größenänderung aktualisieren
		}
	}()

	// Goroutine zum Empfangen von Nachrichten
	go network.ReceiveMessages(conn, AddMessage)

	// Schleife zum Senden von Nachrichten und Aktualisieren der Anzeige
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		network.SendMessage(conn, line, AddMessage)
		rl.SetPrompt("> ")
		updateScreen()
	}
}
