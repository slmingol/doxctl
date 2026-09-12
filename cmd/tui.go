package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type checkStatus int

const (
	statusPending checkStatus = iota
	statusRunning
	statusPass
	statusFail
	statusSkipped
)

type checkItem struct {
	id       string
	name     string
	args     []string
	selected bool
	status   checkStatus
	output   string
	errOut   string
	duration time.Duration
}

type tuiModel struct {
	checks   []checkItem
	cursor   int
	running  bool
	done     bool
	quitting bool
	cfgFile  string
	spinner  spinner.Model
	exported string
}

type checkStartMsg struct{ id string }
type checkDoneMsg struct {
	id       string
	passed   bool
	output   string
	errOut   string
	duration time.Duration
}

var (
	styleBold    = lipgloss.NewStyle().Bold(true)
	styleTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	styleDim     = lipgloss.NewStyle().Faint(true)
	styleGreen   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	styleRed     = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	styleYellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	styleCursor  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	styleChecked = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
)

func initialTUIModel() tuiModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styleYellow

	checks := []checkItem{
		{id: "dns-cfg", name: "DNS Resolver Config", args: []string{"dns", "--domName", "--domSearch", "--domAddr"}, selected: true},
		{id: "dns-ping", name: "DNS Resolver Ping", args: []string{"dns", "--ping"}, selected: true},
		{id: "vpn-iface", name: "VPN Interface", args: []string{"vpn", "--ifReachableChk"}, selected: true},
		{id: "vpn-routes", name: "VPN Routes", args: []string{"vpn", "--routeChk"}, selected: true},
		{id: "svrs", name: "Server Reachability", args: []string{"svrs", "--reach"}, selected: true},
		{id: "net-perf", name: "Network Performance", args: []string{"net", "--perf"}, selected: true},
	}

	return tuiModel{
		checks:  checks,
		spinner: s,
		cfgFile: cfgFile,
	}
}

func (m tuiModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.checks)-1 {
				m.cursor++
			}
		case " ":
			if !m.running {
				m.checks[m.cursor].selected = !m.checks[m.cursor].selected
			}
		case "enter", "r":
			if !m.running {
				m.running = true
				m.done = false
				m.exported = ""
				// reset statuses
				for i := range m.checks {
					if m.checks[i].selected {
						m.checks[i].status = statusPending
						m.checks[i].output = ""
						m.checks[i].errOut = ""
					} else {
						m.checks[i].status = statusSkipped
					}
				}
				return m, startAllSelected(m)
			}
		case "e":
			if m.done && m.exported == "" {
				filename, err := exportResults(m)
				if err != nil {
					m.exported = fmt.Sprintf("Export failed: %v", err)
				} else {
					m.exported = fmt.Sprintf("Exported to %s", filename)
				}
			}
		}

	case checkStartMsg:
		for i, c := range m.checks {
			if c.id == msg.id {
				m.checks[i].status = statusRunning
				break
			}
		}

	case checkDoneMsg:
		allDone := true
		for i, c := range m.checks {
			if c.id == msg.id {
				if msg.passed {
					m.checks[i].status = statusPass
				} else {
					m.checks[i].status = statusFail
				}
				m.checks[i].output = msg.output
				m.checks[i].errOut = msg.errOut
				m.checks[i].duration = msg.duration
			}
			if m.checks[i].status == statusPending || m.checks[i].status == statusRunning {
				allDone = false
			}
		}
		if allDone {
			m.running = false
			m.done = true
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m tuiModel) View() string {
	if m.quitting {
		return ""
	}

	var sb strings.Builder

	border := strings.Repeat("═", 42)
	sb.WriteString(styleTitle.Render("╔" + border + "╗") + "\n")
	sb.WriteString(styleTitle.Render("║  doxctl TUI Dashboard") + styleDim.Render(strings.Repeat(" ", 21)) + styleTitle.Render("║") + "\n")
	sb.WriteString(styleTitle.Render("╚" + border + "╝") + "\n\n")

	hints := styleDim.Render("[space]") + " toggle  " +
		styleDim.Render("[enter]") + " run  " +
		styleDim.Render("[e]") + " export  " +
		styleDim.Render("[q]") + " quit"
	sb.WriteString("  " + hints + "\n\n")

	passing, failing := 0, 0
	for _, c := range m.checks {
		if c.status == statusPass {
			passing++
		} else if c.status == statusFail {
			failing++
		}
	}

	for i, c := range m.checks {
		checkbox := "[ ]"
		if c.selected {
			checkbox = styleChecked.Render("[✓]")
		}

		var icon string
		switch c.status {
		case statusPending:
			icon = styleDim.Render("○")
		case statusRunning:
			icon = styleYellow.Render("◌") + " " + m.spinner.View()
		case statusPass:
			icon = styleGreen.Render("✓")
		case statusFail:
			icon = styleRed.Render("✗")
		case statusSkipped:
			icon = styleDim.Render("-")
		}

		var statusStr string
		switch c.status {
		case statusPending:
			statusStr = styleDim.Render("[pending]")
		case statusRunning:
			statusStr = styleYellow.Render("[running]")
		case statusPass:
			statusStr = styleGreen.Render(fmt.Sprintf("[%.1fs]", c.duration.Seconds()))
		case statusFail:
			statusStr = styleRed.Render(fmt.Sprintf("[%.1fs]", c.duration.Seconds()))
		case statusSkipped:
			statusStr = styleDim.Render("[skipped]")
		}

		name := c.name
		padding := 28 - len(name)
		if padding < 1 {
			padding = 1
		}
		line := fmt.Sprintf("  %s %s %s%s%s", checkbox, icon, name, strings.Repeat(" ", padding), statusStr)

		if i == m.cursor {
			sb.WriteString(styleCursor.Render(line) + "\n")
		} else {
			sb.WriteString(line + "\n")
		}
	}

	sb.WriteString("\n")

	completed := passing + failing
	total := 0
	for _, c := range m.checks {
		if c.status != statusSkipped {
			total++
		}
	}
	progress := fmt.Sprintf("  Progress: %d/%d complete", completed, total)
	if m.done {
		summary := fmt.Sprintf(" | %s %d passed", styleGreen.Render("✓"), passing)
		if failing > 0 {
			summary += fmt.Sprintf(" %s %d failed", styleRed.Render("✗"), failing)
		}
		sb.WriteString(styleBold.Render(progress) + summary + "\n")
		if m.exported != "" {
			sb.WriteString("\n  " + styleGreen.Render(m.exported) + "\n")
		} else {
			sb.WriteString("\n  " + styleDim.Render("Press [e] to export results as JSON") + "\n")
		}
	} else {
		sb.WriteString(styleDim.Render(progress) + "\n")
	}

	return sb.String()
}

