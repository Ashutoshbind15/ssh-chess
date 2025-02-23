package managers

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SessionManager struct {
	UserProgram map[string]*tea.Program
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		UserProgram: make(map[string]*tea.Program),
	}
}
