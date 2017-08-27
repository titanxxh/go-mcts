/*
	Package nim implements a game state usable by go-mcts

	From: http://mcts.ai/
	A state of the game Nim. In Nim, players alternately take 1,2 or 3 chips with the winner being the player
	to take the last chip. In Nim any initial state of the form 4n+k for k = 1,2,3 is a win for player 1
	(by choosing k) chips. Any initial state of the form 4n is a win for player 2.
*/
package nim

import (
	"github.com/titanxxh/go-mcts"
	"log"
)

const (
	_MAX_PICKABLE_CHIPS = 3 // The max chips a player can pick on their move.
)

// NimSate is the state of a Nim game.
type NimState struct {
	// Public members are used by the scoring function and running the game.
	JustMovedPlayerId int // Who just moved, if they took the last chip they win.
	ActivePlayerId    int // This is the player who's current turn it is.

	// Private members are just used for running the game.
	chips uint64 // How many chips are left?
}

// NewNimState creates a new nim game.
func NewNimState(chips uint64) *NimState {
	return &NimState{
		JustMovedPlayerId: 2,
		ActivePlayerId:    1, // The first player starts.
		chips:             chips,
	}
}

// Log reports the current game state.
func (g *NimState) Log() {
	log.Printf("CHIPS: %d", g.chips)
}

// LogWinner reports the winner.
func (g *NimState) LogWinner() {
	log.Printf("PLAYER %d WINS!", g.JustMovedPlayerId)
}

// Clone makes a deep copy of the game state.
func (g *NimState) Clone() mcts.GameState {
	// Return the new state.
	return &NimState{
		JustMovedPlayerId: g.JustMovedPlayerId,
		ActivePlayerId:    g.ActivePlayerId,
		chips:             g.chips,
	}
}

// AvailableMoves returns all the available moves.
func (g *NimState) AvailableMoves() []mcts.Move {
	var maxChipsPickable uint64 = g.chips
	if maxChipsPickable > _MAX_PICKABLE_CHIPS {
		maxChipsPickable = _MAX_PICKABLE_CHIPS
	}
	var moves []mcts.Move
	var pickedChips uint64
	for pickedChips = 1; pickedChips <= maxChipsPickable; pickedChips++ {
		moves = append(moves, newNimMove(g.ActivePlayerId, pickedChips))
	}
	return moves // Will be nil if the game is over (no pickable chips).
}

// MakeMove makes a move in the game state, changing it.
func (g *NimState) MakeMove(move mcts.Move) {
	// Convert the move to a form we can use.
	var nimMove *NimMove = move.(*NimMove)
	g.chips -= nimMove.chips
	// It is now the next player's turn.
	g.JustMovedPlayerId, g.ActivePlayerId = g.ActivePlayerId, g.JustMovedPlayerId
}

// scoreNim scores the game state from a player's perspective, returning 0.0 (lost), 0.5 (in progress), 1.0 (won)
func (g *NimState) Score(player int) float64 {
	// Is the game over or still in progress?
	var moves []mcts.Move = g.AvailableMoves()
	if len(moves) > 0 {
		// The game is still in progress.
		return 0.5 // Consider it a neutral state (0.0-1.0)
	}

	// The game is over.
	if player == g.JustMovedPlayerId {
		// The game is over and we were the last player to move. We win!
		return 1.0
	}
	// We didn't win.
	return 0.0
}

func (g *NimState) PlayerJustMoved() int {
	return g.JustMovedPlayerId
}
