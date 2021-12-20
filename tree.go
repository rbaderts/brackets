package brackets

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

type GameState struct {
	Result *MatchResult `json:"result"`
}

type NodeId int

type MatchResult struct {
	WinningNode NodeId `json:"winningNode"`
	WinningSlot int    `json:"winningSlot"`

	WinningPlayer int `json:"winningPlayer"`
	//	LosingPlayer   int             `json:"losingPlayer"`

	LosingParticipant  ParticipantNumber `json:"losingParticipant"`
	WinningParticipant ParticipantNumber `json:"winningParticipant"`

	DropNode NodeId `json:"dropNode"`

	Time int64 `json:"time"`
}

type NodeReference struct {
	Id   NodeId `json:"id"`
	Kind int    `json:"kind"`
	node *Node
}

/**
  A Bracket is represented as a tree of Nodes.   There are 3 nodes types (Player, Game and Drop)
    *  Each node has an ID unique withing the tournament
    *  Nodes maintain references to both their children and to their parent.
	*  The Node at the root of the tree represents the final game of the tournament
    *  Leaf nodes are associated with a single player (either at tournament start for Player nodes, or after
     	a "drop" occurs for Drop nodes).
    *  All non-leaf nodes represent Games.
	* A Nodes Level represents it's level in the tree where the root is Level 1.
	* A Nodes Depth reprents its its longest subtree,  So the depth of the
	     root represents the longest branch.

	Games:

       Game nodes have 2 slots, one for each participant.   When a game is resolved it's partiipant field references
	      the winning player.    A game can be played when both it's children nodes (either Game nodes, or other),
		  have a non-0 "participant" field. (Meaning it is in a resolved state)


    *  For winners side leaf nodes ("Player" nodes) the player is assigned before play begins
    *  For losers side leaf nodes ("Drop" nodes) the player is the loser of one of the winners side matches
    *  At the completion of a game the winning participant is propogated to on of the "slots" in the parent game node.

    * Neither Drop nor Player nodes are ever parents of other nodes, they always are leaves.

	Drops:

	     A drop is when a player loses their first winner side game and "drop" down to the losers bracket
		 A drop has a source (Game) node and a target "Drop" node.    The "Drop" field of a Node references
		    its associated source or target...

*/

type NodeType int

const (
	PLAYER = 1
	GAME   = 2
	DROP   = 3
	FINAL  = 4
)

func (this NodeType) String() string {
	switch this {
	case PLAYER:
		return "Player"
	case GAME:
		return "Game"
	case DROP:
		return "Drop"
	case FINAL:
		return "Final"

	}
	return "Unknown"

}

type NodeSubType int

const (
	NONE         = 1
	WINNERS_ROOT = 2
	LOSERS_ROOT  = 3
	FIRST_FINAL  = 4
	SECOND_FINAL = 5
	ZOMBIE       = 6 // Game where one of th participants is a buy
)

type TreeContext struct {
	IdCounter          NodeId `json:"idCounter"`
	InvisibleIdCounter NodeId `json:"invisibleIdCounter"`
}

func NewTreeContext() *TreeContext {
	t := new(TreeContext)
	t.IdCounter = NodeId(1)
	t.InvisibleIdCounter = NodeId(1000)
	return t
}

func (this *TreeContext) IncrementCount() {
	this.IdCounter = NodeId(int(this.IdCounter) + 1)
}

type Node struct {
	Id      NodeId        `json:"id"`
	Type    NodeType      `json:"nodeType"`
	SubType NodeSubType   `json:"nodeSubType"`
	Left    NodeReference `json:"left"`
	Right   NodeReference `json:"right"`
	Parent  NodeReference `json:"parent"`

	Span      Span      `json:"span"`
	GridSpan  Span      `json:"gridSpan"`
	GameState GameState `json:"state"`

	Depth int `json:"depth"`
	Level int `json:"level"`
	Tier  int `json:"tier"`

	// For winners side (GAME) nodes this is what node receives the loser
	// For losers side (DROP) nodes it is the winners node the drop came from
	Drop NodeId `json:"drop"`

	// For GAME nodes Participant is the ParticipantNumber of
	//   the game winner, 0 if there is none
	// For Player nodes this is the ParticipantNumber
	// For Drop nodes this is the ParticipantNumber of the
	//   player the dropped , 0 if there isn't one yet
	Participant ParticipantNumber `json:"participant"`

	Incumbant ParticipantNumber

	ChallengerUpOne bool `json:"challengerUpOne"`
}

