package brackets

import (
	"encoding/json"
	"errors"
	_ "errors"
	"fmt"

	//	"github.com/gocraft/dbr/v2/dialect"
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/gocraft/dbr/v2"
)

type BracketType int

// ParticipantNumber identifies each participant in a tournament
type ParticipantNumber int

const (
	SINGLE_ELIMINATION = 1
	DOUBLE_ELIMINATION = 2
	ROUND_ROBIN        = 3
)

// Tournament State
const (
	NEEDS_DRAW     = "NeedsDraw"
	READY_TO_START = "Ready"
	UNDERWAY       = "Underway"
	REGISTRATION   = "Registration"
	COMPLETE       = "Complete"
)

type TournamentRecord struct {
	Id               int64     `json:"id"`
	UserId           int64     `json:"userId"`
	Subject          string    `json:"subject"`
	Name             string    `json:"name"`
	CreatedAt        time.Time `json:"createdAt"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	State            string    `json:"tournamentState"`
	ParticipantCount int       `json:"participantCount"`
}

/**
  Models an elimination tournament.

   A tournament begins with a number of competitors (teams or players).

   The tournament consists of series of matches between pairs of competitors.

   The initial matches can be determined randomly or by a seeding step.
   Depending on the number of competitors some may be granted a "bye"
    (a "bye" is a free win, or a match with noone, which exists to balance
     the tournament, which operates on the priciple that the number
      of competitors is a power of 2).

   The "brackets" are represented as a tree structure, where each node
    in the tree represents a match between 2 competitors.    The root
    node of the tree represents the championship match.
    A double elimination tournament is reprsented by 2 separate tree's
       with a common parent node.

    A tournmanent has a setup phase and a underway phase.  During
      setup competitors can be added and removed at will and each
      change results in a  change to the brackets (tree structure).


*/
type Tournament struct {
	Id        int64       `json:"id"              db:"id"`
	UserId    int64       `json:"userId"          db:"user_id"`
	AccountId int64       `json:"accountId"       db:"account_id"`
	Subject   string      `json:"subject"         db:"subject"`
	Name      string      `json:"name"            db:"tournament_name"`
	Typ       BracketType `json:"tournamentType"  db:"tournament_type"`

	Bracket *Bracket `json:"bracket"`

	Participants                 map[ParticipantNumber]*Participant `json:"participants"`
	ParticipantsByPlayerId       map[int]*Participant
	ParticipantCount             int                 `json:"participant_count"`
	RecyclableParticipantNumbers []ParticipantNumber `json:"recyclable_participant_numbers"`

	EntryFee int `json:"entryFee"`
	TotalPot int `json:"totalPot"`

	State     string    `json:"tournamentState"       db:"tournament_state"`
	CreatedAt time.Time `json:"createdAt"   db:"creation_date"`
	StartTime time.Time `json:"startTime"   db:"end_time"`
	EndTime   time.Time `json:"endTime"     db:"start_time"`

	FinalGame NodeId `json:"finalGame"     db:"final_game"`
	//	participantNumberCount int    `json:"participantNumberCount" db:"participant_number_count"`
	GamesPlayed int `json:"games-played"`
}

type TournamentStatus struct {
	StartTime time.Time         `json:"startTime"`
	EndTime   time.Time         `json:"endTime"`
	Winner    ParticipantNumber `json:"winner"`
	RunnerUp  ParticipantNumber `json:"runnerUp"`
}

/*
   Tournament events:

      Events are produced based on tournament activity:

     tournament:  created, drawn, completed, player.added, player.removed, player.paid,
     game:   setPlayer, decided, nowReady,



     tournament.drawn
     tournament.drawn
*/

type Participant struct {
	Number       ParticipantNumber `json:"participantNumber"`
	TournamentId int               `json:"TournamentId"`
	PlayerId     int               `json:"playerId"`

	// Defaults to the players name but can be assigned any name for the tournament
	Name       string `json:"name"`
	PaidAmount int    `json:"paidAmount"`
	Paid       bool   `json:"paid"`
}

type GameResult struct {
	GameId      int   `json:"gameId"`
	WinningSlot int   `json:"winningSlot"` // 0 means no win
	Time        int64 `json:"time"`        // order
}

func DeleteTournaments(session dbr.SessionRunner, ids []int64) {

	for _, id := range ids {
		fmt.Printf("Deleteting t with id %v\n", id)
		session.DeleteFrom("tournaments").Where(dbr.Eq("id", id)).Exec()
	}

}

func (this *Tournament) GetStatus(session dbr.SessionRunner) (*TournamentStatus, error) {
	status := new(TournamentStatus)
	status.StartTime = this.StartTime
	status.EndTime = this.EndTime

	if this.FinalGame != 0 {
		result := this.GetResult(this.FinalGame)
		status.Winner = result.WinningParticipant
		status.RunnerUp = result.LosingParticipant
	}

	return status, nil
}

func (this *Tournament) Start(session dbr.SessionRunner, subject string) error {

	/*
		for _, p := range this.Participants {

			if (p.PaidAmount < this.EntryFee) {
				return errors.New("One or more participants have not yet paid")
			}

		}
	*/

	this.State = UNDERWAY
	this.StartTime = time.Now()
	if err := this.Store(session, subject); err != nil {
		return errors.New(fmt.Sprintf("Store,Marshall Tournament failed: %v", err))
	}

	return nil
}

func (this *Tournament) SetParticipantPaid(session dbr.SessionRunner, subject string, participantNumber int, paid bool) error {

	participant := this.Participants[ParticipantNumber(participantNumber)]
	participant.PaidAmount = 10
	participant.Paid = true
	//this.ParticipantsByPlayerId[playerId]

	fmt.Printf("set particpantNumber %v to paid\n", participantNumber)

	if err := this.Store(session, subject); err != nil {
		return errors.New(fmt.Sprintf("Store,Marshall Tournament failed: %v", err))
	}

	return nil
}

func (this *Tournament) GetParticipants() map[ParticipantNumber]*Participant {
	return this.Participants
}

type ParticipantList struct {
	Participants []*Participant `json:"participants"`
}

func (this *Tournament) GetParticipantList() *ParticipantList {

	list := make([]*Participant, 0)
	for _, p := range this.Participants {
		list = append(list, p)
	}
	Number := func(p1, p2 *Participant) bool {
		return p1.Number < p2.Number
	}
	By(Number).Sort(list)

	return &ParticipantList{list}
}

func (this *Tournament) FindParticipantByName(name string) *Participant {
	for _, p := range this.Participants {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (this *Tournament) GetResult(gameId NodeId) *MatchResult {
	return this.Bracket.Nodes[gameId].GameState.Result
	//	return this.Nodes[gameId].GameState.Result
}

func (this *Tournament) Complete(session dbr.SessionRunner, winner int, node *Node) error {
	this.EndTime = time.Now()
	this.FinalGame = node.Id
	this.State = COMPLETE
	return nil
}

func (this *Tournament) AddResult(session dbr.SessionRunner, node *Node, winnerSide int) error {

	fmt.Printf("AddResult\n")

	if node.Left.node.Participant == 0 || node.Right.node.Participant == 0 {
		return errors.New("Game can't be completed, state error")
	}

	if node.Type == FINAL && winnerSide == 2 {
		if node.ChallengerUpOne == false {
			node.ChallengerUpOne = true
			return nil
		}
	}

	var winnerNodeId NodeId
	var winner ParticipantNumber
	var winnerPlayerId int
	var loser ParticipantNumber
	var loserPlayerId int
	if winnerSide == 1 {
		winner = node.Left.node.Participant
		winnerNodeId = node.Left.Id
		loser = node.Right.node.Participant
	} else {
		winner = node.Right.node.Participant
		winnerNodeId = node.Right.Id
		loser = node.Left.node.Participant
	}

	fmt.Printf("winner = %d, loser=%d\n", winner, loser)

	if winner > 0 {
		winnerPlayerId = this.Participants[winner].PlayerId
	}
	if loser > 0 {
		loserPlayerId = this.Participants[loser].PlayerId
	}

	fmt.Printf("setting node (%d) player to %d\n", node.Id, winner)
	node.Participant = winner

	var dropNodeId NodeId
	var dropNode *Node
	if node.Drop != 0 {
		dropNodeId = node.Drop
		dropNode = this.Bracket.Nodes[node.Drop]
		if dropNode != nil {
			fmt.Printf("dropNode id = %d, setting player to %d\n", dropNode.Id, loser)
			dropNode.Participant = loser
		}
	}

	res := new(MatchResult)
	res.WinningNode = winnerNodeId
	res.WinningParticipant = winner
	res.DropNode = dropNodeId
	res.LosingParticipant = loser
	res.WinningSlot = winnerSide
	res.WinningPlayer = winnerPlayerId

	res.Time = time.Now().Unix()
	node.GameState.Result = res

	//	_ = winnerPlayerId
	_ = loserPlayerId

	this.State = UNDERWAY
	this.StartTime = time.Now()
	this.GamesPlayed = this.GamesPlayed + 1

	return nil
}

func (this *Tournament) RemoveResult(session dbr.SessionRunner, subject string, gameId NodeId) error {
	node := this.Bracket.GetNode(gameId)

	if (node.Type == FINAL) && (node.GameState.Result == nil) {
		if node.ChallengerUpOne {
			node.ChallengerUpOne = false
			return nil
		}
	}

	fmt.Printf("node = %v\n", node)
	if node.GameState.Result == nil {
		return errors.New("unable to undo game result")
	}
	fmt.Printf("node.Parent = %v\n", node.Parent)
	if node.Parent.Id != 0 && node.Parent.node.GameState.Result != nil {
		return errors.New("unable to undo game result")
	}
	dropNode := this.Bracket.GetNode(node.GameState.Result.DropNode)

	if dropNode != nil {
		dropNode.Participant = 0
		dropParent := dropNode.Parent
		dropParent.node.Participant = 0
	}

	node.Participant = 0

	node.GameState.Result = nil
	this.GamesPlayed = this.GamesPlayed - 1
	if this.GamesPlayed == 0 {
		this.State = REGISTRATION
	}

	//	fmt.Printf("Tournament: Remove Result:   games played now = %d\n", this.GamesPlayed)
	//	if this.GamesPlayed == 0 {
	//		this.State = NEEDS_DRAW
	//	}
	return nil
}

func (this *Tournament) AddParticipant(session dbr.SessionRunner, playerId int64, name string) error {

	num := 0
	lenRecyclable := len(this.RecyclableParticipantNumbers)
	if lenRecyclable > 0 {
		num = int(this.RecyclableParticipantNumbers[lenRecyclable-1])
		this.RecyclableParticipantNumbers =
			append([]ParticipantNumber(nil),
				this.RecyclableParticipantNumbers[0:lenRecyclable-1]...)
	} else {
		num = len(this.Participants) + 1
	}

	p := new(Participant)
	//	p.PlayerId = this.participantNumberCount
	p.TournamentId = int(this.Id)
	p.PlayerId = int(playerId)
	p.Name = name
	p.PaidAmount = 0
	p.Paid = false
	p.Number = ParticipantNumber(num)

	//	if this.Bracket != nil && this.State == UNDERWAY {
	if this.Bracket.Root != nil && this.State == UNDERWAY {
		fmt.Printf("Trying to add to running tournamnt\n")
		result := this.Bracket.AddParticipantIfAble(this.Participants, p.Number)
		if result == false {
			return errors.New("Unable to add new player, sorry")
		}
	}
	this.Participants[ParticipantNumber(num)] = p
	this.ParticipantsByPlayerId[p.PlayerId] = p

	fmt.Printf("AddParticipant: size now = %d\n", len(this.Participants))
	return nil
}

func (this *Tournament) recalcPlayerIdMap() {

	this.ParticipantsByPlayerId = make(map[int]*Participant)
	for _, v := range this.Participants {
		this.ParticipantsByPlayerId[v.PlayerId] = v
	}

}
func (this *Tournament) RemoveParticipant(session dbr.SessionRunner, num ParticipantNumber) {

	fmt.Printf("RemoveParticipant: %d, pariticpant.length=%d\n", num, len(this.Participants))
	playerId := this.Participants[ParticipantNumber(num)].PlayerId
	fmt.Printf("RemoveParticipant: %d, playerId =%d\n", num, playerId)
	delete(this.Participants, ParticipantNumber(num))
	delete(this.ParticipantsByPlayerId, playerId)

	this.RecyclableParticipantNumbers =
		append(this.RecyclableParticipantNumbers, ParticipantNumber(num))

}

func (this *Tournament) RemoveParticipants(session dbr.SessionRunner, subject string, participants []int) error {
	fmt.Printf("RemoveParticpants: %v\n", participants)

	if this.State == UNDERWAY {
		return errors.New("Unable to delete participants after tournament is underway")
	}

	for _, p := range participants {
		this.RemoveParticipant(session, ParticipantNumber(p))
	}
	//this.State = READY_TO_START
	if err := this.Store(session, subject); err != nil {
		return errors.New(fmt.Sprintf("Store,Marshall Tournament failed: %v", err))
	}
	return nil
}

func (this *Tournament) Store(session dbr.SessionRunner, subject string) error {

	user, err := LoadUserBySubject(session, subject)
	userId := user.Id

	//fmt.Printf("Tournament.Store: id = %d\n", this.Id)
	var data []byte

	if data, err = json.Marshal(this); err != nil {
		log.Fatalf("Store,Marshall Tournament failed: %v", err)
		return err
	}

	pCount := len(this.Participants)
	//	fmt.Printf("this.RootNodeId = %v\n", this.bracket.RootNodeId)
	var id int64
	if this.Id == 0 {
		err = session.InsertInto("tournaments").
			Columns("tournament_name", "tournament_data", "user_id", "subject",
				"account_id", "creation_date", "tournament_state", "participant_count").
			Values(this.Name, data, userId, subject, 1, time.Now(), this.State, pCount).
			Returning("id").
			Load(&id)

		if err != nil {
			Logger.Fatalf("Insert Tournaments failed: %v", err)
			return err
		}
		this.Id = id
	} else {
		start := &(this.StartTime)
		if start.IsZero() {
			start = nil
		}
		end := &(this.StartTime)
		if end.IsZero() {
			end = nil
		}
		UpdateStmt := session.Update("tournaments").
			Set("tournament_data", data).
			Set("tournament_name", this.Name).
			Set("participant_count", pCount).
			Set("tournament_state", this.State)

		if start != nil && !start.IsZero() {
			UpdateStmt.Set("start_time", start)
		}
		if end != nil && !end.IsZero() {
			UpdateStmt.Set("end_time", end)
		}

		_, err := UpdateStmt.Where("id = ?", this.Id).Exec()
		if err != nil {
			Logger.Fatalf("Update Tournament failed: %v", err)
			return err
		}
	}

	return nil
}

func ListTournamentsBySubject(session dbr.SessionRunner, subject string, active bool) ([]*TournamentRecord, error) {

	fmt.Printf("selecting tournamentRecords for subject: %s\n", subject)
	fmt.Printf("db = %v\n", session)
	result, err := session.Select("id, account_id, user_id, subject, tournament_name, creation_date, start_time, end_time, final_game, tournament_state, participant_count").From("tournaments").
		Where("subject = ?", subject).OrderBy("id").Rows()

	if err != nil {
		Logger.Fatalf("Select Tournament failed: %v\n", err)
		return nil, err
	}

	var records = make([]*TournamentRecord, 0)

	//c, e := result.Columns()
	//fmt.Printf("result set = %v, err = %v\n", c, e)
	for {
		if result.Next() == false {
			if err := result.Close(); err != nil {
				return records, err
			} else {
				return records, nil
			}
		}
		var id int64
		var accountid int64
		var name string
		var userId int64
		var subject string
		var creationDate time.Time
		var startTime time.Time
		var endTime time.Time
		var finalGame int
		var state string
		var participantCount int
		//fmt.Printf("Scanning record\n")
		if err := result.Scan(&id, &accountid, &userId, &subject, &name, &creationDate, &startTime, &endTime, &finalGame, &state, &participantCount); err != nil {
			_ = result.Close()
			Logger.Fatalf("Scan tournament data failed: %v\n", err)
			return nil, err
		}
		rec := TournamentRecord{id, userId, subject, name,
			creationDate, startTime, endTime, state, participantCount}
		records = append(records, &rec)
	}
}
func ListTournaments(session dbr.SessionRunner, subject string, active bool) ([]*TournamentRecord, error) {

	result, err := session.Select("id, account_id, user_id, subject, tournament_name, creation_date, start_time, end_time, final_game, tournament_state, participant_count").From("tournaments").
		Where("subject = ?", subject).OrderBy("id").Rows()

	if err != nil {
		Logger.Fatalf("Select Tournament failed: %v\n", err)
		return nil, err
	}

	var records = make([]*TournamentRecord, 0)

	//c, e := result.Columns()
	//fmt.Printf("result set = %v, err = %v\n", c, e)
	for {
		if result.Next() == false {
			if err := result.Close(); err != nil {
				return records, err
			} else {
				return records, nil
			}
		}
		var id int64
		var accountid int64
		var userId int64
		var name string
		var sub string
		var creationDate time.Time
		var startTime time.Time
		var endTime time.Time
		var finalGame int
		var state string
		var participantCount int
		//fmt.Printf("Scanning record\n")
		if err := result.Scan(&id, &accountid, &userId, &sub, &name, &creationDate, &startTime, &endTime, &finalGame, &state, &participantCount); err != nil {
			_ = result.Close()
			Logger.Fatalf("Scan tournament data failed: %v\n", err)
			return nil, err
		}
		rec := TournamentRecord{id, userId, sub, name,
			creationDate, startTime, endTime, state, participantCount}
		records = append(records, &rec)
	}
}

func LoadTournament(session dbr.SessionRunner, id int64) (*Tournament, error) {

	result, err := session.Select("tournament_data").From("tournaments").Where("id = ?", id).Rows()
	if err != nil {
		Logger.Fatalf("Select Tournament failed: %v\n", err)
		return nil, err
	}
	if result.Next() == false {
		return nil, nil
	}

	var data []byte
	if err := result.Scan(&data); err != nil {
		Logger.Fatalf("Scan tournament data failed: %v\n", err)
		return nil, nil
	}

	result.Close()

	///	fmt.Printf("data = %v\n", data)
	var t *Tournament
	if err = json.Unmarshal(data, &t); err != nil {
		Logger.Fatalf("Unmarshal tournament data failed: %v", err)
		return nil, nil
	}

	//fmt.Printf("t = %v\n", t)
	t.Id = id

	t.Bracket.internalize()
	t.Bracket.tournament = t
	t.Bracket.Root = t.Bracket.Nodes[t.Bracket.RootNodeId]

	return t, nil

}
func (this *Tournament) getUnassignedParticipants() (error, []*Participant) {

	newPart := make([]*Participant, 0)

	for _, v := range this.Participants {
		if v.Number == 0 {
			newPart = append(newPart, v)
		}
	}
	return nil, nil
}

func (this *Tournament) DrawParticipants(session dbr.SessionRunner, subject string) error {

	vals := make([]int, len(this.Participants))
	for i := 0; i < len(this.Participants); i++ {
		vals[i] = i + 1
	}
	Shuffle(vals)

	keys := make([]int, 0)
	for k, _ := range this.Participants {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	newParticipants := make(map[ParticipantNumber]*Participant)
	index := 0
	for _, k := range keys {
		p := this.Participants[ParticipantNumber(k)]
		p.Number = ParticipantNumber(index)
		newParticipants[p.Number] = p
		index += 1
	}

	this.Participants = newParticipants
	this.recalcPlayerIdMap()
	this.Bracket = new(Bracket)
	if err := this.Store(session, subject); err != nil {
		return errors.New(fmt.Sprintf("Store,Marshall Tournament failed: %v", err))
	}
	return nil
}

func (this *Tournament) BuildBrackets(session dbr.SessionRunner, subject string) (*Tournament, error) {

	Logger.Infof("Enter\n")

	//	if this.Bracket.Root != nil {
	//		return this, nil
	//	}
	participantNumbers := make([]ParticipantNumber, len(this.Participants))
	i := 0
	for k := range this.Participants {
		participantNumbers[i] = k
		i++
	}

	sort.Slice(participantNumbers, func(i, j int) bool { return int(i) < int(j) })

	if len(this.Participants) < 2 {
		return nil, errors.New("Not enough players")
	}

	this.Bracket = NewBracket(this.Participants)
	this.Bracket.tournament = this
	this.Bracket.BuildDoubleElimBracket(this.Participants)

	//this.resolveBuys(session)

	//this.Bracket.CreateBrackets(participantNumbers, this.Typ)

	//	this.State = READY_TO_START
	if err := this.Store(session, subject); err != nil {
		return nil, errors.New(fmt.Sprintf("Store,Marshall Tournament failed: %v", err))
	}

	return this, nil

}

const layoutISO = "2006-01-02"

func NewTournament2(userId int64, subject string, accountId int64) *Tournament {

	Logger.Info("Enter")
	t := new(Tournament)
	t.AccountId = accountId
	t.UserId = userId
	t.Subject = subject
	t.Typ = DOUBLE_ELIMINATION
	t.Participants = make(map[ParticipantNumber]*Participant)
	t.ParticipantsByPlayerId = make(map[int]*Participant)
	t.RecyclableParticipantNumbers = make([]ParticipantNumber, 0)
	t.CreatedAt = time.Now()
	//t.Bracket = new(Bracket)
	t.Bracket = new(Bracket)
	t.State = REGISTRATION
	date := time.Now().Format(layoutISO)
	t.Name = "Tournament" + "_" + date
	return t

}

func (this Tournament) String() string {

	var b strings.Builder

	fmt.Fprintf(&b, "Tournament id: %d\n", this.Id)

	return b.String()
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

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(p1, p2 *Participant) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(participants []*Participant) {
	ps := &participantSorter{
		participants: participants,
		by:           by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// planetSorter joins a By function and a slice of Planets to be sorted.
type participantSorter struct {
	participants []*Participant
	by           func(p1, p2 *Participant) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *participantSorter) Len() int {
	return len(s.participants)
}

// Swap is part of sort.Interface.
func (s *participantSorter) Swap(i, j int) {
	s.participants[i], s.participants[j] = s.participants[j], s.participants[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *participantSorter) Less(i, j int) bool {
	return s.by(s.participants[i], s.participants[j])
}
