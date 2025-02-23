package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Ashutoshbind15/ssh-chess/managers"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
)

const (
	host = "localhost"
	port = "23234"
)

var sessionManager = managers.NewSessionManager()

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			hash := md5.Sum(key.Marshal())
			fingerprint := hex.EncodeToString(hash[:])
			ctx.SetValue("fingerprint", fingerprint)
			return true
		}),
		wish.WithMiddleware(
			myCustomBubbleteaMiddleware(),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func myCustomBubbleteaMiddleware() wish.Middleware {
	newProg := func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		p := tea.NewProgram(m, opts...)
		return p
	}
	teaHandler := func(s ssh.Session) *tea.Program {
		pty, _, active := s.Pty()

		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil
		}

		renderer := bubbletea.MakeRenderer(s)

		bg := "light"
		if renderer.HasDarkBackground() {
			bg = "dark"
		}

		txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
		quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
		msgStyle := renderer.NewStyle().Foreground(lipgloss.Color("9"))
		statusStyle := renderer.NewStyle().Foreground(lipgloss.Color("9"))

		ctx := context.Background()
		fingerprint := s.Context().Value("fingerprint").(string)

		m := model{
			term:      pty.Term,
			profile:   renderer.ColorProfile().Name(),
			bg:        bg,
			txtStyle:  txtStyle,
			quitStyle: quitStyle,
			serverResp: serverRespUIModel{
				status: statusStyle,
				msg:    msgStyle,
			},
			ctx:         ctx,
			fingerprint: fingerprint,
		}
		pg := newProg(m, append(bubbletea.MakeOptions(s), tea.WithAltScreen())...)
		sessionManager.UserProgram[fingerprint] = pg
		return pg
	}
	return bubbletea.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
}

type model struct {
	term        string
	profile     string
	bg          string
	txtStyle    lipgloss.Style
	quitStyle   lipgloss.Style
	status      int
	msg         string
	serverResp  serverRespUIModel
	fingerprint string
	ctx         context.Context
	currentPage string
}

type serverRespUIModel struct {
	status lipgloss.Style
	msg    lipgloss.Style
}

func PingServer() tea.Msg {
	fmt.Println("ping from the client")
	time.Sleep(time.Second)
	// simulates a sort of async behavior (todo: call a js chess engine to validate the move)
	return pingpong{msg: "pong"}
}

type pingpong struct {
	msg string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case pingpong:
		fmt.Println("pong from the server")
		m.status = 1
		m.msg = msg.msg
		fmt.Println("m.fingerprint", m.fingerprint)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			fmt.Println("sessionManager.UserProgram", sessionManager.UserProgram)
			sessionManager.UserProgram[m.fingerprint] = nil
			fmt.Println("sessionManager.UserProgram", sessionManager.UserProgram)
			return m, tea.Quit
		case "p":
			m.status = -1
			m.msg = "loading"
			return m, PingServer
		}
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Your term is %s\nBackground: %s\nColor Profile: %s", m.term, m.bg, m.profile)
	return m.txtStyle.Render(s) + "\n\n" + m.quitStyle.Render("Press 'q' to quit") + "\n\n" +
		"yo" + "\n\n" + m.serverResp.msg.Render(m.msg) + "\n\n" + m.serverResp.status.Render(strconv.Itoa(m.status))
}
