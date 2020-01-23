package brackets

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type GameState struct {
	Result       *MatchResult   `json:"result"`
}

var IdCounter = 0

type MatchResult struct {
	WinningNode    int             `json:"winningNode"`
	WinningSlot    int             `json:"winningSlot"`
	WinningPlayer  int             `json:"winningPlayer"`

	LosingPlayer   int             `json:"losingPlayer"`
	DropNode       int             `json:"dropNode"`

	Time           int64           `json:"time"`
}

type NodeReference struct {
	Id          int               `json:"id"`
	Kind        int               `json:"kind"`
	node        *Node
}

/**
    Each node has a unique Id.

    Leaf nodes are associated with a single player (either at tournament start, or at a drop)
    All other nodes represent a game.    3 nodes types (Player, Game and Drop)
    FOr winners side leaf nodes the player is assigned before play begins
    For losers side leaf nodes the player is the loser of one of the winners side matches
    For all others the player is the winner of a match between the nodes left and right children nodes
 */


type NodeType int

const (
	PLAYER     = 1
	GAME       = 2
	DROP       = 3
)


type Node struct {
	Id          int                 `json:"id"`
	Type        NodeType            `json:"nodeType"`
	Left        NodeReference       `json:"left"`
	Right       NodeReference       `json:"right"`
	Parent      NodeReference       `json:"parent"`

	Span        Span               	`json:"span"`
	GameState   GameState			`json:"state"`

	Tier        int                 `json:"tier"`       // 1 for winners 2 for losers
	Level       int					`json:"level"`

	Drop        int                 `json:"drop"`       // what game drops to here

	Player     int	    			`json:"player"`

}


func (this *Node) internalize(t *Tournament) {
	this.Left.node = t.GetNode(this.Left.Id)
	this.Right.node = t.GetNode(this.Right.Id)
	this.Parent.node = t.GetNode(this.Parent.Id)
}

func (this *Node) PrintTree() string {
	return this.printNode(new(bytes.Buffer), true, new(bytes.Buffer)).String()
}

func (this *Node) printNode(prefix *bytes.Buffer, isTail bool, buf *bytes.Buffer) *bytes.Buffer {

	if( this.Right.Id != 0) {
		b := new(bytes.Buffer)
		p := "|    "
		if isTail {
			p = "     "
		}

		b.Write(prefix.Bytes())
		b.WriteString(p)
		if (this.Right.node != nil) {
			this.Right.node.printNode(b, false, buf)
		}
	}

	t := "└── "
	if ! isTail {
		t = "┌── "
	}
	buf.WriteString(prefix.String())
	buf.WriteString(t)
	buf.WriteString(fmt.Sprintf("%d(lvl:%d, player:%d \n", this.Id, this.Level, this.Player))

	if (this.Left.Id != 0) {
		b := new(bytes.Buffer)
		p := "|    "
		if isTail {
			p = "     "
		}

		b.Write(prefix.Bytes())
		b.WriteString(p)
		if (this.Left.node != nil) {
		    this.Left.node.printNode(b, true, buf)
		}
	}
	return buf;
}

/*
func (this Node) String() string {
	return this.PrintTree()
}
 */

func (this Node) String() string {
//    var b strings.Builder

    return fmt.Sprintf("Node %d:  Left: %d, Right: %d, Drop: %d, ",
    	this.Id, this.Left.Id, this.Right.Id, this.Drop)

}

func (this *Node) GetInnerNodes() []int {
	tmpinners := make([]int, 0)
	this.innerNodes(&tmpinners)
	sort.Ints(tmpinners)
	return tmpinners
}

func (this *Node) innerNodes(nodes *[]int) {

	if (this.Right.Id == 0 && this.Left.Id == 0) {
		return
	} else {

		fmt.Printf("adding node %d level %d\n", this.Id, this.Level)
		(*nodes) = append((*nodes), this.Id)
		//inners = append(inners, this.Id)

		if (this.Right.Id != 0) {
			this.Right.node.innerNodes(nodes)
		}
		if (this.Left.Id != 0) {
			this.Left.node.innerNodes(nodes)
		}
	}
}

func (this *Node) SetParent(parent *Node) {
	this.Parent = NodeReference{parent.Id, 1, parent}
	/*`
	if (parent.Level >= 0) {
		this.Level = parent.Level - 1
	}
	 */
//	this.Level = parent.Level + 1
}

func (this *Node) SetLeft(left *Node) {
	kind := 1
	if left.Type == PLAYER || left.Type == DROP {
	     kind = 2
	}
	this.Left = NodeReference{left.Id, kind, left}
	this.Level = left.Level - 1
	left.SetParent(this)
}

func (this *Node) SetRight(right *Node) {
	kind := 1
	if right.Type == PLAYER || right.Type == DROP {
		kind = 2
	}
	this.Right = NodeReference{right.Id, kind, right}
	this.Level = right.Level - 1
	right.SetParent(this)
}


func NewNode(t NodeType,  parent *Node, dropGameId int, playerId int, tier int, level int) *Node {

	node := new(Node)
	node.Type = t
	node.Tier = tier
	node.Drop = dropGameId
    node.Player = playerId
	node.Id = IdCounter
	IdCounter += 1
	node.Left = NodeReference{0 , 0, nil}
	node.Right = NodeReference{0 , 0, nil}

	if (parent != nil) {
		node.Parent = NodeReference{parent.Id, 1, parent}
	} else {
//		node.Level = 0
	}

    node.Span.Upper = 10
	node.Span.Lower = 10
	node.Level = level

	return node

}

func (this *Node) calculateSpans () Span {

	if (this.Left.Id != 0) {
		if (this.Left.Kind  == 1) {
			span := this.Left.node.calculateSpans()
			this.Span.Upper += (span.Upper + span.Lower) - 10
		} else if (this.Left.Kind == 2) {
			this.Span.Upper = 24
		} else {
//			this.Span.Upper = 12
			this.Span.Upper = 20
		}
	} else {
		this.Span.Upper = 20
	}


	if (this.Right.Id != 0) {
		if (this.Right.Kind  == 1) {
			span := this.Right.node.calculateSpans()
			this.Span.Lower += (span.Upper + span.Lower) - 10
		} else if (this.Right.Kind == 2) {
			this.Span.Lower = 24
		} else {
//			this.Span.Lower = 12
			this.Span.Lower = 20
		}

	} else {
		this.Span.Lower = 20
	}

	fmt.Printf("nodeId: %d, span = %d %d\n", this.Id, this.Span.Upper, this.Span.Lower)

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

	(*s)[len(*s)-1] = nil   // Erase last element (write zero value).
	*s = (*s)[:len(*s)-1] // Truncate slice.

	return res

}

func (this NodeQueue) String() string {
	var s strings.Builder

	for _, r := range ([]*Node(this)){
		fmt.Fprintf(&s, "%d ", r.Id)
	}
	fmt.Fprintf(&s, "\n")
	return s.String()

}


