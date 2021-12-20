package brackets

import (
	"fmt"
	"math/bits"

	"github.com/gocraft/dbr/v2"
)

func (this *Bracket) GenerateLosersGames(dropQ *NodeQueue, gameQ *NodeQueue, level int, count int, useDrops bool) (*NodeQueue, []NodeId) {

	fmt.Printf("GenerateLosersGames: level = %d, dropQ.size = %d, gameQ.size = %d, count = %d, useDrops = %v\n",
		level, dropQ.Size(), gameQ.Size(), count, useDrops)
	newQ := new(NodeQueue)
	drops := make([]NodeId, 0)
	c := 0

	dropsToUse := (count * 2) - gameQ.Size()
	if useDrops == false {
		dropsToUse = 0
	}

	for {
		fmt.Printf("c = %d\n", c)
		if c == count {
			break
		}

		var n1 *Node
		var n2 *Node
		if gameQ.Size() > 0 {
			n1 = gameQ.Remove()
		} else {
			n1 = dropQ.Remove()
		}
		if dropsToUse > 0 {
			n2 = dropQ.Remove()
			dropsToUse -= 1
		} else {
			n2 = gameQ.Remove()
		}
		if n1.Type == DROP {
			drops = append(drops, n1.Id)
		}
		if n2.Type == DROP {
			drops = append(drops, n2.Id)
		}

		gameNode := this.newNode(GAME, nil, 0, 0, 2, 0)
		gameNode.SetRight(n1)
		gameNode.SetLeft(n2)
		newQ.Add(gameNode)
		c = c + 1
	}
	fmt.Printf("New Games: %d\n", newQ.Size())

	return newQ, drops

}

func (this *Bracket) GenerateWinnersGames(playerQ *NodeQueue, gameQ *NodeQueue, level int, max int) *NodeQueue {

	fmt.Printf("GenerateWinnersGames: level = %d, playerQ.size = %d, gameQ.size = %d, max = %d\n", level, playerQ.Size(), gameQ.Size(), max)
	newQ := new(NodeQueue)

	gameCount := 0

	numPlayerSources := 0
	if gameQ.Size() < max*2 {
		numPlayerSources = max*2 - gameQ.Size()
	}

	playerSlots := make([]int, max*2, max*2)
	fmt.Printf("len(playerSlots) = %d\n", max)
	for i := 0; i < numPlayerSources; i++ {
		index := 0
		if i%2 == 0 {
			index = i / 2
		} else {
			index = ((max * 2) - 1) - (i / 2)
		}
		fmt.Printf("palyerSlots[%d]\n", index)
		playerSlots[index] = 1

	}

	cnt := 0
	for {

		if gameCount == max {
			break
		}

		var n1 *Node
		var n2 *Node
		if playerSlots[cnt] == 1 {
			n1 = playerQ.Remove()
		} else {
			n1 = gameQ.Remove()
		}

		cnt = cnt + 1

		if playerSlots[cnt] == 1 {
			n2 = playerQ.Remove()
		} else {
			n2 = gameQ.Remove()
		}
		cnt = cnt + 1

		// func (this *Bracket) newNode(t NodeType, parent *Node, dropGameId NodeId, participantNumber ParticipantNumber, tier int, level int) *Node {
		gameNode := this.newNode(GAME, nil, 0, 0, 1, 0)
		gameNode.SetLeft(n1)
		gameNode.SetRight(n2)

		newQ.Add(gameNode)
		gameCount += 1
	}
	fmt.Printf("New Games: %d\n", newQ.Size())

	return newQ
}

func (this *Bracket) BuildDoubleElimBracket() *Node {

	root := this.GenerateEmptyDoubleElimBracket()

	this.CalculateLevels()

	fmt.Printf("\nAfter levels calculated:   %s\n", this.Root.PrintTree())

	fmt.Printf("\n%s\n", this.Root.PrintTree())

	root.calculateSpans()
	root.Right.node.Span.Upper += 20
	this.RootNodeId = root.Id

	return root
}

func (this *Bracket) CalculateLevels() {
	fmt.Printf("# leaaves: %d\n", len(this.Root.GetLeafNodes()))
	for _, v := range this.Root.GetLeafNodes() {
		this.GetNode(v).calcLevel()
	}

	this.Root.CalculateDepth()

}

