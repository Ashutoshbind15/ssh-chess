package managers

import (
	"github.com/google/uuid"
	"github.com/notnil/chess"

	types "github.com/Ashutoshbind15/ssh-chess/common"
	tea "github.com/charmbracelet/bubbletea"
)

type Game struct {
	ID            string
	Players       []string
	CurrentPlayer string
	ChessGame     *chess.Game
}

type GameManager struct {
	games []Game
}

func NewGameManager() *GameManager {
	return &GameManager{
		games: []Game{},
	}
}

func (gm *GameManager) CreateGame(players []string, currentPlayer string) *Game {
	game := Game{
		ID:            uuid.New().String(),
		Players:       players,
		CurrentPlayer: currentPlayer,
		ChessGame:     chess.NewGame(),
	}
	gm.games = append(gm.games, game)
	return &game
}

func (gm *GameManager) GetGame(id string) *Game {
	for _, game := range gm.games {
		if game.ID == id {
			return &game
		}
	}
	return nil
}

func (gm *GameManager) move(id, from, to string) tea.Msg {
	game := gm.GetGame(id)
	if game == nil {
		return types.ChessBoardError{Error: "Game not found"}
	}

	chessGame := game.ChessGame

	err := chessGame.MoveStr(from + to)
	if err != nil {
		return types.ChessBoardError{Error: err.Error()}
	}

	return types.BoardUpdateResponse{Fen: chessGame.FEN()}
}

func (gm *GameManager) MakeMove(id, from, to string) tea.Cmd {
	return func() tea.Msg {
		return gm.move(id, from, to)
	}
}
