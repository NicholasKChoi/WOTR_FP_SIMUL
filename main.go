package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	charR CharDiceNeededTable
	corrR CorruptionInflictedTable
	reavR RevealsTable

	huntPool []*HuntTile
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

func DoHunt(eyes int) {
	var tile *HuntTile
	if len(huntPool) == 0 {
		fmt.Println("Hunt Pool exhausted!")
	} else {
		tile, huntPool = huntPool[0], huntPool[1:]
		fmt.Printf("%v doing %d dmg and IsReveal %t\n",
			tile,
			tile.GetDmg(eyes),
			tile.IsReveal())
	}
}

func main() {
	fmt.Println("Hello world")
	huntPool = InitDefaultHuntPool()
	for len(huntPool) > 0 {
		DoHunt(2)
	}
}
