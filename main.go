package main

import (
	"fmt"
	"unsafe"

	"github.com/fsamin/go-dump"
)

var (
	charR CharDiceNeededTable
	corrR CorruptionInflictedTable
	reavR RevealsTable
)

// hunt pool stuff
func InitDefaultHuntPool() []*HuntTile {
	h := []*HuntTile{
		{3, "3"}, {3, "3"}, {3, "3"}, {2, "2"}, {2, "2"}, {1, "1"}, {1, "1"},
		{-1, "eye"}, {-1, "eye"}, {-1, "eye"}, {-1, "eye"}, {2, "2r"},
		{2, "2r"}, {1, "1r"}, {1, "1r"}, {0, "0r"}, {0, "0r"},
	}
	return h
}

func main() {
	fmt.Println("Hello world")
	huntPool := InitDefaultHuntPool()
	dump.Dump(huntPool)
	fmt.Printf("%d\n", unsafe.Sizeof(huntPool[0]))
}
