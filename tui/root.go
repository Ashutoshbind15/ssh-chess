package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/Ashutoshbind15/ssh-chess/common"
	"github.com/Ashutoshbind15/ssh-chess/managers"
	"github.com/Ashutoshbind15/ssh-chess/theme"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
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

	gameManager *managers.GameManager
	statusText  string
	chessBoard  string
	gameId      string
	color       bool
	selected    string

	renderer        *lipgloss.Renderer
	theme           theme.Theme
	viewport        viewport.Model
	isViewportReady bool
	zone            *zone.Manager
}

func NewModel(renderer *lipgloss.Renderer, fingerprint string, sessionManager *managers.SessionManager, gameManager *managers.GameManager) model {

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
		gameManager:    gameManager,
		currentPage:    "home",
		renderer:       renderer,
		theme:          theme.BasicTheme(renderer),
		chessBoard:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		gameId:         "",
		zone:           zone.New(),
		selected:       "",
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
		m.status = 1
		m.msg = msg.Msg

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// Todo: cleanup the user from the session manager
			m.sessionManager.UserProgram[m.fingerprint] = nil
			m.zone.Close()
			return m, tea.Quit
		}

	case common.PairedResponse:
		m.statusText = "paired with " + msg.Opponent
		m.currentPage = "chess"
		m.color = msg.Color
		m.gameId = msg.GameID

	case common.BoardUpdateResponse:
		m.chessBoard = msg.Fen

	case tea.WindowSizeMsg:

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.isViewportReady {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.getContent())
			m.isViewportReady = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch m.currentPage {
	case "home":
		m, cmd = m.IntroPageUpdate(msg)
	case "chess":
		m, cmd = m.chessUpdate(msg)
	}

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.viewport.SetContent(m.getContent())
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.isViewportReady || m.zone == nil {
		return "initializing..."
	}
	return m.zone.Scan(lipgloss.JoinVertical(lipgloss.Center, m.headerView(), m.viewport.View(), m.footerView()))
}

func (m model) getContent() string {
	switch m.currentPage {
	case "home":
		return m.IntroPageRenderer()
	case "chess":
		return m.RenderChessPage()
	}
	return ""
}

func (m model) headerView() string {
	return ""
}

func (m model) footerView() string {
	return ""
}
