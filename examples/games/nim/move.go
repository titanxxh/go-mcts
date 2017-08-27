package nim

import (
	"log"
)

// NimMove is a single move in a nim game.
type NimMove struct {
	playerId int    // Which of the players is taking this move.
	chips    uint64 // The number of chips the player is taking, 1 - 3.
}

// newNimMove creates a new Nim move.
func newNimMove(playerId int, chips uint64) *NimMove {
	if chips < 1 || chips > 3 {
		log.Panicf("Nim move cannot be to remove %d chips.", chips)
	}
	return &NimMove{
		playerId: playerId,
		chips:    chips,
	}
}

// Log reports the current game state.
func (m *NimMove) Log() {
	log.Printf("player %d takes %d chips", m.playerId, m.chips)
}
