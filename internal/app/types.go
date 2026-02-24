package app

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

type Phase int

const (
	PhaseSetup Phase = iota
	PhaseWorking
	PhaseShortBreak
	PhaseLongBreak
)

type Model struct {
	phase        Phase
	workDuration time.Duration
	shortBreak   time.Duration
	longBreak    time.Duration
	remaining    time.Duration
	sessionCount int
	paused       bool
	width        int
	height       int
	inputs       []textinput.Model
	focusIndex   int
	spinner      spinner.Model
}
