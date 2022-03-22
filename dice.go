package main

import "math/rand"

type Side int

const (
	MusterArmy Side = iota
	Muster
	Army
	Palantir
	Character
	Eye
	WillOfWest
)

type ShadowDieResults struct {
	NumAttacks int
	NumEyes    int
}

func RollShadowDice(numDice int) ShadowDieResults {
	sd := GetShadowDice()
	results := ShadowDieResults{}
	for i := 0; i < numDice; i++ {
		if side := sd.Roll(); side == MusterArmy || side == Army || side == Character {
			results.NumAttacks++
		} else if side == Eye {
			results.NumEyes++
		}
	}
	return results
}

type FreeDieResults struct {
	NumMoves int
}

func RollFreeDice(numDice int) FreeDieResults {
	fd := GetFreeDice()
	results := FreeDieResults{}
	for i := 0; i < numDice; i++ {
		if side := fd.Roll(); side == Character || side == WillOfWest {
			results.NumMoves++
		}
	}
	return results
}

type Dice struct {
	Sides []Side
}

func (d *Dice) Roll() Side {
	rand.Seed(initialSeed)
	initialSeed++
	s := d.Sides // reassign to variable in order to make next line shorter
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s[0]
}

func GetShadowDice() *Dice {
	return &Dice{[]Side{Muster, MusterArmy, Army, Palantir, Character, Eye}}
}

func GetFreeDice() *Dice {
	return &Dice{[]Side{WillOfWest, Character, Character, Muster, MusterArmy, Palantir}}
}
