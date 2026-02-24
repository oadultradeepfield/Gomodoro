package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oadultradeepfield/gomodoro/internal/notify"
)

type TickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type NotifyMsg struct {
	Title   string
	Message string
}

func notifyCmd(title, message string) tea.Cmd {
	return func() tea.Msg {
		notify.Send(title, message)
		return NotifyMsg{Title: title, Message: message}
	}
}
