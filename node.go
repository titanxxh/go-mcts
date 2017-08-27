package mcts

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

//rnd = rand.New(rand.NewSource(37))//fix seed

// A node in the (action, state) game tree. Wins are from the viewpoint of the player-just-moved.
type treeNode struct {
	parent         *treeNode   // What node contains this node? Root node's parent is nil.
	move           Move        // What move lead to this node? Root node's action is nil.
	state          GameState   // What is the game state at this node?
	totalOutcome   float64     // What is the sum of all outcomes computed for this node and its children? From the point of view of a single player.
	visits         uint64      // How many times has this node been studied? Used with totalValue to compute an average value for the node.
	untriedMoves   []Move      // What moves have not yet been explored from this state?
	children       []*treeNode // The children of this node, can be many.
	ucbC           float64     // The UCB constant used in selection calculations.
	selectionScore float64     // The computed score for this node used in selection, balanced between exploitation and exploration.
	player         int         // from whose perspective to score
}

// newTreeNode creates a new well-formed tree node.
func newTreeNode(parent *treeNode, move Move, state GameState, ucbC float64) *treeNode {
	// Construct the new node.
	node := treeNode{
		parent:         parent,
		move:           move,
		state:          state,
		totalOutcome:   0.0,                    // No outcome yet.
		visits:         0,                      // No visits yet.
		untriedMoves:   state.AvailableMoves(), // Initially the node starts with every node unexplored.
		children:       nil,                    // No children yet.
		ucbC:           ucbC,                   // Whole tree uses same constant.
		selectionScore: 0.0,                    // No value yet.
		player:         state.PlayerJustMoved(),
	}

	// We're working with pointers.
	return &node
}

// getVisits returns the visits to a node, 0 if the node doesn't exist (for when a root checks its parent).
func (n *treeNode) getVisits() uint64 {
	if n == nil {
		return 0
	}
	return n.visits
}

// computeSelectionScore prepares the selection score of a single child.
func (n *treeNode) computeSelectionScore() {
	n.selectionScore = upperConfidenceBound(n.totalOutcome, n.ucbC, n.parent.getVisits(), n.visits)
}

// selectChild picks the child with the highest selection score (balancing exploration and exploitation).
func (n *treeNode) selectChild() *treeNode {
	// Sort the children by their UCB, balances winning children with unexplored children.
	sort.Sort(bySelectionScore(n.children))
	return n.children[0]
}

// update adds the outcome value from a computation involving the node or one of its children.
// Every outcome value in the tree is from the perspective of a particular player. Higher outcomes mean better
// winning situations for the player.
func (n *treeNode) update(score func(player int) float64) {
	// Allow the root to call this on its parent with no ill effect.
	if n != nil {
		outcome := score(n.player)
		// Update this node's data.
		n.totalOutcome += outcome
		n.visits++
		// Pass the value up to the parent as well.
		n.parent.update(score) // Will recurse up the tree to the root.
		// Now that the parent is also updated
		n.computeSelectionScore()
	}
}

// get and remove a random optional move
func (n *treeNode) getRandomMoves() Move {
	i := rnd.Intn(len(n.untriedMoves))
	move := n.untriedMoves[i]
	n.untriedMoves = append(n.untriedMoves[:i], n.untriedMoves[i+1:]...)
	return move
}

func (n *treeNode) printNode(depth int) {
	fmt.Println(fmt.Sprintf("d%d>p%d move:%v oc%f vst%d", depth, n.player, n.move, n.totalOutcome, n.visits))
	if len(n.children) > 0 {
		for _, child := range n.children {
			child.printNode(depth + 1)
		}
	}
}

// bySelectionScore implements sort.Interface to sort *descending* by selection score.
// Example: sort.Sort(bySelectionScore(nodes))
type bySelectionScore []*treeNode

func (a bySelectionScore) Len() int           { return len(a) }
func (a bySelectionScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bySelectionScore) Less(i, j int) bool { return a[i].selectionScore > a[j].selectionScore }

// byVisits implements sort.Interface to sort *descending* by visits.
// Example: sort.Sort(byVisits(nodes))
type byVisits []*treeNode

func (a byVisits) Len() int           { return len(a) }
func (a byVisits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byVisits) Less(i, j int) bool { return a[i].visits > a[j].visits }