func (this *Bracket) GenerateEmptyDoubleElimBracket() *Node {

	fmt.Print("GenerateEmptyDoubleElimBracket\n")

	level := this.Degree + 1
	this.WinnersDepth = level

	// Winners Side

	// Create all Player Nodes
	playerNodeQ := new(NodeQueue)

	for _, p := range this.Participants {
		newNode := this.newNode(PLAYER, nil, 0, ParticipantNumber(p), 1, 0)
		playerNodeQ.Add(newNode)
	}
	level = level - 1

	gameWidth := this.Size / 2
	firstLevelGames := gameWidth - this.Buys
	gameQ := new(NodeQueue)
	winnersSideGameNodes := make([]NodeId, 0)

	iteration := 1
	for {

		gameWidth := this.Size >> iteration
		maxGames := gameWidth
		if iteration == 1 {
			maxGames = firstLevelGames
		}
		fmt.Printf("WinnersGameLoop: itertion =%d, gameWidth = %d, maxGames = %d\n", iteration, gameWidth, maxGames)
		gameQ = this.GenerateWinnersGames(playerNodeQ, gameQ, 0, maxGames)

		winnersSideGameNodes = append(winnersSideGameNodes, gameQ.NodeIds()...)
		if gameQ.Size() == 1 && playerNodeQ.Size() == 0 {
			break
		}
		iteration += 1
	}

	fmt.Printf("GameNodes %v\n", winnersSideGameNodes)
	winnersRoot := gameQ.Remove()
	winnersRoot.SubType = WINNERS_ROOT
	winnersRoot.Level = 2

	// Losers Side

	// Create all drops
	dropQ := new(NodeQueue)
	for _, winnerNode := range winnersSideGameNodes {
		dropNode := this.newNode(DROP, nil, winnerNode, 0, 2, 0)
		this.SetNodeDropById(winnerNode, dropNode.Id)
		dropQ.Add(dropNode)
	}

	dropNodes := make([]NodeId, 0)

	extras := len(this.Participants) - this.Size/2
	fmt.Printf("extras = %d, loserFirstLevelGames = %d\n", extras, firstLevelGames)
	gameWidth = this.Size / 4
	if firstLevelGames > gameWidth {
		firstLevelGames = firstLevelGames - gameWidth
	}
	count := firstLevelGames
	fmt.Printf("count = %d\n", count)
	fmt.Printf("Size = %d\n", this.Size)
	fmt.Printf("Participants = %d\n", len(this.Participants))
	gameQ = new(NodeQueue)
	iteration = 1

	firstRoundGames := len(this.Participants) - (this.Size / 2)
	width := this.Size / 4

	useDrops := true
	for {

		//fmt.Printf("totalGames = %d, gamesLeft = %d, calcedGamesLeft=%d\n", totalGames, gamesLeft, calcedGamesLeft)
		//fmt.Printf("iteration: %d, count = %d, gamWidth = %d\n", iteration, count, gameWidth)
		var drops []NodeId

		if iteration == 1 {
			count = firstRoundGames % width
			firstRoundGames -= count
		} else if iteration == 2 {
			if firstRoundGames == width {
				count = width
			} else {
				count = width / 2
			}
		} else {
			if dropQ.Size() > gameQ.Size() {
				count = count
			} else {
				if count > 1 {
					count = count / 2
				}
				if gameQ.Size() == count && !(gameQ.Size() == 1 && dropQ.Size() == 1) {
					useDrops = false
				}
			}
		}
		fmt.Printf("iteration: %d, games = %d\n", iteration, count)
		gameQ, drops = this.GenerateLosersGames(dropQ, gameQ, 0, count, useDrops)
		useDrops = true
		//		gamesLeft -= gameQ.Size()
		//		if iteration == 1 {
		//			count = gameWidth
		//		} else {
		//			count = count / 2
		// }

		dropNodes = append(dropNodes, drops...)
		if gameQ.Size() == 1 && dropQ.Size() == 0 {
			break
		}

		iteration += 1
	}
	losersRoot := gameQ.Remove()
	losersRoot.SubType = LOSERS_ROOT
	losersRoot.Level = 2

	fmt.Printf("LosersRoot: \n%s\n", losersRoot.PrintTree())
	fmt.Printf("DropNdoes %v\n", dropNodes)

	root := new(Node)
	root.Id = this.Context.IdCounter
	this.Context.IdCounter += 1
	root.Tier = 1
	root.Level = 1
	root.Type = FINAL
	root.SubType = FIRST_FINAL
	this.Nodes[root.Id] = root
	root.SetLeft(winnersRoot)
	root.SetRight(losersRoot)

	this.LosersDepth = losersRoot.Depth
	this.WinnersDepth = winnersRoot.Depth

	this.Root = root
	this.RootNodeId = root.Id
	return this.Root
}