func runCheck(m tuiModel, item checkItem) tea.Cmd {
	return func() tea.Msg {
		bin, err := os.Executable()
		if err != nil {
			return checkDoneMsg{id: item.id, passed: false, errOut: err.Error()}
		}

		args := append(item.args, "--output", "json")
		if m.cfgFile != "" {
			args = append(args, "--config", m.cfgFile)
		}

		start := time.Now()
		var stdout, stderr bytes.Buffer
		cmd := exec.Command(bin, args...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		runErr := cmd.Run()
		elapsed := time.Since(start)

		passed := runErr == nil
		return checkDoneMsg{
			id:       item.id,
			passed:   passed,
			output:   stdout.String(),
			errOut:   stderr.String(),
			duration: elapsed,
		}
	}
}

func startAllSelected(m tuiModel) tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range m.checks {
		if c.selected {
			c := c
			cmds = append(cmds, tea.Sequence(
				func() tea.Msg { return checkStartMsg{id: c.id} },
				runCheck(m, c),
			))
		}
	}
	return tea.Batch(cmds...)
}

type tuiExportResult struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Passed   bool   `json:"passed"`
	Duration string `json:"duration"`
	Output   string `json:"output,omitempty"`
	ErrOut   string `json:"errOut,omitempty"`
}

type tuiExport struct {
	Timestamp string            `json:"timestamp"`
	Results   []tuiExportResult `json:"results"`
}

func exportResults(m tuiModel) (string, error) {
	exp := tuiExport{
		Timestamp: time.Now().Format(time.RFC3339),
	}
	for _, c := range m.checks {
		if c.status == statusSkipped {
			continue
		}
		exp.Results = append(exp.Results, tuiExportResult{
			ID:       c.id,
			Name:     c.name,
			Passed:   c.status == statusPass,
			Duration: fmt.Sprintf("%.2fs", c.duration.Seconds()),
			Output:   c.output,
			ErrOut:   c.errOut,
		})
	}

	filename := fmt.Sprintf("doxctl-tui-results-%s.json", time.Now().Format("20060102-150405"))
	data, err := json.MarshalIndent(exp, "", "  ")
	if err != nil {
		return "", err
	}
	return filename, os.WriteFile(filename, data, 0600)
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Interactive TUI dashboard for running diagnostics",
	Long: `Launch an interactive terminal dashboard to run and monitor doxctl diagnostics.

Use arrow keys to navigate, space to toggle checks, enter to run, e to export results.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialTUIModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
