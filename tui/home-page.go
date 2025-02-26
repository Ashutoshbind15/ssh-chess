package tui

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) IntroPageRenderer() string {
	s := fmt.Sprintf("Your term is %s\nBackground: %s\nColor Profile: %s", m.term, m.bg, m.profile)
	return m.txtStyle.Render(s) + "\n\n" + m.quitStyle.Render("Press 'q' to quit") + "\n\n" +
		"yo" + "\n\n" + m.serverResp.msg.Render(m.msg) + "\n\n" + m.serverResp.status.Render(strconv.Itoa(m.status)) + "\n\n" +
		m.txtStyle.Render(m.statusText)
}

func (m model) IntroPageUpdate(msg tea.Msg) (model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			m.status = -1
			m.msg = "loading"
			cmd = PingServer
		case "s":
			m.statusText = "waiting for pairing up"
			cmd = m.sessionManager.StartPairing(m.fingerprint)
		}
	}

	return m, cmd
}