func (this *Bracket) computeVariables() {

	degree := 2
	size := 2

	for {
		if size >= len(this.Participants) {
			this.Size = size
			this.Degree = degree - 1
			this.Depth = this.Degree + 3
			break
		}
		size = size << 1
		degree += 1
	}
	this.Buys = this.Size - len(this.Participants)

	fmt.Printf("Participants = %d, Size = %d, Buys = %d, Degree = %d, Depth = %d\n",
		len(this.Participants), this.Size, this.Buys, this.Degree, this.Depth)

}

func calcGamesLeft(width int, drops int, hasDrops bool) int {
	w := width
	total := 0

	i := 0
	for {
		if i == 0 && hasDrops {
			drops -= w
		}
		if w == 0 {
			break
		}
		total += w * 2
		w = w >> 1
		i++
	}

	total += drops + 1
	return total / 2

}

func NewBracket(participants []ParticipantNumber) *Bracket {

	bracket := new(Bracket)
	bracket.Nodes = make(map[NodeId]*Node)
	bracket.Drops = make(map[NodeId]NodeId)
	bracket.Participants = participants

	bracket.Context = NewTreeContext()
	bracket.computeVariables()

	return bracket
}

func prevPowerOf2(num uint32) int {
	zeros := bits.LeadingZeros32(num)
	res := 1 << (31 - zeros)
	return res
}

func (this *Bracket) findLosersSpot(node *Node, winnersNode *Node) *Node {

	// Now find a cooresponding Loser side place to insert the new drop

	right := node.Right
	left := node.Left

	if (right.Kind == DROP && left.Kind == GAME) ||
		(right.Kind == GAME && left.Kind == DROP) {

		if node.Parent.node.GameState.Result == nil {
			return node
		}

	}

	if right.Kind == GAME {
		result := this.findLosersSpot(right.node, winnersNode)
		if result != nil {
			return result
		}
	}
	if left.Kind == GAME {
		result := this.findLosersSpot(left.node, winnersNode)
		if result != nil {
			return result
		}
	}

	return nil

}

func (this *Bracket) findAddLocation(node *Node) *Node {

	fmt.Printf("findAddLocation: node = %v\n", node)
	///	fmt.Printf("%v", node.PrintTree())
	right := node.Right
	left := node.Left

	if node.Type != GAME {
		return nil
	}

	fmt.Printf("findAddLocation: node.level = %d, WinnrsDepth = %d\n", node.Level, this.WinnersDepth)
	if node.Level >= this.WinnersDepth {
		return nil
	}

	if node.GameState.Result != nil {
		fmt.Printf("findAddLocation: Game already has a winner\n")
		return nil
	} else {
		fmt.Printf("Node has no match result")
	}

	if right.Kind == GAME && left.Kind == GAME {
		result := this.findAddLocation(right.node)
		if result != nil {
			return result
		}
		result = this.findAddLocation(left.node)
		if result != nil {
			return result
		}
	}

	if (right.Kind == PLAYER) && (left.Kind == PLAYER) {

		if node.GameState.Result == nil {
			fmt.Printf("Found spot to insert new user at Node: %v\n", node.Id)
			return node
		}

	}

	if right.Kind == GAME && left.Kind == PLAYER {
		result := this.findAddLocation(right.node)
		if result != nil {
			return result
		}
		if node.GameState.Result == nil {
			fmt.Printf("Found spot to insert new user at Node: %v\n", node.Id)
			return node
		}
	}
	if left.Kind == GAME && right.Kind == PLAYER {
		result := this.findAddLocation(left.node)
		if result != nil {
			return result
		}
		if node.GameState.Result == nil {
			fmt.Printf("Found spot to insert new user at Node: %v\n", node.Id)
			return node
		}
	}

	return nil
}

