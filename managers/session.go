package managers

import (
	tea "github.com/charmbracelet/bubbletea"

	"math/rand"

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
	gameManager *GameManager
}

func (s *SessionManager) AddUser(fingerprint string, program *tea.Program) {
	s.Users = append(s.Users, User{fingerprint: fingerprint, currentStatus: "unpaired", program: program})
}

func NewSessionManager(gameManager *GameManager) *SessionManager {
	return &SessionManager{
		UserProgram: make(map[string]*tea.Program),
		gameManager: gameManager,
	}
}

func (s SessionManager) StartPairing(fingerprint string) tea.Cmd {

	lof := func() tea.Msg {

		currentUserIndex := -1
		isPaired := false

		for i, user := range s.Users {
			if user.currentStatus == "unpaired" && user.fingerprint != fingerprint {

				randomColor := rand.Intn(2) == 0

				whitePlayer := ""
				blackPlayer := ""

				if randomColor {
					whitePlayer = user.fingerprint
					blackPlayer = fingerprint
				} else {
					whitePlayer = fingerprint
					blackPlayer = user.fingerprint
				}

				game := s.gameManager.CreateGame([]string{whitePlayer, blackPlayer}, whitePlayer)

				user.program.Send(common.PairedResponse{Opponent: blackPlayer, Color: randomColor, GameID: game.ID})
				cuserProgram := s.UserProgram[fingerprint]
				cuserProgram.Send(common.PairedResponse{Opponent: whitePlayer, Color: !randomColor, GameID: game.ID})

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
