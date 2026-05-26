package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00d4ff"))
	green   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff88")).Bold(true)
	gray    = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	white   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Bold(true)
	red     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4444"))
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00d4ff")).
			Padding(0, 1)
)

type Model struct {
	messages    []string
	cursor      int
	selected    int
	editing     bool
	editBuffer  string
	committed   bool
	commitMsg   string
	quitting    bool
	files       []string
	insertions  int
	deletions   int
}

type CommitResult struct {
	Message string
	Quit    bool
}

func NewModel(messages []string, files []string, insertions int, deletions int) Model {
	return Model{
		messages:   messages,
		cursor:     0,
		selected:   -1,
		files:      files,
		insertions: insertions,
		deletions:  deletions,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "enter":
				m.selected = m.cursor
				m.commitMsg = m.editBuffer
				m.committed = true
				return m, tea.Quit
			case "esc":
				m.editing = false
				m.editBuffer = ""
			case "backspace":
				if len(m.editBuffer) > 0 {
					m.editBuffer = m.editBuffer[:len(m.editBuffer)-1]
				}
			default:
				if len(msg.String()) == 1 {
					m.editBuffer += msg.String()
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.messages)-1 {
				m.cursor++
			}
		case "1":
			if len(m.messages) > 0 {
				m.selected = 0
				m.commitMsg = m.messages[0]
				m.committed = true
				return m, tea.Quit
			}
		case "2":
			if len(m.messages) > 1 {
				m.selected = 1
				m.commitMsg = m.messages[1]
				m.committed = true
				return m, tea.Quit
			}
		case "3":
			if len(m.messages) > 2 {
				m.selected = 2
				m.commitMsg = m.messages[2]
				m.committed = true
				return m, tea.Quit
			}
		case "e":
			m.editing = true
			if m.cursor < len(m.messages) {
				m.editBuffer = m.messages[m.cursor]
			}
		case "enter":
			if m.cursor < len(m.messages) {
				m.selected = m.cursor
				m.commitMsg = m.messages[m.cursor]
				m.committed = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return gray.Render("  cancelled\n")
	}

	if m.committed {
		return green.Render("  ✓ committed: ") + white.Render(m.commitMsg) + "\n"
	}

	var sb strings.Builder


	sb.WriteString("\n")
	sb.WriteString(cyan.Render("  ◆ KOMMIT") + gray.Render(" — ai commit messages\n\n"))

	sb.WriteString(fmt.Sprintf("  %s %d files  %s +%d  %s -%d\n\n",
		gray.Render("staged:"),
		len(m.files),
		green.Render(""),
		m.insertions,
		red.Render(""),
		m.deletions,
	))

	for i, msg := range m.messages {
		prefix := fmt.Sprintf("  %d. ", i+1)
		if i == m.cursor {
			sb.WriteString(cyan.Render(prefix) + white.Render(msg) + "\n")
		} else {
			sb.WriteString(gray.Render(prefix) + msg + "\n")
		}
	}

	sb.WriteString("\n")

	if m.editing {
		sb.WriteString(cyan.Render("  edit: ") + white.Render(m.editBuffer) + "█\n\n")
		sb.WriteString(gray.Render("  enter: confirm  esc: cancel\n"))
	} else {
		sb.WriteString(gray.Render("  [1/2/3] pick  [↑↓] navigate  [enter] select  [e] edit  [q] quit\n"))
	}

	return sb.String()
}

func (m Model) Result() CommitResult {
	return CommitResult{
		Message: m.commitMsg,
		Quit:    m.quitting,
	}
}

func RunUI(messages []string, files []string, insertions int, deletions int) (string, bool, error) {
	m := NewModel(messages, files, insertions, deletions)
	p := tea.NewProgram(m)

	result, err := p.Run()
	if err != nil {
		return "", false, err
	}

	final := result.(Model)
	if final.quitting {
		return "", true, nil
	}

	return final.commitMsg, false, nil
}
