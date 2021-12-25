package brackets

import (
	"testing"
)

func Test1(t *testing.T) {

	/*
		participants := make([]ParticipantNumber, 14)
		for i := 0; i < 14; i++ {
			participants[i] = ParticipantNumber(i + 1)
		}

		bracket := NewBracket(participants)
		root := bracket.BuildDoubleElimBracket()
		root.debugNode()

		fmt.Printf("%s\n", root.PrintTree())
	*/

}

/*
func Test1(t *testing.T) {

	db, mock, err := sqlmock.New()
	_ = mock
	require.NoError(t, err)

	conn := &dbr.Connection{
		DB:            db,
		EventReceiver: &dbr.NullEventReceiver{},
		Dialect:       dialect.MySQL,
	}
	sess := conn.NewSession(nil)

	tournament := NewTournament2(1, "test", 1)

	for i := 0; i < 14; i++ {
		name := fmt.Sprintf("Player%d", i+1)
		tournament.AddParticipant(sess, int64(i+1), name)
	}

	_, err = tournament.BuildBrackets2(sess, "test")

	root := tournament.Bracket.Root
	root.debugNode()

	fmt.Printf("%s\n", root.PrintTree())

}

*/
