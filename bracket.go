package brackets

import (
	"bytes"
	"encoding/json"
	"errors"
	_ "errors"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"log"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const BUY = -1

type TournamentRecord struct {
	Id        int64     `json:"id" db:"id"`
	UserId    int64     `json:"userId" db:"user_id"`
	Name      string    `json:"name" db:"tournament_name"`
	CreatedAt time.Time `json:"createdAt" db:"creation_date"`
}

type Tournament struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"tournament_name"`

	Root  *Node       `json:"root"`
	Drops map[int]int `json:"drops"` // map of Node Id to Node Id

	//WinnersBracket *Bracket `json:"winnersBracket"`
	//LosersBracket  *Bracket `json:"losersBracket"`
	//WholeBracket   *Bracket `json:"wholeBracket"`
	//FinalGameId    int      `json:"finalGameId"`

	//Games          map[int]*Game `json:"games"`
	Nodes      map[int]*Node `json:"nodes"`
	RootNodeId int           `json:"rootNodeId"`
	Size       int           `json:"size"`

	Players map[int]*Player `json:"players"`
	PlayersIndex []int

	//Players map[int]*Player `json:"players"`

	Degree int `json:"degree"`
	Buys   int `json:"buys"`
	Depth  int `json:"depth"`


	EntryFee int `json:"entryFee"`
	TotalPot int `json:"totalPot"`

}


/*
    Tournament events:

       Events are produced based on tournament activity:

      tournament:  created, drawn, completed, player.added, player.removed, player.paid,
      game:   setPlayer, decided, nowReady,



      tournament.drawn
      tournament.drawn
*/

/*
type Bracket struct {
	tournament     *Tournament
	RootGame       *Game
	leafSlots      map[BigKey]*Slot
	leafSlotsOrder []BigKey
	gameIdQueue    *IntQueue

	leafSlotCounter int
	gameIdCounter   int

	//Games            map[string]string
}
*/

type Player struct {
	Id     int    `json:"id"`
	Number int    `json:"number"`
	Name   string `json:"name"`
	Paid   bool   `json:"paid"`
}

type Game struct {
	Id         int            `json:"id"`
	Left       *GameReference `json:"left"`
	Right      *GameReference `json:"right"`
	WinnerGame *GameReference `json:"winnerGame"`
	LoserGame  *GameReference `json:"loserGame"`
	Span       *Span          `json:"span"`
	Result     *GameResult    `json:"result"`

	Player1      int  `json:"player1"`
	Player2      int  `json:"player2"`
	IsLosersSide bool `json:"isLosersSide"`

	tournament *Tournament
}

type GameReference struct {
	game   *Game
	GameId int `json:"id"`
	Var    int `json:"var"`
}

type GameResult struct {
	GameId      int   `json:"gameId"`
	WinningSlot int   `json:"winningSlot"` // 0 means no win
	Time        int64 `json:"time"`        // order
}

type Slot struct {
	GameId int
	Slot   int
}

type SlotObj struct {
	GameObject *Game
	Slot       int
}

type PlayerList struct {
	Players []*Player `json:"players"`
}

/*
func (this *Tournament) BuildGames(node *Node, level int, side int) *Game {


	var *Node winners
	if level == 0 {
		winners = node.Left
	}

	if (node.Left != nil && node.Right != nil) {
		game := new(Game)
		game.Id = node.Id
		left := this.BuildGames(node.Left)
		ref := GameReference{left, left.Id, 0}

		game.Left = &ref

	}

}

*/

func (this *Tournament) GetPlayers() *PlayerList {
	list := make([]*Player, 0)
	for _, p := range this.Players {
		list = append(list, p)
	}
	return &PlayerList{list}
}

func (this *Tournament) GetResult(gameId int) *MatchResult {
	return this.Nodes[gameId].GameState.Result
}

//func (this *Tournament) AddResult(result MatchResult) error {
func (this *Tournament) AddResult(node *Node, winnerSide int) error {
	//	fmt.Printf("game(%d) before adding result: %v\n", result.NodeId, this. Nodes[result.NodeId])

	if node.Left.node.Player == 0 || node.Right.node.Player == 0 {
		return errors.New("Game can't be completed, state error")
	}

	var winnerNodeId int
	var loser int
	var winner int
	if winnerSide == 1 {
		winner = node.Left.node.Player
		winnerNodeId = node.Left.Id
		loser = node.Right.node.Player
	} else {
		winner = node.Right.node.Player
		winnerNodeId = node.Right.Id
		loser = node.Left.node.Player
	}

	fmt.Printf("setting node (%d) player to %d\n", node.Id, winner)
	node.Player = winner

	var dropNode *Node = nil
	dropNodeId, has := this.Drops[node.Id]
	if has {
		dropNode = this.GetNode(dropNodeId)
	}
	if dropNode != nil {
		dropNode.Player = loser
	}

	//winningPlayer := winnerNode.GameState.Player
	//losingPlayer := winnerNode.GameState.Player

	res := new(MatchResult)
	res.WinningNode = winnerNodeId
	res.WinningPlayer = winner
	res.DropNode = dropNodeId
	res.LosingPlayer = loser
	res.WinningSlot = winnerSide

	if dropNode != nil {
		fmt.Printf("dropNode id = %d, setting player to %d\n", dropNode.Id, loser)
		dropNode.Player = loser
	}

	/*
		if node.Parent.Id > 0 {
			parent, has := this.Nodes[node.Parent.Id]
			if has {
				parent.Player = winner
			}
		}
	*/

	res.Time = time.Now().Unix()
	node.GameState.Result = res

	//	node.AddResult(*res)

	//	if (game.Player1 == 0 || game.Player2 == 0) {
	//	}
	return nil
}

func (this *Tournament) RemoveResult(gameId int) error {
	node := this.GetNode(gameId)
	if node.GameState.Result == nil {
		return errors.New("unable to complete game")
	}
	fmt.Printf("node.Parent = %v\n", node.Parent)
	if node.Parent.Id != 0 && node.Parent.node.GameState.Result != nil {
		return errors.New("unable to complete game")
	}
	//	result := node.GameState.Result
	dropNode := this.GetNode(node.GameState.Result.DropNode)
	if dropNode != nil {
		dropNode.Player = 0
	}

	node.Player = 0

	node.GameState.Result = nil

	//	node.RemoveResult()

	return nil
}

func (this *Tournament) AddPlayer(name string) {

	num := len(this.Players) + 1
	player := new(Player)
	player.Id = num
	player.Number = num
	player.Name = name
	player.Paid = false
	this.Players[num] = player
	fmt.Printf("AddPlayer: size now = %d\n", len(this.Players))
}

func (this *Tournament) AddPlayers(players []string) {
	for _, p := range players {
		this.AddPlayer(p)
	}
	fmt.Printf("AddPlayers: size = %d\n", len(this.Players))
}

func (this *Tournament) RemovePlayer(name string) {
	fmt.Printf("Remove Player %s\n", name)
	ids := make([]int, 0)
	for k, p := range this.Players {
		fmt.Printf("player = %s\n", p.Name)
		if p.Name == name {
			ids = append(ids, k)
			fmt.Printf("removing player with Id %d\n", p.Id)
		}
	}
	for _, mapKey := range ids {
		delete(this.Players, mapKey)
	}
}

func (this *Tournament) RemovePlayers(players []string) {
	fmt.Printf("RemovePlayers: %v\n", players)
	for _, p := range players {
		this.RemovePlayer(p)
	}
}

func (this *Tournament) Store(db *dbr.Session, userId int64) error {

	fmt.Printf("Tournament.Store: id = %d\n", this.Id)
	var err error
	var data []byte

	if data, err = json.Marshal(this); err != nil {
		log.Fatalf("Store,Marshall Tournament failed: %v", err)
		return err
	}

	var id int64
	if this.Id == 0 {
		err = db.InsertInto("tournaments").
			Columns("tournament_name", "tournament_data", "user_id", "creation_date").
			Values("test", data, 1, time.Now()).
			Returning("id").
			Load(&id)

		if err != nil {
			log.Fatalf("Insert Tournaments failed: %v", err)
			return err
		}
		this.Id = id
	} else {
		_, err = db.Update("tournaments").Set("tournament_data", data).
			Where("id = ?", this.Id).Exec()
		if err != nil {
			log.Fatalf("Update Tournament failed: %v", err)
			return err
		}
	}

	return nil
}

func ListTournaments(db *dbr.Session, userId int64) ([]*TournamentRecord, error) {

	fmt.Printf("selecting tournamentRecords for user: %d\n", userId)
	result, err := db.Select("id, user_id, tournament_name, creation_date").From("tournaments").
		Where("user_id = ?", userId).Rows()

	if err != nil {
		log.Fatalf("Select Tournament failed: %v\n", err)
		return nil, err
	}

	var records = make([]*TournamentRecord, 0)
	//if err := dec.Decode(&val); err != nil {

	for {
		if result.Next() == false {
			if err := result.Close(); err != nil {
				return records, err
			} else {
				return records, nil
			}
		}
		var id int64
		var uid int64
		var name string
		var creationDate time.Time
		if err := result.Scan(&id, &uid, &name, &creationDate); err != nil {
			_ = result.Close()
			log.Fatalf("Scan tournament data failed: %v\n", err)
			return nil, err
		}
		rec := TournamentRecord{id, uid, name, creationDate}
		records = append(records, &rec)
	}
}

func LoadTournament(db *dbr.Session, id int64) (*Tournament, error) {

	fmt.Printf("selecting tournament with id: %d\n", id)
	result, err := db.Select("tournament_data").From("tournaments").Where("id = ?", id).Rows()
	if err != nil {
		log.Fatalf("Select Tournament failed: %v\n", err)
		return nil, err
	}
	if result.Next() == false {
		return nil, nil
	}

	var data []byte
	if err := result.Scan(&data); err != nil {
		log.Fatalf("Scan tournament data failed: %v\n", err)
		return nil, nil
	}

	result.Close()

	var t *Tournament
	if err = json.Unmarshal(data, &t); err != nil {
		log.Fatalf("Unmarshal tournament data failed: %v", err)
		return nil, nil
	}

	fmt.Printf("t = %v\n", t)
	t.Id = id

	for _, n := range t.Nodes {
		n.internalize(t)
	}
	t.Root = t.Nodes[t.RootNodeId]
	/*
		if (t.FinalGameId != 0) {
			t.GetGame(t.FinalGameId).internalize(t)
		}
	*/

	return t, nil

}

func (this *Player) String() string {
	return fmt.Sprintf("%d:%s", this.Number, this.Name)
}

func (this *Game) SetLeftDropSource(left *Game) {
	this.Left = &GameReference{left, left.Id, 1}
}

func (this *Game) SetLeft(left *Game) {
	this.Left = &GameReference{left, left.Id, 0}
}

func (this *Game) SetRight(right *Game) {
	this.Right = &GameReference{right, right.Id, 0}
}

func (this *Game) encode() string {
	return fmt.Sprintf("%d:Left(%v):Right(%v):Winner(%v):Loser(%v)\n",
		this.Id, this.Left, this.Right, this.WinnerGame, this.LoserGame)
}

func (this *Game) PrintTree() string {
	return this.printTreeNode(new(bytes.Buffer), true, new(bytes.Buffer)).String()
}

func (this *Game) printTreeNode(prefix *bytes.Buffer, isTail bool, buf *bytes.Buffer) *bytes.Buffer {

	fmt.Printf("printTreeNode %v\n", this)

	if this.Right != nil {

		if this.Right.Winner() == true {
			newPrefix := new(bytes.Buffer)
			tail := "    "
			if isTail {
				tail = "|   "
			}
			newPrefix.WriteString(prefix.String())
			newPrefix.WriteString(tail)
			this.Right.game.printTreeNode(newPrefix, false, buf)
		} else {
			buf.WriteString(prefix.String())
			buf.WriteString("    --- ")
			buf.WriteString(fmt.Sprintf("(%d)\n", this.Right.game.Id))
		}
	}
	t := "└── "
	if !isTail {
		t = "┌── "
	}
	buf.WriteString(prefix.String())
	buf.WriteString(t)
	buf.WriteString(fmt.Sprintf("%d\n", this.Id))

	if this.Left != nil {

		if this.Left.Winner() == true {
			newPrefix := new(bytes.Buffer)
			tail := "|   "
			if isTail {
				tail = "    "
			}
			newPrefix.WriteString(prefix.String())
			newPrefix.WriteString(tail)
			this.Left.game.printTreeNode(newPrefix, true, buf)
		} else {
			buf.WriteString(prefix.String())
			buf.WriteString("    --- ")
			buf.WriteString(fmt.Sprintf("(%d)\n", this.Left.game.Id))
		}
	}
	return buf
}

func (this *Tournament) GetNode(id int) *Node {
	node := this.Nodes[id]
	return node
}

func (this *Tournament) GetGame(id int) *Game {
	/*
		game := this.Games[id]
		game.tournament = this
		return game
	*/
	return nil
}

//   Encoded string:    GameId.Var

/*
   Meanings of Var:

       For tree.go structure association (Left and Right) Var means:
           0 - means its a direct parent child relationship
           1 - means a association from a Loser tier game to a winners tier game (a drop down)

       For associations from a game to the games the winner and loser propagate to (WinnersGame and LosersGame)
       Var means:

           The slot # in the target game the player goes to.
*/

func (this *GameReference) Slot() int {
	return this.Var
}

func (this *GameReference) Winner() bool {
	if this.Var == 0 {
		return true
	}
	return false
}

func (this *GameReference) GetGame(t *Tournament) *Game {
	/*
		if this.GameId != 0 && this.game == nil {
			this.game = t.Games[this.GameId]
		}
		return this.game
	*/
	return nil
}

func (this *GameReference) UnmarshalJSON(data []byte) error {
	var f interface{}
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	values := f.(map[string]interface{})
	this.GameId = int(values["id"].(float64))
	this.Var = int(values["var"].(float64))

	return nil
}

func (this GameReference) MarshallJSON() ([]byte, error) {

	data := make(map[string]interface{}, 0)
	data["id"] = this.game.Id
	data["var"] = this.Var
	var b []byte
	var err error
	if b, err = json.Marshal(data); err != nil {
		panic(err)
	}
	return b, err

}

func (this GameReference) String() string {
	return fmt.Sprintf("%d.%d", this.game.Id, this.Var)
}

var GameCounter int = 1

func (this *Tournament) Draw(session *dbr.Session, userId int64) error {

	fmt.Printf("Draw:   # Players: %d\n", len(this.Players))
	vals := make([]int, len(this.Players))
	for i := 0; i < len(this.Players); i++ {
		vals[i] = i + 1
	}
	Shuffle(vals)

	newPlayers := make(map[int]*Player)
	count := 0
	for _, player := range this.Players {
		newPlayers[player.Id] = player
		player.Number = vals[count]
		count += 1
	}
	fmt.Printf("Players: %v\n", this.Players)
	fmt.Printf("newPlayers: %v\n", newPlayers)

	this.Players = newPlayers
	if err := this.Store(session, userId); err != nil {
		return errors.New(fmt.Sprintf("Store,Marshall Tournament failed: %v", err))
	}

    return nil
}

func (this *Tournament) BuildBrackets(session *dbr.Session, userId int64) (*Tournament, error) {

	IdCounter = 1
	GameCounter = 1
	if len(this.Players) < 2 {
		return nil, errors.New("Not enough players")
	}
	this.Nodes = make(map[int]*Node)
	this.Drops = make(map[int]int)

	degree := 2
	size := 2
	for {
		if size >= len(this.Players) {
			this.Size = size
			this.Degree = degree - 1
			break
		}
		size = size << 1
		degree += 1
	}
	//	this.Draw()

	vals := make([]int, len(this.Players))
	for i := 0; i < len(this.Players); i++ {
		vals[i] = i + 1
	}
//	Shuffle(vals)

	this.Buys = this.Size - len(this.Players)

	root := this.CreateDoubleElimBracket()
	root.calculateSpans()
	root.Right.node.Span.Upper += 20

	if err := this.Store(session, userId); err != nil {
		return nil, errors.New(fmt.Sprintf("Store,Marshall Tournament failed: %v", err))
	}

	return this, nil
}

func NewTournament2() *Tournament {

	t := new(Tournament)
	//	t.Games = make(map[int]*Game)
	t.Nodes = make(map[int]*Node)
	t.Drops = make(map[int]int)
	t.Players = make(map[int]*Player)
	return t

}

func (this Tournament) String() string {

	var b strings.Builder

	fmt.Fprintf(&b, "Tournament id: %d\n", this.Id)
	//	fmt.Fprintf(&b, "Games = %v\n", this.Games)

	return b.String()
}

func (this *Tournament) CreateDoubleElimBracket() *Node {

	fmt.Printf("# Players = %d\n", len(this.Players))

	spots := make([]int, len(this.Players))

	for i := 1; i <= len(this.Players); i++ {
		spots[i-1] = i
	}

	this.createPlayersIndex()
	winnerRoot := this.CreateBracket(spots, 1)

	fmt.Printf("winners:\n\n")
	fmt.Printf("%v", winnerRoot.PrintTree())
	losers := winnerRoot.GetInnerNodes()

	fmt.Printf("inner nodes: %v\n", losers)

	loserRoot := this.CreateBracket(losers, 2)
	fmt.Printf("losers:\n\n")
	fmt.Printf("%v", loserRoot.PrintTree())

	root := new(Node)
	root.Id = IdCounter
	IdCounter += 1
	root.Tier = 1
	root.Level = 0
	this.Nodes[root.Id] = root
	root.SetLeft(winnerRoot)
	root.SetRight(loserRoot)

	fmt.Printf("all\n\n")
	fmt.Printf("%v\n", root.PrintTree())

	this.RootNodeId = root.Id

	return root
}

func (this *Tournament) NodesString() string {
	var b strings.Builder
	a := make([]int, len(this.Nodes))
	for i, n := range this.Nodes {
		a[i-1] = n.Id
	}

	sort.Ints(a)

	fmt.Fprintf(&b, "%v ", a)
	fmt.Fprintf(&b, "\n")
	return b.String()
}

func (this *Tournament) CompleteBracket(q *NodeQueue, level int, tier int) *Node {

	var last *Node
	l := level

	fmt.Printf("WTierTotal = %d\n", ComputeWTierTotal(len(this.Players)))
	fmt.Printf("WTier1 = %d\n", ComputeWTier1(this.Size, len(this.Players)))
	fmt.Printf("WTier2 = %d\n", ComputeWTier2(this.Size, len(this.Players)))
	for {

		if q.Size() <= 1 {
			if q.Size() == 1 {
				last = q.Remove()
			}
			break
		}

		p1 := q.Remove()
		p2 := q.Remove()
		fmt.Printf("Pulled node %d from Queue \n", p1.Id)
		fmt.Printf("Pulled node %d from Queue \n", p2.Id)

		p3 := this.newNode(GAME, nil, 0, 0, tier, l)
		p3.SetLeft(p1)
		p3.SetRight(p2)
		q.Add(p3)
		fmt.Printf("added game %d to Queue\n", p3.Id)
		last = p3
	}
	level -= 1
	return last
}

func (this *Tournament) createPlayersIndex() {

	this.PlayersIndex = make([]int, 0)
	for k, _ := range this.Players {
		this.PlayersIndex = append(this.PlayersIndex, k)
	}
}

func (this *Tournament) getPlayerAtIndex(i int) *Player {
	return this.Players[this.PlayersIndex[i-1]]
}

func (this *Tournament) GenerateWinnersBracket(q *NodeQueue, nodeIds []int) *Node {

	level := this.Degree + 1
	for _, id := range nodeIds {
		var n *Node
		n = this.newNode(PLAYER, nil, 0,
			this.getPlayerAtIndex(id).Id, 1, level)
		q.Add(n)
	}

	level -= 1

	return this.CompleteBracket(q, level, 1)
}

/*
   Given a set of sets of ids, where each set are the ids introduced at a given level, generate the tree.

    The tree is generated by creating nodes for each new id, adding them to the queue, then poping nodes 2 at a time
    from the Queue, create a new node from them  and put it back in the queue, the last node is the tree root.

    When there are no more ids to introduce finish processing the Queue to generate games

*/
func (this *Tournament) CreateBracket(leafs []int, tier int) *Node {

	fmt.Printf("CreateBracket: leafs sizeof: %d\n", len(leafs))
	fmt.Printf("CreateBracket:  Degree = %d\n", this.Degree)
	q := new(NodeQueue)

	if tier == 1 {
		winners := this.GenerateWinnersBracket(q, leafs)
		fmt.Printf("%v\n", winners.PrintTree())
		return winners
	}

	fmt.Printf("len of leafs = %d\n", len(leafs))

	losers := this.GenerateLosersBracket2(leafs)

	return losers

}

func (this *Tournament) newNode(t NodeType, parent *Node, dropGameId int, playerId int, tier int, level int) *Node {
	node := NewNode(t, parent, dropGameId, playerId, tier, level)
	this.Nodes[node.Id] = node
	return node
}

func (this *Tournament) GenerateLosersBracket2(nodeIds []int) *Node {
	multiplier := 1
	level := 1
	root := this.newNode(GAME, nil, 0, 0, 2, level)
	level += 1

	parentTier := make(map[int]*Node, 0)
	parentTier[root.Id] = root
	usedDrops := 0
	totalDrops := len(nodeIds)
	totalGames := 0
	usedGames := 0

//	loserSideGames := ComputeLTierTotal(len(this.Players))

	for loop := 1; ; loop++ {

		if (usedDrops == totalDrops) {
			break
		}
		remainingDrops := totalDrops - usedDrops

		/*
		dontExpand := false
		if (remainingDrops <= 4*multiplier) {
			dontExpand = true
		}
		 */

		firstMultiplierTier := false
		if loop%2 == 1 {
			firstMultiplierTier = true
		}
		if loop != 1 && loop%2 == 1 {
			multiplier = multiplier * 2
		}
		M := multiplier

//		if len(parentTier) != M {
//			log.Fatal("wrong number of parent tier nodes")
//		}

		mTier := 1
		if !firstMultiplierTier {
			mTier = 2
		}
		fmt.Printf("Tier: %d.%d\n", M, mTier)
		fmt.Printf("Drops Remaining: %d\n", totalDrops-usedDrops)

		tmpTier := make(map[int]*Node, 0)


		//children := M / 2 + 1
//		tierGames := M
		tierDrops := M
		if (remainingDrops == len(parentTier) *2) {
			tierDrops = remainingDrops
		} else if (remainingDrops < M) {
			log.Fatal("Shouldn't happen")
		} else if remainingDrops > M && remainingDrops <= 2*M {
//			if (remainingDrops % 2 == 1) {
				tierDrops = remainingDrops - len(parentTier)
//			}
///			tierDrops = 2 * int(math.Pow(float64(M), float64(2)))
//			tierDrops = tierDrops - remainingDrops
			firstMultiplierTier = true
			//tierDrops = remainingDrops - M
//			tierGames = M*2 - tierDrops
		} else if remainingDrops > M && remainingDrops < 4*M {
			tierDrops = 4 * len(parentTier) - remainingDrops;
			firstMultiplierTier = true
		} else {
			tierDrops = len(parentTier)
		}



		tierGames := len(parentTier) * 2 - tierDrops
		fmt.Printf("tierGames: %d, tierDrops: %d\n", tierGames, tierDrops)
		childGamesUsed := 0
		tierDropsUsed := 0
		if firstMultiplierTier || (remainingDrops == len(parentTier) *2) {
			leftlean := true
			for _, n := range parentTier {
				fmt.Printf("childGames: %d, tierGames: %d\n", childGamesUsed, tierGames)
				if leftlean {
					if (tierDropsUsed < tierDrops) {
						id := nodeIds[len(nodeIds)-1-usedDrops]
						left := this.newNode(DROP, n, id, 0, 2, level)
						this.Drops[id] = left.Id
						usedDrops += 1
						n.SetLeft(left)
						fmt.Printf("      new L Drop: %d\n", left.Id)
						tierDropsUsed += 1
						if usedDrops == totalDrops {
							break
						}
					} else {
						left := this.newNode(GAME, n, 0, 0, 2, level)
						tmpTier[left.Id] = left
						n.SetLeft(left)
						totalGames += 1
						fmt.Printf("      new L game: %d\n", left.Id)
						childGamesUsed += 1
					}

					if childGamesUsed < tierGames {
						right := this.newNode(GAME, n, 0, 0, 2, level)
						tmpTier[right.Id] = right
						n.SetRight(right)
						totalGames += 1
						fmt.Printf("      new R game: %d\n", right.Id)
						childGamesUsed += 1
					} else {
						id := nodeIds[len(nodeIds)-1-usedDrops]
						right := this.newNode(DROP, n, id, 0, 2, level)
						n.SetRight(right)
						this.Drops[id] = right.Id
						usedDrops += 1
						tierDropsUsed += 1
						fmt.Printf("      new R Drop: %d\n", right.Id)
					}

					if usedDrops == totalDrops {
						break
					}
				} else {
					if usedDrops == totalDrops {
						break
					}

					fmt.Printf("childGames: %d, tierGames: %d\n", childGamesUsed, tierGames)
					if childGamesUsed < tierGames {
						left := this.newNode(GAME, n, 0, 0, 2, level)
						tmpTier[left.Id] = left
						n.SetLeft(left)
						totalGames += 1
						fmt.Printf("      new L game: %d\n", left.Id)
						childGamesUsed += 1
					} else {

						id := nodeIds[len(nodeIds)-1-usedDrops]
						left := this.newNode(DROP, n, id, 0, 2, level)
						this.Drops[id] = left.Id
						usedDrops += 1
						n.SetLeft(left)
						fmt.Printf("      new L Drop: %d\n", left.Id)
						tierDropsUsed += 1
					}
					if usedDrops == totalDrops {
						break
					}

					if (tierDropsUsed < tierDrops) {
						id := nodeIds[len(nodeIds)-1-usedDrops]
						right := this.newNode(DROP, n, id, 0, 2, level)
						n.SetRight(right)
						this.Drops[id] = right.Id
						usedDrops += 1
						tierDropsUsed += 1
						fmt.Printf("      new R Drop: %d\n", right.Id)
						if usedDrops == totalDrops {
							break
						}
					} else {

						right := this.newNode(GAME, n, 0, 0, 2, level)
						tmpTier[right.Id] = right
						n.SetRight(right)
						totalGames += 1
						fmt.Printf("      new R game: %d\n", right.Id)
						childGamesUsed += 1
					}

				}
				leftlean = !leftlean

			}
		} else {

			for _, n := range parentTier {
				left := this.newNode(GAME, n, 0, 0, 2, level)
				tmpTier[left.Id] = left
				n.SetLeft(left)
				totalGames += 1
				fmt.Printf("      new L Expansion game: %d\n", left.Id)
				if usedDrops == totalDrops {
					break
				}

				right := this.newNode(GAME, n, 0, 0, 2, level)
				tmpTier[right.Id] = right
				n.SetRight(right)
				totalGames += 1
				fmt.Printf("      new R Expansion game: %d\n", right.Id)
				if usedDrops == totalDrops {
					break
				}
			}

		}

		parentTier = tmpTier
		level += 1
		fmt.Printf("End of loop: nextLoopNodes: %d, usedGames: %d, usedDrops: %d\n", len(parentTier), usedGames, usedDrops)

	}

	this.Depth = level

	return root
}

func (this *Tournament) GenerateLosersBracket(nodeIds []int) *Node {

	multiplier := 1
	level := 1
	root := this.newNode(GAME, nil, 0, 0, 2, level)
	level += 1

	lastTier := make(map[int]*Node, 0)
	lastTier[root.Id] = root

	games := 0
	dropLevel := true
	usedDrops := 0

	for tier := 1; ; tier++ {
		if tier%2 == 1 && tier != 1 {
			multiplier = multiplier * 2
		}

		fmt.Printf("LEVEL: %d, TIER_SIZE: %d\n", level, len(lastTier))

		dropsRemaining := len(nodeIds) - usedDrops
		fmt.Printf("Drops remaining: %v\n", dropsRemaining)

		if len(lastTier) == 0 {
			break
		}

		tmpTier := make(map[int]*Node, 0)

		fmt.Printf(" Creating %d child nodes for Tier\n", 2*len(lastTier))
		for _, n := range lastTier {
			left := this.newNode(DROP, n, 0, 0, 2, level)
			n.SetLeft(left)

			right := this.newNode(DROP, n, 0, 0, 2, level)
			n.SetRight(right)

		}

		// N := multiplier
		// tier := the tier number
		// D := number of remaining drops
		// if (tier % 2 == 1) this is the first multiplier tier
		// if (tier % 2 == 0) this is the second multiplier tier
		//
		// For both tiers of a multiplier there are N nodes and 2*N slots
		//
		//   S = total # of slots in a tier
		//   C = the # of slots that point to child game nodes
		//   T = the # of slots that point to drop nodes (drops or players)
		//
		//

		// D = number of the child references that are drops

		couldNotUseAllDrops := false
		if dropsRemaining < len(lastTier)*2 {
			couldNotUseAllDrops = true
			fmt.Printf(" Could use all drops this level\n")
		}
		_ = couldNotUseAllDrops

		gamesRequiredThisLevel := len(lastTier)
		if dropsRemaining < (len(lastTier) * 4) {
			dropLevel = true
			gamesRequiredThisLevel = dropsRemaining % (len(lastTier))
		}

		//if  nearEnd {!dropLevel {
		if !dropLevel {
			fmt.Printf("  NEAR END LEVEL: %d \n", level)
			fmt.Printf("  Branching %d nodes\n", 2*len(lastTier))
			for _, n := range lastTier {
				left := this.Nodes[n.Left.Id]
				left.Type = GAME
				tmpTier[left.Id] = left
				n.SetLeft(left)
				games += 1
				fmt.Printf("      new L Expansion game: %d\n", left.Id)

				right := this.Nodes[n.Right.Id]
				right.Type = GAME
				tmpTier[right.Id] = right
				n.SetRight(right)
				games += 1
				fmt.Printf("      new R Expansion game: %d\n", right.Id)
			}

		} else {
			usedGames := 0
			ct := 0
			for _, n := range lastTier {

				fmt.Printf("    Processing node %d on level %d\n", n.Id, level)
				fmt.Printf("    USED DROPS: %d, # nodes: %d\n", usedDrops, len(nodeIds))

				if usedDrops < len(nodeIds) {

					if (level+ct)%2 == 0 {
						fmt.Printf("      Setting L oritented drop Node\n")

						left := this.Nodes[n.Left.Id]
						id := nodeIds[len(nodeIds)-1-usedDrops]
						this.Drops[id] = left.Id
						usedDrops += 1
						left.Drop = id
						n.SetLeft(left)

						if usedGames < gamesRequiredThisLevel {
							right := this.Nodes[n.Right.Id]
							right.Type = GAME
							tmpTier[right.Id] = right
							n.SetRight(right)
							usedGames += 1
							games += 1
						} else {
							right := this.Nodes[n.Right.Id]
							id := nodeIds[len(nodeIds)-1-usedDrops]
							this.Drops[id] = right.Id
							usedDrops += 1
							right.Drop = id
							n.SetRight(right)
						}
					} else {

						fmt.Printf("      Setting R oritented drop Node\n")

						right := this.Nodes[n.Right.Id]
						id := nodeIds[len(nodeIds)-1-usedDrops]
						this.Drops[id] = right.Id
						usedDrops += 1
						right.Drop = id
						n.SetRight(right)

						if usedGames < gamesRequiredThisLevel {
							left := this.Nodes[n.Left.Id]
							left.Type = GAME
							tmpTier[left.Id] = left
							n.SetLeft(left)
							usedGames += 1
							games += 1
						} else {
							left := this.Nodes[n.Left.Id]
							id := nodeIds[len(nodeIds)-1-usedDrops]
							this.Drops[id] = left.Id
							usedDrops += 1
							left.Drop = id
							n.SetLeft(left)
						}

					}
				} else {
					fmt.Printf("LESS DROPS than NODES: %d\n", level)
				}
				ct += 1
			}

		}

		lastTier = tmpTier
		//		tmpTier := make([]*Node, 0)
		dropLevel = !dropLevel
		level += 1
		fmt.Printf("End of loop: nextLoopNodes: %d, usedGames: %d, usedDrops: %d\n", len(lastTier), games, usedDrops)
	}

	return root
}

/*
func NewBracket(t *Tournament) *Bracket {
	b := new(Bracket)
	b.tournament = t
	b.leafSlots = make(map[BigKey]*Slot)
	b.leafSlotsOrder = make([]BigKey, 0)
	b.gameIdQueue = new(IntQueue)
	b.leafSlotCounter = 0
	b.gameIdCounter = 0
	return b
}


func (this *Bracket) CompressLeafSlotOrder() []BigKey {

	t := make([]BigKey, 0)

	for _, v := range this.leafSlotsOrder {
		if int64(v) != 0 {
			t = append(t, v)
		}
	}

	return t
}
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

func NewSpan(upper int, lower int) *Span {
	span := new(Span)
	span.Upper = upper
	span.Lower = lower
	return span
}

func Shuffle(vals []int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}

// Game number calculations

func ComputeDegree(size int) int {

	s := size
	d := 1
	for {
		if s == 0 {
			fmt.Printf("degree for size %d  = %d\n", size, d)
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
	vm := int(math.Min(float64((size)-(size-players)), float64(size>>2)))
	fmt.Printf("vm = %d\n", vm)

	return size >> 2

	//vm := int(math.Min(float64( (size >> 1) - (size-players)), float64(size>>2) ))
}

func ComputeLTier1(size int, players int) int {
	v := int(math.Max(float64(0), float64((size>>2)-(size-players))))
	return v
}

func ComputeLTier2(size int, players int) int {
	//    v := int(math.Max(float64(0), float64((size>>2) - ((size - players)>>2))))

	//	v := (size >> 1) - ((size-players+1) )
	//	fmt.Printf("v = %d\n", v)
	vm := int(math.Min(float64((size>>1)-(size-players)), float64(size>>2)))
	//	fmt.Printf("vm = %d\n", vm)
	return vm
	//	v := int(math.Min(float64(size >> 2), float64(size>>2 - (size-players)-1)))
	//	return v
	//	max (0, (size >> 2) - ((size-players) >> 2) )
}

//type BigKey uint64

/*
func MakeBigKey(p1 int, p2 int) BigKey {
	key := uint64(p1) | (uint64(p2) << 32)
	return BigKey(key)
}

func (this *BigKey) getPart1() int {
	return int(uint64(*this) & 0x00000000FFFFFFFF)
}

func (this *BigKey) getPart2() int {
	return int(uint64(*this) >> 32)
}

func (this BigKey) print() {
	fmt.Printf("p1: %d, p2: %d\n", this.getPart1(), this.getPart2())

}
func (this BigKey) String() string {
	return fmt.Sprintf("GameId: %d, Slot: %d\n", this.getPart1(), this.getPart2())
}

*/

func (this Slot) String() string {
	return fmt.Sprintf("GameId: %d, Slot: %d\n", this.GameId, this.Slot)
}
