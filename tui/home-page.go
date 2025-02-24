package tui

import (
	"fmt"
	"strconv"
)

func (m model) IntroPageRenderer() string {
	s := fmt.Sprintf("Your term is %s\nBackground: %s\nColor Profile: %s", m.term, m.bg, m.profile)
	return m.txtStyle.Render(s) + "\n\n" + m.quitStyle.Render("Press 'q' to quit") + "\n\n" +
		"yo" + "\n\n" + m.serverResp.msg.Render(m.msg) + "\n\n" + m.serverResp.status.Render(strconv.Itoa(m.status)) + "\n\n" +
		m.txtStyle.Render(m.statusText)
}
