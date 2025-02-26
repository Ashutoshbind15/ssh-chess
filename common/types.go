package common

type Pingpong struct {
	Msg string
}

type PairedResponse struct {
	Opponent string
	Color    bool
	GameID   string
}

type BoardUpdateResponse struct {
	Fen string
}

type ChessBoardError struct {
	Error string
}

type ChessVictory struct {
	Message string
}
