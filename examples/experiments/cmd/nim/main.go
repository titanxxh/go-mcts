/*
	nim the mcts algorithm for the game of nim.
*/
package main

import (
	"flag"
	"github.com/titanxxh/go-mcts"
	"github.com/titanxxh/go-mcts/examples/games/nim"
	"log"
	"os"
)

func main() {
	// Pass the configuration files as parameters:
	//
	//   bin/nim -ucbc=1.0 -chips=10
	//   bin/nim -h
	//
	var ucbC *float64 = flag.Float64("ucbc", 1.0, "the constant biasing exploitation vs exploration")
	var chips *uint64 = flag.Uint64("chips", 100, "the number of chips in the starting state")
	flag.Parse()

	var experimentName string = os.Args[0]
	log.Printf("Experiment game: '%s'\n", experimentName)

	opt := mcts.UctOption{
		Iterations:  1000,
		Simulations: 100,
		UcbConstant: *ucbC,
	}

	// Create the initial game state.
	var state *nim.NimState = nim.NewNimState(*chips)

	// Play until the game is over (no more available moves).
	for len(state.AvailableMoves()) > 0 {

		// Log the current game state.
		state.Log()

		// What is the next active player's move?
		var move mcts.Move = mcts.Uct(state, opt)
		state.MakeMove(move)

		// Report the action taken.
		var nimMove *nim.NimMove = move.(*nim.NimMove)
		nimMove.Log()
	}

	// Report winner.
	state.LogWinner()

	log.Println("Experiment Complete.")
	os.Exit(0)
}
