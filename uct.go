/*
	package mcts is an implementation of a Monte Carlo Tree Search.

	More information and example code can be found here: http://mcts.ai/
*/
package mcts

import (
	"math"
	"math/rand"
	"sort"
)

type UctOption struct {
	Iterations  uint
	Simulations uint
	UcbConstant float64 // large->bfs, small->dfs
}

// Uct is an Upper Confidence Bound Tree search through game stats for an optimal move, given a starting game state.
func Uct(state GameState, opt UctOption) Move {
	// Find the best move given a fixed number of state explorations.
	root := newTreeNode(nil, nil, state, opt.UcbConstant)
	for i := 0; i < int(opt.Iterations); i++ {
		node := root
		simulatedState := root.state.Clone()

		// 1. Select.
		// Find the node we wish to explore next.
		// While we have complete nodes, dig deeper for a new state to explore.
		for len(node.untriedMoves) == 0 && len(node.children) > 0 {
			// This node has no more moves to try but it does have children.
			// Move the focus to its most promising child.
			node = node.selectChild()
			simulatedState.MakeMove(node.move)
		}

		// 2. Expand.
		// Can we explore more about this particular state? Are there untried moves?
		if len(node.untriedMoves) > 0 {
			// This creates a new child node with cloned game state.
			// NOTE: node becomes the new child
			move := node.getRandomMoves()
			simulatedState.MakeMove(move)
			child := newTreeNode(node, move, simulatedState.Clone(), node.ucbC)
			node.children = append(node.children, child)
			node = child
		}

		// 3. Simulation.
		// From the new child, make many simulated random steps to get a fuzzy idea of how good
		// the move that created the child is.
		for j := 0; j < int(opt.Simulations); j++ {
			// What moves can further the game state?
			availableMoves := simulatedState.AvailableMoves()
			// Is the game over?
			if len(availableMoves) == 0 {
				break
			}
			// Pick a random move (could be any player).
			randomIndex := rand.Intn(len(availableMoves))
			move := availableMoves[randomIndex]
			simulatedState.MakeMove(move)
		}

		// 4. Back propagate.
		// Our simulated state may be good or bad in the eyes of our player of interest.
		// todo no need to score for each layer of node
		node.update(simulatedState.Score) // Will internally propagate up the tree.
	}

	//root.printNode(0)

	// The best move to take is going to be the root nodes most visited child.
	sort.Sort(byVisits(root.children))
	return root.children[0].move // Descending by visits.
}

// upperConfidenceBound calculates the value of a child node (relative to its parent) for selection.
// c is a bias parameter, higher favors exploration (value children that have not been explored much),
// lower favors exploitation (value children for the scores they've already accumulated).
func upperConfidenceBound(
	childAggregateOutcome float64,
	ucbC float64,
	parentVisits uint64,
	childVisits uint64) float64 {
	return childAggregateOutcome/float64(childVisits) +
		ucbC*math.Sqrt(2*math.Log(float64(parentVisits))/float64(childVisits))
}
