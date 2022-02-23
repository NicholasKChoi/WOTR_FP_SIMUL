package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const (
	MORDOR_POSITION = 5
	MAX_CORRUPTION  = 12
)

var (
	charR CharDiceNeededTable
	corrR CorruptionInflictedTable
	reavR RevealsTable

	CHECKS = []VictoryChecker{&FrodoRingChecker{}, &ShadowRingChecker{}}

	// TODO: only put global settings here

	// TODO: move these to a gamestate structure that is in the mainloop
	HuntPool           []*HuntTile
	IsFrodoRevealed    bool
	CurrMovesUsed      int
	CurrCorruptionLvl  int
	CurrMordorPosition int
)

// hunt pool stuff
func InitDefaultHuntPool() []*HuntTile {
	h := []*HuntTile{
		{3, "3"}, {3, "3"}, {3, "3"}, {2, "2"}, {2, "2"}, {1, "1"}, {1, "1"},
		{-1, "eye"}, {-1, "eye"}, {-1, "eye"}, {-1, "eye"}, {2, "2r"},
		{2, "2r"}, {1, "1r"}, {1, "1r"}, {0, "0r"}, {0, "0r"},
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(h), func(i, j int) { h[i], h[j] = h[j], h[i] })
	return h
}

/// Does the following
/// 1. removes a tile from the hunt pool
/// 2. calculates damage based on the tile type
/// 3. determines if frodo is revealed based on tile type
func DoHunt(eyes int, huntPool []*HuntTile) (h []*HuntTile, t *HuntTile) {
	if len(huntPool) == 0 {
		fmt.Println("Hunt Pool exhausted!")
		h = huntPool
	} else {
		// randomize the hunt drawing at the time of drawing the tile
		t, h = huntPool[0], huntPool[1:]
		fmt.Printf("%v doing %d dmg and IsReveal %t\n",
			t,
			t.GetDmg(eyes),
			t.IsReveal())
	}
	return
}

// order of victory conditions matters
/// 1. if frodo > 12 corruption, shadow wins
/// 2. if frodo is at mount doom, frodo wins
/// 3. if shadow has 10 victory points, shadow wins
/// 4. if frodo has 4 victory points, frodo wins
type VictoryChecker interface {
	IsVictory(currMP int, currCorr int) bool
	PrintVictoryMessage()
}

type FrodoRingChecker struct{}
type ShadowRingChecker struct{}

func (*FrodoRingChecker) IsVictory(currMP int, currCorr int) bool {
	return currMP >= MORDOR_POSITION
}

func (*FrodoRingChecker) PrintVictoryMessage() {
	fmt.Println("FRODO DUNKS THE RING")
}

func (*ShadowRingChecker) IsVictory(currMP int, currCorr int) bool {
	return currCorr >= MAX_CORRUPTION
}

func (*ShadowRingChecker) PrintVictoryMessage() {
	fmt.Println("SAURON GETS THE RING")
}

func isGameOver(currMP int, currCorr int) bool {
	for _, checker := range CHECKS {
		if checker.IsVictory(CurrMordorPosition, CurrCorruptionLvl) {
			checker.PrintVictoryMessage()
			return true
		}
	}
	return false
}

func main() {

	// todo roll for number of eyes
	numEyes := 2
	fmt.Println("Hello world")
	HuntPool = InitDefaultHuntPool()
	for len(HuntPool) > 0 {
		CurrMovesUsed++
		if IsFrodoRevealed {
			fmt.Println("Frodo hid! Sneaky hobbits!")
			IsFrodoRevealed = false
		} else {
			fmt.Println("Frodo moved! HUNT HIM DOWN!!")
			var drawnTile *HuntTile
			HuntPool, drawnTile = DoHunt(numEyes, HuntPool)
			IsFrodoRevealed = drawnTile.IsReveal()
			CurrCorruptionLvl += drawnTile.GetDmg(numEyes)
			CurrMordorPosition++
		}
		fmt.Printf("After that isFrodoRevealed = %t and dmg is at %d\n",
			IsFrodoRevealed,
			CurrCorruptionLvl)

		if isGameOver(CurrMordorPosition, CurrCorruptionLvl) {
			break
		}
	}
	spew.Dump(HuntPool)
}
