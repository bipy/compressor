package tui

import (
	"fmt"
	"strings"
	"time"

	"compressor/internal/reporter"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

type Model struct {
	reporter   *reporter.Reporter
	viewport   viewport.Model
	quitting   bool
	paused     bool
	showHelp   bool
	width      int
	height     int
	done       chan struct{}
}

func NewModel(r *reporter.Reporter, done chan struct{}) Model {
	return Model{
		reporter: r,
		viewport: viewport.New(80, 20),
		done:     done,
		width:    80,
		height:   24,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		tea.EnterAltScreen,
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 8
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "p":
			m.paused = !m.paused
		case "v":
			m.showHelp = !m.showHelp
		case "?":
			m.showHelp = !m.showHelp
		}

	case tickMsg:
		if m.done != nil {
			select {
			case <-m.done:
				m.quitting = true
				return m, tea.Quit
			default:
			}
		}
		return m, tickCmd()
	}

	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	snap := m.reporter.Snapshot()

	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	b.WriteString(headerStyle.Render("Compressor") + "\n")

	// Stats line
	statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	stats := fmt.Sprintf(
		"Total: %d | %s | %s | Elapsed: %s",
		snap.Processed,
		successStyle.Render(fmt.Sprintf("✓ %d", snap.Successes)),
		failStyle.Render(fmt.Sprintf("✗ %d", snap.Failures)),
		snap.Elapsed.Round(time.Second),
	)
	b.WriteString(statsStyle.Render(stats) + "\n")

	// Throughput line
	throughputStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	throughput := fmt.Sprintf(
		"Input: %.1fMB | Output: %.1fMB | Ratio: %.2f | Speed: %.1fMB/s",
		float64(snap.InputBytes)/(1024*1024),
		float64(snap.OutputBytes)/(1024*1024),
		snap.CompressionRatio,
		snap.BytesPerSecond/(1024*1024),
	)
	b.WriteString(throughputStyle.Render(throughput) + "\n")

	// Progress bar
	if snap.Discovered > 0 {
		pct := float64(snap.Processed) / float64(snap.Discovered) * 100
		barWidth := m.width - 20
		if barWidth < 10 {
			barWidth = 10
		}
		filled := int(pct / 100 * float64(barWidth))
		if filled > barWidth {
			filled = barWidth
		}

		barStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		b.WriteString(barStyle.Render(fmt.Sprintf("[%s] %.1f%%", bar, pct)) + "\n")
	}

	b.WriteString("\n")

	// Recent jobs
	recentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	b.WriteString(recentStyle.Render("Recent:") + "\n")
	for i := len(snap.RecentJobs) - 1; i >= 0 && i >= len(snap.RecentJobs)-5; i-- {
		r := snap.RecentJobs[i]
		icon := "✓"
		style := successStyle
		if r.Status != "done" {
			icon = "✗"
			style = failStyle
		}
		line := fmt.Sprintf("  %s %s (%.1fKB)",
			icon,
			truncatePath(r.Job.SourcePath, 40),
			float64(r.OutputBytes)/1024,
		)
		b.WriteString(style.Render(line) + "\n")
	}

	// Failures detail
	if m.showHelp && snap.Failures > 0 {
		b.WriteString("\n")
		b.WriteString(failStyle.Render(fmt.Sprintf("Failures (%d):", snap.Failures)) + "\n")
		for _, f := range m.reporter.Failures() {
			b.WriteString(failStyle.Render(fmt.Sprintf("  %s: %v", f.Job.SourcePath, f.Err)) + "\n")
		}
	}

	// Help
	if m.showHelp {
		b.WriteString("\n")
		helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		b.WriteString(helpStyle.Render("q: quit | p: pause | v/f: toggle failures | ?: help") + "\n")
	}

	return b.String()
}

func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
}

func IsTTY() bool {
	return isTTY()
}