func (this *MatchResult) String() string {
	return fmt.Sprintf("MatchResult:  WinningNode = %d:  WinningSlot = %d, WinningParticipant = %d\n",
		this.WinningNode, this.WinningSlot, this.WinningParticipant)
}

func (this *Node) CalculateDepth() int {

	if this.Left.Id == 0 && this.Right.Id == 0 {
		this.Depth = 0
		return this.Depth
	}

	leftDepth := 0
	rightDepth := 0
	if this.Left.Id != 0 {
		leftDepth = this.Left.node.CalculateDepth()
	}
	if this.Right.Id != 0 {
		rightDepth = this.Right.node.CalculateDepth()
	}

	if leftDepth > rightDepth {
		this.Depth = leftDepth + 1
	}
	if rightDepth > leftDepth {
		this.Depth = rightDepth + 1
	}
	return this.Depth

}

func (this *Node) calcLevel() int {
	fmt.Printf("calcLevel id: %d\n", this.Id)

	if this.Level != 0 {
		fmt.Printf("calcLevel this.Level = %d\n", this.Level)
		return this.Level
	}

	if this.Parent.Id == 0 {
		fmt.Printf("calcLevel return 1\n")
		return 1
	}

	// cache the calculaated level in the node
	parentLevel := this.Parent.node.calcLevel()
	this.Level = parentLevel + 1
	fmt.Printf("calcLevel cached Level %d, for id: %d\n", this.Level, this.Id)
	return this.Level

}

func (this *Node) internalize(b *Bracket) {
	this.Left.node = b.GetNode(this.Left.Id)
	this.Right.node = b.GetNode(this.Right.Id)
	this.Parent.node = b.GetNode(this.Parent.Id)
}

func (this *Node) PrintTree() string {
	return this.printNode(new(bytes.Buffer), true, new(bytes.Buffer)).String()
}

func (this *Node) debugNode() {

	fmt.Printf("Id=%d, Right=(%d,kind=%d), LeftId=(%d,kind=%d)\n", this.Id, this.Right.Id, this.Right.Kind, this.Left.Id, this.Left.Kind)

}

func (this *Node) printNode(prefix *bytes.Buffer, isTail bool, buf *bytes.Buffer) *bytes.Buffer {

	if this.Right.Id != 0 {
		b := new(bytes.Buffer)
		p := "|    "
		if isTail {
			p = "     "
		}

		b.Write(prefix.Bytes())
		b.WriteString(p)
		if this.Right.node != nil {
			this.Right.node.printNode(b, false, buf)
		}
	}

	t := "└── "
	if !isTail {
		t = "┌── "
	}
	buf.WriteString(prefix.String())
	buf.WriteString(t)
	buf.WriteString(fmt.Sprintf(Red+"%d"+Reset+" (l:%d, d: %d,, pl:%d typ:%v D:%d\n",
		this.Id, this.Level, this.Depth, this.Participant, this.Type, this.Drop))

	if this.Left.Id != 0 {
		b := new(bytes.Buffer)
		p := "|    "
		if isTail {
			p = "     "
		}

		b.Write(prefix.Bytes())
		b.WriteString(p)
		if this.Left.node != nil {
			this.Left.node.printNode(b, true, buf)
		}
	}
	return buf
}

/*
func (this Node) String() string {
	return this.PrintTree()
}
*/

