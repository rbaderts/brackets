package brackets

import (
	"testing"
)

func TestTierCalcs(t *testing.T) {

	tourn := NewTournament2()
	tourn.AddPlayers([]string{"A","B", "C", "D", "E", "F", "G"})
	//tourn.Size = 17
	_ = tourn.CreateDoubleElimBracket()
//	fmt.Printf("%v", root.PrintTree())

	/*
	for players := 16; players > 5; players-- {
		size := ComputeSize(players)
		fmt.Printf("TotalW: %d TotalL: %d, %d/%d   W1: %d, W2: %d, L1: %d, L2:%d\n",
			ComputeWTierTotal(players), ComputeLTierTotal(players), players, size,
			ComputeWTier1(size, players), ComputeWTier2(size, players),
			ComputeLTier1(size, players), ComputeLTier2(size, players))


	}
	 */



}
