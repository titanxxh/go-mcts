package mcts

// Move represents a move in the game.
type Move interface {
}

// GameState is the interface a game supports to satisfy the MCTS.
type GameState interface {
	Clone() GameState       // Clone the game state, a deep copy.
	AvailableMoves() []Move // Return all the viable moves given the current game state. For a finished game, nil.
	MakeMove(move Move)     // Take an action, changing the game state.
	Score(player int) float64
	PlayerJustMoved() int
}