func (this Node) String() string {
	//    var b strings.Builder

	rightId := 0
	rightKind := 0
	leftId := 0
	leftKind := 0
	if true {
		rightId = int(this.Right.Id)
		rightKind = this.Right.Kind
	}
	if true {
		leftId = int(this.Left.Id)
		leftKind = this.Left.Kind
	}

	pt1 := fmt.Sprintf("Id=%d, Depth: %d, Leel: %d, Parent: %d, Right:(%d,kind=%d), Left:(%d,kind=%d) Drop: %d\n",
		this.Id, this.Depth, this.Level, this.Parent.Id, rightId, rightKind, leftId, leftKind, this.Drop)

	var pt2 string
	if this.GameState.Result != nil {
		pt2 = fmt.Sprintf("   %v\n", this.GameState.Result)
	} else {
		pt2 = fmt.Sprintf("   No Game Result\n")
	}

	return fmt.Sprintf("%s%s\n", pt1, pt2)
}

func (this *Node) GetInnerNodes() []NodeId {
	tmpinners := make([]NodeId, 0)
	this.innerNodes(&tmpinners)
	sort.Slice(tmpinners, func(i, j int) bool { return int(i) < int(j) })

	return tmpinners
}

func (this *Node) GetLeafNodes() []NodeId {
	tmpleafs := make([]NodeId, 0)
	this.leafNodes(&tmpleafs)
	sort.Slice(tmpleafs, func(i, j int) bool { return int(i) < int(j) })
	return tmpleafs
}

func (this *Node) leafNodes(nodes *[]NodeId) {
	if this.Right.Id != 0 {
		this.Right.node.leafNodes(nodes)
	}
	if this.Left.Id != 0 {
		this.Left.node.leafNodes(nodes)
	}

	if this.Left.Id == 0 && this.Right.Id == 0 {
		(*nodes) = append((*nodes), this.Id)
	}
}
func (this *Node) innerNodes(nodes *[]NodeId) {

	if this.Right.Id == 0 && this.Left.Id == 0 {
		return
	} else {

		fmt.Printf("adding node %d level %d\n", this.Id, this.Level)
		(*nodes) = append((*nodes), this.Id)
		//inners = append(inners, this.Id)

		if this.Right.Id != 0 {
			this.Right.node.innerNodes(nodes)
		}
		if this.Left.Id != 0 {
			this.Left.node.innerNodes(nodes)
		}
	}
}

func (this *Node) SetParent(parent *Node) {
	this.Parent = NodeReference{parent.Id, 1, parent}
	fmt.Printf("Setting parent of Node %d of type %s to node %d\n", this.Id,
		this.Type, parent.Id)

	/*`
	if (parent.Level >= 0) {
		this.Level = parent.Level - 1
	}
	*/
	//	this.Level = parent.Level + 1
}

func (this *Node) SetParticipant(p ParticipantNumber) {
	this.Participant = p
	this.UpdateParent()
}

func (this *Node) UpdateParent() {
	if this.Parent.node.Left.node.Participant == -1 ||
		this.Parent.node.Right.node.Participant == -1 {
		this.Parent.node.SubType = ZOMBIE
	} else {
		this.Parent.node.SubType = NONE
	}

}

func (this *Node) SetLeft(left *Node) {
	kind := left.Type
	this.Left = NodeReference{left.Id, int(kind), left}

	if this.Level != 0 {
		left.Level = this.Level + 1
	}
	//	this.Level = left.Level - 1

	if left.Depth > this.Depth-1 {
		this.Depth = left.Depth + 1
	}

	fmt.Printf("Setting Left child of node %d to node %d\n", this.Id, left.Id)
	left.SetParent(this)
}

func (this *Node) SetRight(right *Node) {
	kind := right.Type
	this.Right = NodeReference{right.Id, int(kind), right}
	//this.Level = right.Level - 1
	if this.Level != 0 {
		right.Level = this.Level + 1
	}
	if right.Depth > this.Depth-1 {
		this.Depth = right.Depth + 1
	}
	fmt.Printf("Setting Right child of node %d to node %d\n", this.Id, right.Id)
	right.SetParent(this)
}

