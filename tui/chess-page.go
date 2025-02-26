package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func convertToChessboardPosition(x, y int, color bool) string {
	// convert 0,0 to a8 and 7,7 to h1
	if color {
		return fmt.Sprintf("%c%d", 'a'+x, 8-y)
	}

	// h to a and 1 to 8
	return fmt.Sprintf("%c%d", 'h'-x, y+1)
}

func (m model) RenderBoard(board [8][8]string) string {

	cellStyle := m.renderer.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))
	var cell string

	rows := make([]string, 8)
	for i := 0; i < 8; i++ {
		cells := make([]string, 8)
		for j := 0; j < 8; j++ {
			if m.selected == convertToChessboardPosition(j, i, m.color) {
				cell = cellStyle.UnsetBorderForeground().BorderForeground(lipgloss.Color("190")).Render(board[i][j])
			} else {
				cell = cellStyle.Render(board[i][j])
			}
			cells[j] = m.zone.Mark(convertToChessboardPosition(j, i, m.color), cell)
		}
		rows[i] = lipgloss.JoinHorizontal(lipgloss.Left, cells...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m model) RenderBoardFromFen(fen string) string {
	board := ParseFEN(fen)
	return m.RenderBoard(board)
}

func ParseFEN(fen string) [8][8]string {
	var board [8][8]string
	ranks := strings.Split(fen, " ")[0] // Get only the piece placement part of FEN
	rows := strings.Split(ranks, "/")

	for i, row := range rows {
		fileIndex := 0
		for _, ch := range row {
			if ch >= '1' && ch <= '8' { // Empty squares
				emptyCount := int(ch - '0')
				for j := 0; j < emptyCount; j++ {
					board[i][fileIndex] = " "
					fileIndex++
				}
			} else { // Piece character
				board[i][fileIndex] = string(ch)
				fileIndex++
			}
		}
	}
	return board
}

func InitRepresentation() [8][8]string {
	board := [8][8]string{}

	board[0] = [8]string{"r", "n", "b", "q", "k", "b", "n", "r"}
	board[1] = [8]string{"p", "p", "p", "p", "p", "p", "p", "p"}
	board[2] = [8]string{" ", " ", " ", " ", " ", " ", " ", " "}
	board[3] = [8]string{" ", " ", " ", " ", " ", " ", " ", " "}
	board[4] = [8]string{" ", " ", " ", " ", " ", " ", " ", " "}
	board[5] = [8]string{" ", " ", " ", " ", " ", " ", " ", " "}
	board[6] = [8]string{"P", "P", "P", "P", "P", "P", "P", "P"}
	board[7] = [8]string{"R", "N", "B", "Q", "K", "B", "N", "R"}

	return board
}

func ReverseBoardSide(board [8][8]string) [8][8]string {
	reversed := [8][8]string{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			reversed[7-i][7-j] = board[i][j]
		}
	}
	return reversed
}

func (m model) RenderChessPage() string {

	res := ""

	if m.color {
		res += m.RenderBoardFromFen(m.chessBoard)
	} else {
		res += m.RenderBoard(ReverseBoardSide(ParseFEN(m.chessBoard)))
	}

	res += "\n\n"
	res += m.txtStyle.Render(m.fingerprint)

	return res
}

func (m model) chessUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:

		if msg.Action != tea.MouseActionRelease || msg.Button != tea.MouseButtonLeft {
			return m, nil
		}

		// iterate over all the zones and check if the mouse is over any of them, each zone being a string from a8 to h1

		doesClick := false

		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if m.zone.Get(convertToChessboardPosition(j, i, m.color)).InBounds(msg) {
					doesClick = true
					if m.selected == "" {
						m.selected = convertToChessboardPosition(j, i, m.color)
					} else {
						tcmd := m.gameManager.MakeMove(m.gameId, m.selected, convertToChessboardPosition(j, i, m.color))
						m.selected = ""
						return m, tcmd
					}
				}
			}
		}

		if !doesClick {
			m.selected = ""
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.currentPage = "home"
			return m, nil
		}
	}

	return m, nil
}
