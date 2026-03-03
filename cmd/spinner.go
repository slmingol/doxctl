package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type spinnerModel struct {
	spinner spinner.Model
	message string
	done    bool
	err     error
}

type doneMsg struct {
	err error
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil
	case doneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m spinnerModel) View() string {
	if m.done {
		if m.err != nil {
			return fmt.Sprintf("❌ %s: %v\n", m.message, m.err)
		}
		return fmt.Sprintf("✓ %s\n", m.message)
	}
	return fmt.Sprintf("%s %s...\n", m.spinner.View(), m.message)
}

// RunWithSpinner runs a function with a spinner showing progress
// Falls back to simple progress messages if TTY is not available
func RunWithSpinner(message string, fn func() error) error {
	if outputFormat != "table" {
		// Don't show spinner for JSON/YAML output
		return fn()
	}

	// ASCII spinner frames - bouncing ball
	frames := []string{"●     ", " ●    ", "  ●   ", "   ●  ", "    ● ", "     ●", "    ● ", "   ●  ", "  ●   ", " ●    "}

	// Track whether to stop the spinner
	var wg sync.WaitGroup
	done := make(chan bool)

	// Start the spinner animation
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Hide cursor
		fmt.Fprintf(os.Stderr, "\033[?25l")
		i := 0
		for {
			select {
			case <-done:
				// Clear the spinner line and show cursor
				fmt.Fprintf(os.Stderr, "\r\033[K\033[?25h")
				return
			default:
				// Print spinner frame
				fmt.Fprintf(os.Stderr, "\r%s %s", message, frames[i%len(frames)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Run the function
	err := fn()

	// Stop the spinner
	close(done)
	wg.Wait()

	// Show completion message
	if err != nil {
		fmt.Fprintf(os.Stderr, "\r\033[K❌ %s: %v\n", message, err)
	} else {
		// Ocean theme teal divider
		fmt.Fprintf(os.Stderr, "\r\033[K\033[38;2;0;128;128m%s\033[0m\n", strings.Repeat("─", 80))
		fmt.Fprintf(os.Stderr, "✓ %s\n", message)
		fmt.Fprintf(os.Stderr, "\033[38;2;0;128;128m%s\033[0m\n", strings.Repeat("─", 80))
	}

	return err
}
