package main

import "strings"

type HuntTile struct {
	Dmg  int
	Name string
}

func (h *HuntTile) GetDmg(numHits int) int {
	if h.Name == "eye" {
		return numHits
	} else {
		return h.Dmg
	}
}

func (h *HuntTile) IsReveal() bool {
	return h.Name == "eye" || strings.HasSuffix(h.Name, "r")
}

func (h *HuntTile) Init(dmg int, name string) {

}

type CharDiceNeededTable struct {
	Results []struct {
		NumCharDie int
		Frequency  int
		Sum        int
	}
}

type CorruptionInflictedTable struct {
	Results []struct {
		NumCorruptionInflicted int
		Frequency              int
		Sum                    int
	}
}

type RevealsTable struct {
	Results []struct {
		NumReveals int
		Frequency  int
		Sum        int
	}
}
