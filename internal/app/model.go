package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oadultradeepfield/gomodoro/internal/ui"
)

func New() Model {
	inputs := make([]textinput.Model, 3)

	for i := range inputs {
		t := textinput.New()
		t.CharLimit = 2
		t.Width = 4
		inputs[i] = t
	}

	inputs[0].Placeholder = "25"
	inputs[0].SetValue("25")
	inputs[0].Focus()

	inputs[1].Placeholder = "5"
	inputs[1].SetValue("5")

	inputs[2].Placeholder = "15"
	inputs[2].SetValue("15")

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = ui.SpinnerStyle

	return Model{
		phase:        PhaseSetup,
		workDuration: 25 * time.Minute,
		shortBreak:   5 * time.Minute,
		longBreak:    15 * time.Minute,
		sessionCount: 1,
		inputs:       inputs,
		spinner:      s,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case TickMsg:
		return m.handleTick()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case NotifyMsg:
		return m, nil
	}

	if m.phase == PhaseSetup {
		return m.updateInputs(msg)
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "enter":
		if m.phase == PhaseSetup {
			return m.startTimer()
		}
	case " ":
		if m.phase != PhaseSetup {
			m.paused = !m.paused
			if !m.paused {
				return m, tea.Batch(tickCmd(), m.spinner.Tick)
			}
		}
	case "tab", "shift+tab":
		if m.phase == PhaseSetup {
			return m.cycleFocus(msg.String() == "shift+tab")
		}
	}
	return m, nil
}

func (m Model) cycleFocus(reverse bool) (tea.Model, tea.Cmd) {
	if reverse {
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
	} else {
		m.focusIndex++
		if m.focusIndex >= len(m.inputs) {
			m.focusIndex = 0
		}
	}

	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}

	return m, nil
}

func (m Model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) startTimer() (tea.Model, tea.Cmd) {
	m.workDuration = m.parseDuration(m.inputs[0].Value(), 25)
	m.shortBreak = m.parseDuration(m.inputs[1].Value(), 5)
	m.longBreak = m.parseDuration(m.inputs[2].Value(), 15)

	m.phase = PhaseWorking
	m.remaining = m.workDuration
	m.sessionCount = 1

	return m, tea.Batch(tickCmd(), m.spinner.Tick)
}

func (m Model) parseDuration(value string, defaultVal int) time.Duration {
	var minutes int
	if _, err := fmt.Sscanf(value, "%d", &minutes); err != nil || minutes <= 0 {
		minutes = defaultVal
	}
	if minutes > 99 {
		minutes = 99
	}
	return time.Duration(minutes) * time.Minute
}

func (m Model) handleTick() (tea.Model, tea.Cmd) {
	if m.paused || m.phase == PhaseSetup {
		return m, nil
	}

	m.remaining -= time.Second

	if m.remaining <= 0 {
		return m.transitionPhase()
	}

	return m, tickCmd()
}

func (m Model) transitionPhase() (tea.Model, tea.Cmd) {
	var title, message string

	switch m.phase {
	case PhaseWorking:
		if m.sessionCount >= 4 {
			m.phase = PhaseLongBreak
			m.remaining = m.longBreak
			m.sessionCount = 0
			title = "Long Break"
			message = fmt.Sprintf("Take a %d minute break!", int(m.longBreak.Minutes()))
		} else {
			m.phase = PhaseShortBreak
			m.remaining = m.shortBreak
			title = "Short Break"
			message = fmt.Sprintf("Take a %d minute break!", int(m.shortBreak.Minutes()))
		}
	case PhaseShortBreak, PhaseLongBreak:
		m.phase = PhaseWorking
		m.remaining = m.workDuration
		m.sessionCount++
		title = "Work Time"
		message = fmt.Sprintf("Focus for %d minutes! Session %d/4", int(m.workDuration.Minutes()), m.sessionCount)
	}

	return m, tea.Batch(tickCmd(), m.spinner.Tick, notifyCmd(title, message))
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	switch m.phase {
	case PhaseSetup:
		return m.viewSetup()
	default:
		return m.viewTimer()
	}
}

func (m Model) viewSetup() string {
	var b strings.Builder

	title := ui.DimStyle.Render("Gomodoro")
	b.WriteString(title + "\n\n")

	labels := []string{"Work duration", "Short break", "Long break"}
	for i, input := range m.inputs {
		label := ui.LabelStyle.Render(fmt.Sprintf("%-15s", labels[i]))
		field := input.View()
		suffix := ui.DimStyle.Render(" min")
		b.WriteString(fmt.Sprintf("%s  %s%s\n", label, field, suffix))
	}

	b.WriteString("\n")
	b.WriteString(ui.DimStyle.Render("Press Enter to start · Tab to switch fields · q to quit"))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String())
}

func (m Model) viewTimer() string {
	var b strings.Builder

	statusBar := m.renderStatusBar()
	b.WriteString(statusBar + "\n\n")

	timer := m.formatTime(m.remaining)
	timerDisplay := ui.RenderStyledASCIITime(timer, ui.TimerStyle)
	b.WriteString(timerDisplay + "\n\n")

	var help string
	if m.paused {
		help = "Press Space to resume · q to quit"
	} else {
		help = "Press Space to pause · q to quit"
	}
	b.WriteString(ui.DimStyle.Render(help))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String())
}

func (m Model) renderStatusBar() string {
	var phaseText string
	var style lipgloss.Style

	switch m.phase {
	case PhaseWorking:
		phaseText = "Working"
		style = ui.PhaseWorkingStyle
	case PhaseShortBreak:
		phaseText = "Short Break"
		style = ui.PhaseBreakStyle
	case PhaseLongBreak:
		phaseText = "Long Break"
		style = ui.PhaseBreakStyle
	}

	var spinnerView string
	if !m.paused {
		spinnerView = m.spinner.View() + " "
	} else {
		spinnerView = ui.DimStyle.Render("⏸ ")
	}

	phase := style.Render(phaseText)
	session := ui.DimStyle.Render(fmt.Sprintf("%d/4", m.sessionCount))

	return spinnerView + phase + "  " + session
}

func (m Model) formatTime(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
