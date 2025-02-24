package tui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

func RenderBoard(board [8][8]string) string {
	renderer := lipgloss.NewRenderer(os.Stdout)
	cellStyle := renderer.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))

	rows := make([]string, 8)
	for i := 0; i < 8; i++ {
		cells := make([]string, 8)
		for j := 0; j < 8; j++ {
			cells[j] = cellStyle.Render(board[i][j])
		}
		rows[i] = lipgloss.JoinHorizontal(lipgloss.Left, cells...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
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

func RenderChessPage(board [8][8]string, color bool) string {
	if color {
		return RenderBoard(ReverseBoardSide(board))
	}
	return RenderBoard(board)
}
