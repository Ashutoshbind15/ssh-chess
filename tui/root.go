package tui

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Ashutoshbind15/ssh-chess/common"
	"github.com/Ashutoshbind15/ssh-chess/managers"
)

type serverRespUIModel struct {
	status lipgloss.Style
	msg    lipgloss.Style
}

type model struct {
	term           string
	profile        string
	bg             string
	txtStyle       lipgloss.Style
	quitStyle      lipgloss.Style
	status         int
	msg            string
	serverResp     serverRespUIModel
	fingerprint    string
	ctx            context.Context
	currentPage    string
	sessionManager *managers.SessionManager
	statusText     string
	chessBoard     [8][8]string
	color          bool
}

func NewModel(renderer *lipgloss.Renderer, fingerprint string, sessionManager *managers.SessionManager) model {

	bg := "light"
	if renderer.HasDarkBackground() {
		bg = "dark"
	}

	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	msgStyle := renderer.NewStyle().Foreground(lipgloss.Color("9"))
	statusStyle := renderer.NewStyle().Foreground(lipgloss.Color("9"))

	ctx := context.Background()

	m := model{
		profile:   renderer.ColorProfile().Name(),
		bg:        bg,
		txtStyle:  txtStyle,
		quitStyle: quitStyle,
		serverResp: serverRespUIModel{
			status: statusStyle,
			msg:    msgStyle,
		},
		ctx:            ctx,
		fingerprint:    fingerprint,
		sessionManager: sessionManager,
		currentPage:    "home",
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func PingServer() tea.Msg {
	fmt.Println("ping from the client")
	time.Sleep(time.Second)
	// simulates a sort of async behavior (todo: call a js chess engine to validate the move)
	return common.Pingpong{Msg: "pong"}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case common.Pingpong:
		fmt.Println("pong from the server")
		m.status = 1
		m.msg = msg.Msg
		fmt.Println("m.fingerprint", m.fingerprint)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// Todo: cleanup the user from the session manager
			m.sessionManager.UserProgram[m.fingerprint] = nil
			return m, tea.Quit
		case "p":
			m.status = -1
			m.msg = "loading"
			return m, PingServer
		case "s":
			m.statusText = "waiting for pairing up"
			return m, m.sessionManager.StartPairing(m.fingerprint)

		}

	case common.PairedResponse:
		m.statusText = "paired with " + msg.Opponent
		m.currentPage = "chess"
		m.color = msg.Color
		return m, nil

	}

	return m, nil
}

func (m model) View() string {
	switch m.currentPage {
	case "home":
		return m.IntroPageRenderer()
	case "chess":
		m.chessBoard = InitRepresentation()
		return RenderChessPage(m.chessBoard, m.color)
	}
	return ""
}
