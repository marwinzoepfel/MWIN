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

	"client/network" // Adjust if the path is different
)

var (
	messages []string
	mutex    sync.Mutex
	width    int
    height   int // Available space for messages (excluding input line)
)

// AddMessage adds a message to the message list and updates the screen

func AddMessage(message string) {
	mutex.Lock()
	messages = append(messages, message)
	mutex.Unlock()
	updateScreen()
}

// updateScreen clears the screen and displays the most recent messages
func updateScreen() {
    	// Get new terminal size
	width, height = getTerminalSize()

	fmt.Print("\033[H\033[2J") // Bildschirm löschen

	mutex.Lock()
	// Calculate number of messages to display
	numMessagesToShow := height
	if len(messages) < height {
		numMessagesToShow = len(messages)
	}
	startIdx := len(messages) - numMessagesToShow

	// Display only the relevant messages
	for i := startIdx; i < len(messages); i++ {
		fmt.Println(messages[i])
	}
	mutex.Unlock()

	fmt.Printf("\r\033[%d;1H> ", height+1) // \r for cursor return
}

func getTerminalSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return 80, 24 // Default values
	}
	wh := strings.Split(string(out), " ")
	width, height := 80, 24 // Default values
	if len(wh) == 2 {
		fmt.Sscan(wh[1], &width)
		fmt.Sscan(wh[0], &height)
	}
	return width, height
}

// RunChat initiates the chat interface, handling message sending/receiving and screen updates
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
			updateScreen() // Update screen on resize
		}
	}()

    	// Goroutine to receive messages
	go network.ReceiveMessages(conn, AddMessage)

	// Loop to send messages and update the display
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
