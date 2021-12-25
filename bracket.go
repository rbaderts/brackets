package brackets

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
)

type Bracket struct {
	Root       *Node             `json:"root"`
	RootNodeId NodeId            `json:"rootNodeId"`
	Drops      map[NodeId]NodeId `json:"drops"` // map of Node Id to Node Id

	Nodes map[NodeId]*Node `json:"nodes"`

	Size int `json:"size"`

	//	Participants []ParticipantNumber `json:"participants"`

	//	WinnerSideRootNode NodeId `json:"winnersSideRootNode"`
	//	LosersSideRootNode NodeId `json:"losersSideRootNode"`

	Degree       int          `json:"degree"`
	Buys         int          `json:"buys"`
	Depth        int          `json:"depth"`
	WinnersDepth int          `json:"winnersDepth"`
	LosersDepth  int          `json:"losersDepth"`
	Context      *TreeContext `json:"treeContext"`
	tournament   *Tournament
}

func (this *Bracket) RecordResult(node *Node, winnerSide int) error {

	if node.Left.node.Participant == 0 || node.Right.node.Participant == 0 {
		return errors.New("Game can't be completed, state error")
	}

	if node.Type == FINAL && winnerSide == 2 {
		if node.ChallengerUpOne == false {
			node.ChallengerUpOne = true
			return nil
		}
	}

	//	var winnerNodeId NodeId
	var winner ParticipantNumber
	var loser ParticipantNumber
	if winnerSide == 1 {
		winner = node.Left.node.Participant
		//		winnerNodeId = node.Left.Id
		loser = node.Right.node.Participant
	} else {
		winner = node.Right.node.Participant
		//	winnerNodeId = node.Right.Id
		loser = node.Left.node.Participant
	}
	//winnerPlayerId = this.Participants[winner].PlayerId
	//loserPlayerId = this.Participants[loser].PlayerId

	//fmt.Printf("setting node (%d) player to %d\n", node.Id, winner)
	node.Participant = winner

	var dropNode *Node = nil
	dropNodeId, has := this.Drops[node.Id]
	if has {
		dropNode = this.GetNode(dropNodeId)
	}
	if dropNode != nil {
		//fmt.Printf("dropNode id = %d, setting player to %d\n", dropNode.Id, loser)
		dropNode.Participant = loser
	}

	mr := new(MatchResult)
	mr.LosingParticipant = -1
	node.GameState.Result = mr

	return nil

}

/*
   Meanings of Var:

       For tree.go structure association (Left and Right) Var means:
           0 - means its a direct parent child relationship
           1 - means a association from a Loser tier game to a winners tier game (a drop down)
       For associations from a game to the games the winner and loser propagate to (WinnersGame and LosersGame)
       Var means:

           The slot # in the target game the player goes to.
*/

func (this *Bracket) findWinnersBuyNode() *Node {

	winnersRoot := this.Root.Left.node
	leaves := winnersRoot.GetLeafNodes()
	for _, leaf := range leaves {
		node := this.Nodes[leaf]
		if node.Level == this.Depth+1 {
			return node
		}
	}
	return nil
}

func (this *Bracket) findLosersBuyNode() *Node {

	winnersRoot := this.Root.Left.node
	leaves := winnersRoot.GetLeafNodes()

	for _, leaf := range leaves {
		node := this.Nodes[leaf]

		if node.Level == this.Depth+1 {
			return node
		}
	}

	return nil

}

func (this *Bracket) AddPlayer(pNumber ParticipantNumber) error {

	buyNode := this.findWinnersBuyNode()
	losersBuyNode := this.findLosersBuyNode()
	if buyNode == nil {
		return errors.New("No buy nodes")
	}

	losersNode := this.newNode(GAME, losersBuyNode, 0, 0, 2, losersBuyNode.Level+1)
	right := this.newNode(GAME, losersBuyNode, 0, 0, 2, losersBuyNode.Level+1)
	left := this.newNode(GAME, buyNode, losersBuyNode.Id, 0, 2, buyNode.Level+1)

	_ = losersNode
	_ = right
	_ = left
	newNode := this.newNode(PLAYER, nil, 0, ParticipantNumber(pNumber), 1, buyNode.Level+1)
	_ = newNode

	//	this.newNode()/new
	//this.Participants = append(this.Participants, pNumber)

	return nil

}

func (this *Bracket) deleteNode(id NodeId) {
	delete(this.Nodes, id)
}

func (this *Bracket) AddNode(node *Node) {

	this.Nodes[node.Id] = node

}

