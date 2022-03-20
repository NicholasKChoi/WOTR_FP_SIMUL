package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type Simulator struct {
	Results  TotalResultContainer
	CurrGame *Game
}

func (s *Simulator) Init() error {
	if err := s.setupResults(); err != nil {
		return errors.Wrap(err, "failed to setup results table")
	}
	return nil
}

func (s *Simulator) RunSimulation() error {
	s.CurrGame = &Game{}
	if err := s.CurrGame.Init(GameLogLevel); err != nil {
		return errors.Wrap(err, "failed due to game could not init")
	} else if err = s.CurrGame.Run(); err != nil {
		return errors.Wrap(err, "game failed to run")
	}
	if err := s.CurrGame.UpdateResults(s.Results); err != nil {
		return errors.Wrap(err, "failed to update results with latest game")
	}
	return nil
}

func (s *Simulator) PrintResult() error {
	fmt.Println(s.Results.CharTable.ToString("CharDice"))
	fmt.Println(s.Results.CorruptionTable.ToString("Corrupt"))
	fmt.Println(s.Results.RevealTable.ToString("Reveal"))
	return nil
}

func (s *Simulator) setupResults() error {
	s.Results = TotalResultContainer{
		CorruptionTable: make(map[int]*Result),
		CharTable:       make(map[int]*Result),
		RevealTable:     make(map[int]*Result),
	}
	return nil
}