func (this *TreeContext) NewNode(t NodeType, parent *Node, dropGameId NodeId, participantNumber ParticipantNumber, tier int, level int) *Node {

	parentId := -1
	if parent != nil {
		parentId = int(parent.Id)
	}

	id := this.IdCounter
	if t != GAME {
		id = this.InvisibleIdCounter
		this.InvisibleIdCounter += 1
	} else {
		this.IdCounter += 1

	}
	fmt.Printf("NewNode: Id: %d, type: %v, parentId: %d, drop: %d, level: %d, participant: %d\n",
		id, t, parentId, dropGameId, level, participantNumber)

	node := new(Node)
	node.Type = t
	node.SubType = NONE
	node.Tier = tier
	node.Drop = dropGameId
	//node.Player = playerId
	node.Participant = participantNumber

	node.Id = id
	node.Left = NodeReference{0, 0, nil}
	node.Right = NodeReference{0, 0, nil}

	if parent != nil {
		node.Parent = NodeReference{parent.Id, int(parent.Type), parent}
	} else {
		//		node.Level = 0
	}

	node.GridSpan.Upper = 1
	node.GridSpan.Lower = 1
	node.Span.Upper = 10
	node.Span.Lower = 10
	node.Level = level

	return node

}

/*

 */
func (this *Node) calculateGridSpans() Span {

	if this.Left.Id != 0 {
		if this.Left.Kind == GAME {
			leftSpan := this.Left.node.calculateGridSpans()
			this.GridSpan.Upper += int(float32(leftSpan.Upper+leftSpan.Lower)) - 1
		} else {
			this.GridSpan.Upper = 1
		}
	} else {
		this.GridSpan.Upper = 1
	}
	if this.Right.Id != 0 {
		if this.Right.Kind == GAME {
			rightSpan := this.Right.node.calculateGridSpans()
			this.GridSpan.Lower += int(float32(rightSpan.Upper+rightSpan.Lower)) - 1
		} else {
			this.GridSpan.Lower = 1
		}
	} else {
		this.GridSpan.Lower = 1
	}
	fmt.Printf("nodeId: %d, gridspan = %v\n", this.Id, this.GridSpan)
	return this.GridSpan
}

func (this *Node) calculateSpans() Span {

	this.Span = Span{}
	if this.Left.Id != 0 {
		if this.Left.Kind == GAME {
			span := this.Left.node.calculateSpans()
			this.Span.Upper += (span.Upper + span.Lower)
		} else if this.Left.Kind == PLAYER ||
			this.Left.Kind == DROP {
			this.Span.Upper = 24
		} else {
			//			this.Span.Upper = 12
			this.Span.Upper = 20
		}
	} else {
		this.Span.Upper = 20
	}

	if this.Right.Id != 0 {
		if this.Right.Kind == GAME {
			span := this.Right.node.calculateSpans()
			this.Span.Lower += (span.Upper + span.Lower)
		} else if this.Right.Kind == PLAYER ||
			this.Right.Kind == DROP {
			this.Span.Lower = 24
		} else {
			//			this.Span.Lower = 12
			this.Span.Lower = 20
		}

	} else {
		this.Span.Lower = 20
	}

	fmt.Printf("nodeId: %d, span = %v\n", this.Id, this.Span)

	return this.Span
}

type NodeQueue []*Node

func (s *NodeQueue) Size() int {
	return len(*s)
}

func (s *NodeQueue) Add(v *Node) {
	*s = append(*s, v)
}

func (s *NodeQueue) AddAll(other *NodeQueue) {
	for _, v := range *other {
		*s = append(*s, v)
	}
}

func (s *NodeQueue) Remove() *Node {
	res := (*s)[0]
	*s = (*s)[1:len(*s)]
	return res
}

func (s *NodeQueue) RemoveHead() *Node {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *NodeQueue) RemoveAt(i int) *Node {
	res := (*s)[i]

	copy((*s)[i:], (*s)[i+1:]) // Shift a[i+1:] left one index.

	(*s)[len(*s)-1] = nil // Erase last element (write zero value).
	*s = (*s)[:len(*s)-1] // Truncate slice.

	return res

}

func (this NodeQueue) String() string {
	var s strings.Builder

	for _, r := range []*Node(this) {
		fmt.Fprintf(&s, "%d ", r.Id)
	}
	fmt.Fprintf(&s, "\n")
	return s.String()

}

func (s *NodeQueue) NodeIds() []NodeId {
	ids := make([]NodeId, 0)
	for _, n := range []*Node(*s) {
		ids = append(ids, n.Id)
	}

	return ids

}