func (this *Bracket) AddParticipantIfAble(participantNumber ParticipantNumber) bool {

	fmt.Printf("%s\n", this.Root.PrintTree())
	// Winners Side
	node := this.Root.Left.node
	match := this.findAddLocation(node)

	if match == nil {
		return false
	}
	fmt.Printf("found location: %d\n", match.Id)

	// Setup new game node and move player down to it

	var newGame *Node
	fmt.Printf("Before: %s", this.Root.Left.node.PrintTree())
	fmt.Printf("match: %v", match)
	existingDrop := this.GetNode(match.Drop)
	fmt.Printf("Done \n")

	newGame = this.newNode(GAME, match, 0, 0, 1, match.Level+1)
	fmt.Printf("newGame = %v\n", newGame)

	newPlayerNode := this.newNode(PLAYER, newGame, 0, participantNumber, 1, match.Level+2)
	fmt.Printf("newPlayerNode = %v\n", newPlayerNode)

	var existingPlayer *Node
	if match.Right.Kind == PLAYER {
		existingPlayer = match.Right.node
		match.SetRight(newGame)
	} else if match.Left.Kind == PLAYER {
		existingPlayer = match.Left.node
		match.SetLeft(newGame)
	}

	newGame.SetLeft(newPlayerNode)
	newGame.SetRight(existingPlayer)

	fmt.Printf("Setup new winners-side game\n")
	fmt.Printf("After %s", this.Root.Left.node.PrintTree())

	///	losersMatch := this.findLosersSpot(participantNumber, this.Root.Right.node, newGame)

	var newLosersGame *Node
	//	if losersMatch == nil {
	//		return false
	//	}

	fmt.Printf("Setup new losers-side game\n")
	fmt.Printf("Before: %s", this.Root.Right.node.PrintTree())

	losersMatch := existingDrop.Parent.node
	newLosersGame = this.newNode(GAME, losersMatch, 0, 0, 2, losersMatch.Level+1)

	fmt.Printf("losersMatch = %v\n", losersMatch)
	fmt.Printf("newLosersGame = %v\n", newLosersGame)

	if losersMatch.Right.Id == existingDrop.Id {
		losersMatch.SetRight(newLosersGame)
	} else {
		losersMatch.SetLeft(newLosersGame)
	}

	newDropNode := this.newNode(DROP, newLosersGame, newGame.Id, 0, 2, losersMatch.Level+2)
	fmt.Printf("newDropNode = %v\n", newLosersGame)
	this.SetNodeDrop(newGame, newDropNode.Id)

	fmt.Printf("existingDrop = %v\n", existingDrop)
	newLosersGame.SetRight(existingDrop)
	newLosersGame.SetLeft(newDropNode)

	///	newDropNode.Drop = newGame.Id
	//newGame.Drop = newDropNode.Id
	existingDrop.SetParent(newLosersGame)
	fmt.Printf("After %s", this.Root.Right.node.PrintTree())

	this.Root.calculateSpans()

	return true
}

func (this *Tournament) DeleteParticipantIfAble(
	session dbr.SessionRunner,
	participantNumber ParticipantNumber) bool {

	node := this.Bracket.findParticipantNode(this.Bracket.Root, participantNumber)

	parentNode := node.Parent.node
	side := 0
	if parentNode.Right.Id == node.Id {
		side = 1
	}
	this.AddResult(session, node, side)
	return true
}

/*
		if node != nil {
			parentGame := node.Parent.node
			if (parentGame)
		}

}

*/

func (this *Bracket) findParticipantNode(node *Node, num ParticipantNumber) *Node {
	if node.Type == PLAYER && node.Participant == num {
		return node
	}
	if node.Left.node != nil {
		node := this.findParticipantNode(node.Left.node, num)
		if node != nil {
			return node
		}
	}
	if node.Right.node != nil {
		node := this.findParticipantNode(node.Right.node, num)
		if node != nil {
			return node
		}
	}
	return nil
}
