package managers

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Ashutoshbind15/ssh-chess/common"
)

type User struct {
	fingerprint   string
	program       *tea.Program
	currentStatus string
}

type SessionManager struct {
	UserProgram map[string]*tea.Program
	Users       []User
}

func (s *SessionManager) AddUser(fingerprint string, program *tea.Program) {
	s.Users = append(s.Users, User{fingerprint: fingerprint, currentStatus: "unpaired", program: program})
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		UserProgram: make(map[string]*tea.Program),
	}
}

func (s SessionManager) StartPairing(fingerprint string) tea.Cmd {

	lof := func() tea.Msg {

		currentUserIndex := -1
		isPaired := false

		for i, user := range s.Users {
			if user.currentStatus == "unpaired" && user.fingerprint != fingerprint {
				user.program.Send(common.PairedResponse{Opponent: fingerprint})
				cuserProgram := s.UserProgram[fingerprint]
				cuserProgram.Send(common.PairedResponse{Opponent: user.fingerprint})

				// update the current status of the current user and the opponent
				s.Users[i].currentStatus = "paired"
				isPaired = true
			}

			if user.fingerprint == fingerprint {
				currentUserIndex = i
			}
		}

		if isPaired {
			s.Users[currentUserIndex].currentStatus = "paired"
		}
		return nil
	}

	return lof
}
