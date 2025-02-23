package managers

import "github.com/google/uuid"

type Game struct {
	ID            string
	Players       []string
	Board         [][]string
	CurrentPlayer string
}

type GameManager struct {
	games []Game
}

func (gm *GameManager) CreateGame(players []string) *Game {
	game := Game{
		ID:            uuid.New().String(),
		Players:       players,
		Board:         make([][]string, 3),
		CurrentPlayer: players[0],
	}
	gm.games = append(gm.games, game)
	return &game
}
