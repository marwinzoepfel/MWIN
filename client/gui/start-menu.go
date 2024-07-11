package gui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type model struct {
	width, height int
	buttons       []string // Start Menu items
	selected      int
	hovered       int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc": // Allow "esc" to exit like a Start Menu
			return m, tea.Quit
		case "up", "left":
			m.selected = max(0, m.selected-1)
			m.hovered = m.selected
		case "down", "right":
			m.selected = min(len(m.buttons)-1, m.selected+1)
			m.hovered = m.selected
		case "enter": // Simulate clicking the selected item
			return m, tea.Quit // Exit after selection
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			for i, button := range m.buttons {
				if m.isClickInButtonArea(msg.X, msg.Y, button, i) {
					m.selected = i
					return m, tea.Quit // Exit on click
				}
			}
		}
		for i, button := range m.buttons {
			if m.isClickInButtonArea(msg.X, msg.Y, button, i) {
				m.hovered = i
				break
			}
		}
	}
	return m, nil
}

func (m model) isClickInButtonArea(x, y int, buttonLabel string, buttonIndex int) bool {
	buttonWidth := len(buttonLabel) + 4 // Add some padding
	buttonHeight := 1
	buttonSpacing := 5 // Adjust spacing between buttons

	centerX := m.width / 2
	centerY := m.height / 2

	buttonLeft := centerX - ((buttonWidth+buttonSpacing)*len(m.buttons))/2 + (buttonWidth+buttonSpacing)*buttonIndex
	buttonRight := buttonLeft + buttonWidth
	buttonTop := centerY - buttonHeight/2
	buttonBottom := centerY + buttonHeight/2

	return x >= buttonLeft && x <= buttonRight && y >= buttonTop && y <= buttonBottom
}

func (m model) View() string {
	var buttons []string
	for i, label := range m.buttons {
		style := lipgloss.NewStyle().
			Width(15). // Maintain consistent button width
			Height(1).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("235")) // Light text on dark background for better contrast

		if i == m.selected {
			style = style.Background(lipgloss.Color("21")) // Darker blue for selected
		} else if i == m.hovered {
			style = style.Background(lipgloss.Color("27")) // Lighter blue for hovered
		}

		button := style.Render(wordwrap.String(label, 20))
		buttons = append(buttons, button)
	}

	// Calculate horizontal offset for centering (same as before)
	offsetX := 0

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().PaddingLeft(offsetX).Render( // Manual padding
			lipgloss.JoinHorizontal(lipgloss.Left, buttons...)),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("0")),
		lipgloss.WithWhitespaceChars(" "))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func StartMenu() string {
	p := tea.NewProgram(model{buttons: []string{"Auto Start", "Manual Start"}}) // Replace with your menu items

	// Start the program and wait for it to finish
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	// The final model is returned by p.Run()
	finalModel := m.(model)
	return finalModel.buttons[finalModel.selected] // Return the clicked button
}