func (this *Bracket) newNode(t NodeType, parent *Node, dropGameId NodeId, participantNumber ParticipantNumber, tier int, level int) *Node {
	node := this.Context.NewNode(t, parent, dropGameId, participantNumber, tier, level)
	this.Nodes[node.Id] = node
	if level > this.Depth {
		this.Depth = this.Depth + 1
	}

	return node
}

func (this *Bracket) NodesString() string {
	var b strings.Builder
	a := make([]NodeId, len(this.Nodes))
	for i, n := range this.Nodes {
		a[i-1] = n.Id
	}

	sort.Slice(a, func(i, j int) bool { return int(i) < int(j) })

	fmt.Fprintf(&b, "%v ", a)
	fmt.Fprintf(&b, "\n")
	return b.String()
}

func (this *Bracket) internalize() {
	for _, n := range this.Nodes {
		n.internalize(this)
	}
	//	this.Root = this.GetNode(this.RootNodeId)

	//	this.Context = NewTreeContext()
	//	this.context.IdCounter = this.RootNodeId + 1
}

func (this *Bracket) GetNode(id NodeId) *Node {
	node := this.Nodes[id]
	return node
}

func (this *Bracket) SetNodeDropById(nodeId NodeId, dropNodeId NodeId) {
	node := this.GetNode(nodeId)
	this.SetNodeDrop(node, dropNodeId)
}

func (this *Bracket) SetNodeDrop(node *Node, dropNodeId NodeId) {

	node.Drop = dropNodeId
	dropNode := this.GetNode(dropNodeId)
	dropNode.Drop = node.Id
	Logger.Infof("Setting Node %d to drop to %d\n", node.Id, dropNodeId)

	Logger.Infof("Setting Node %d as recipient of drop from %d\n", dropNodeId, node.Id)

}

// The # of total games in the Winners bracket

func ComputeWTierTotal(players int) int {
	size := ComputeSize(players)
	degree := ComputeDegree(size)
	tier1 := ComputeWTier1(size, players)
	total := (2 << (degree - 2)) + tier1 - 1
	return total
}

func ComputeLTierTotal(players int) int {

	size := ComputeSize(players)
	degree := ComputeDegree(size)
	tier1 := ComputeLTier1(size, players)
	tier2 := ComputeLTier2(size, players)
	total := (2 << (degree - 2)) - 2
	total += tier1 + tier2
	return total

}

func ComputeWTier1(size int, players int) int {
	if players <= (size >> 1) {
		log.Fatalf("bad size:  size: %d, players; %d\n", size, players)
	}
	num := (size >> 1) - (size - players)
	return num
}

func ComputeWTier2(size int, players int) int {
	//vm := int(math.Min(float64((size)-(size-players)), float64(size>>2)))
	//fmt.Printf("vm = %d\n", vm)
	return size >> 2
}

func ComputeLTier1(size int, players int) int {
	v := int(math.Max(float64(0), float64((size>>2)-(size-players))))
	return v
}

func ComputeLTier2(size int, players int) int {
	vm := int(math.Min(float64((size>>1)-(size-players)), float64(size>>2)))
	return vm
}

// Game number calculations

func ComputeDegree(size int) int {

	s := size
	d := 1
	for {
		if s == 0 {
			//fmt.Printf("degree for size %d  = %d\n", size, d)
			return d
		}
		d += 1
		s = s >> 2
	}
}
func ComputeSize(players int) int {
	size := 2
	for {
		if players <= size {
			break
		}
		size = size << 1
	}
	return size
}

/*
   Given a set of sets of ids, where each set are the ids introduced at a given level, generate the tree.

    The tree is generated by creating nodes for each new id, adding them to the queue, then poping nodes 2 at a time
    from the Queue, create a new node from them  and put it back in the queue, the last node is the tree root.

    When there are no more ids to introduce finish processing the Queue to generate games

*/

type IntQueue []int

func (s *IntQueue) Add(v int) {
	*s = append(*s, v)
}

func (s *IntQueue) Remove() int {
	res := (*s)[0]
	*s = (*s)[1:len(*s)]
	return res
}

func (s *IntQueue) RemoveHead() int {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *IntQueue) RemoveAt(i int) int {
	res := (*s)[i]

	copy((*s)[i:], (*s)[i+1:]) // Shift a[i+1:] left one index.

	(*s)[len(*s)-1] = 0   // Erase last element (write zero value).
	*s = (*s)[:len(*s)-1] // Truncate slice.

	return res

}

type Span struct {
	Upper int `json:"upper"`
	Lower int `json:"lower"`
}

func (s Span) String() string {
	return fmt.Sprintf("{upper:%d,lower:%d}", s.Upper, s.Lower)
}

func NewSpan(upper int, lower int) *Span {
	span := new(Span)
	span.Upper = upper
	span.Lower = lower
	return span
}
