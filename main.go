package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var (
	charR CharDiceNeededTable
	corrR CorruptionInflictedTable
	reavR RevealsTable

	huntPool             []*HuntTile
	isFrodoRevealed      bool
	currentCorruptionLvl int
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
func DoHunt(eyes int) {
	var tile *HuntTile
	if len(huntPool) == 0 {
		fmt.Println("Hunt Pool exhausted!")
	} else {
		// randomize the hunt drawing at the time of drawing the tile
		tile, huntPool = huntPool[0], huntPool[1:]
		fmt.Printf("%v doing %d dmg and IsReveal %t\n",
			tile,
			tile.GetDmg(eyes),
			tile.IsReveal())
		isFrodoRevealed = tile.IsReveal()
		currentCorruptionLvl += tile.GetDmg(eyes)
	}
}

func main() {
	fmt.Println("Hello world")
	huntPool = InitDefaultHuntPool()
	for len(huntPool) > 0 {
		DoHunt(2)
		fmt.Printf("After that isFrodoRevealed = %t and dmg is at %d\n",
			isFrodoRevealed,
			currentCorruptionLvl)
		isFrodoRevealed = false
	}
	spew.Dump(huntPool)
}
