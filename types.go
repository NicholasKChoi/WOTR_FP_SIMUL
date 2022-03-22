package main

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=FspPositionType

type FspPositionType int

const (
	// FspPositionTypes
	OutMordor FspPositionType = 0
	InMordor  FspPositionType = 1
)

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
	return h.Name == "eye" || (strings.HasSuffix(h.Name, "r") && !isGollumGuide)
}

func (h *HuntTile) Init(dmg int, name string) {

}

type Result struct {
	Frequency int
	Sum       int
}

type FspPosition struct {
	CurrPosition int
	ModeType     FspPositionType
}

type ResultTable map[int]*Result

func (t ResultTable) RegisterResult(value int) error {
	for i := 0; i <= value; i++ {
		if _, ok := t[i]; !ok {
			t[i] = &Result{}
		}
	}
	result := t[value]
	result.Frequency += 1
	return nil
}

func (t ResultTable) ToString(name string) string {
	s := "+" + strings.Repeat("-", 32) + "+"
	s += fmt.Sprintf("\n| %10s | %7s | %7s |", name, "Freq.", "Sum")
	sum := 0
	for i := 0; i < len(t); i++ {
		sum += t[i].Frequency
		s += fmt.Sprintf("\n| %10d | %7d | %7d |", i, t[i].Frequency, sum)
	}
	s += "\n+" + strings.Repeat("-", 32) + "+"
	return s
}

type TotalResultContainer struct {
	CorruptionTable ResultTable
	CharTable       ResultTable
	RevealTable     ResultTable
	AttacksTable    ResultTable
	RolledEyesTable ResultTable
	TurnsTable      ResultTable
}
