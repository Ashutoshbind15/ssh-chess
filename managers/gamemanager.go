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
	games          []Game
	sessionManager *SessionManager
}

func (gm *GameManager) SetSessionManager(sessionManager *SessionManager) {
	gm.sessionManager = sessionManager
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
		ChessGame:     chess.NewGame(chess.UseNotation(chess.UCINotation{})),
	}
	gm.games = append(gm.games, game)
	return &game
}

func (gm *GameManager) GetGame(id string) *Game {
	for i := range gm.games {
		if gm.games[i].ID == id {
			return &gm.games[i]
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

	var opponent string

	if game.CurrentPlayer == game.Players[0] {
		opponent = game.Players[1]
	} else {
		opponent = game.Players[0]
	}

	game.CurrentPlayer = opponent

	oppProgram := gm.sessionManager.UserProgram[opponent]
	if oppProgram == nil {
		return types.ChessBoardError{Error: "Opponent not found"}
	}

	oppProgram.Send(types.BoardUpdateResponse{Fen: chessGame.FEN()})
	return types.BoardUpdateResponse{Fen: chessGame.FEN()}
}

func (gm *GameManager) MakeMove(id, from, to string) tea.Cmd {
	return func() tea.Msg {
		return gm.move(id, from, to)
	}
}
