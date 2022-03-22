package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/op/go-logging"
	"github.com/pkg/errors"
)

const (
	MOUNT_DOOM_POSITION     = 5
	MAX_CORRUPTION          = 12
	NUM_TILES_BEFORE_GOLLUM = 2
)

var (
	isGollumGuide        = false
	numTilesBeforeGollum = NUM_TILES_BEFORE_GOLLUM
)

var (
	CHECKS = []VictoryChecker{&ShadowRingChecker{}, &FrodoRingChecker{}}

	defaultHuntPool = []*HuntTile{
		{3, "3"}, {3, "3"}, {3, "3"}, {2, "2"}, {2, "2"}, {1, "1"}, {1, "1"},
		{-1, "eye"}, {-1, "eye"}, {-1, "eye"}, {-1, "eye"}, {2, "2r"},
		{2, "2r"}, {1, "1r"}, {1, "1r"}, {0, "0r"}, {0, "0r"},
	}

	initialSeed = time.Now().UnixNano()
)

type Game struct {
	FellowshipPosition FspPosition
	Corruption         int
	CharacterDieUsed   int
	HuntPool           []*HuntTile

	isFspRevealed bool
	logger        *logging.Logger

	numShadowDice    int
	numFreeDice      int
	numReveals       int
	numShadowAttacks int
	numRolledEyes    int
	numTurns         int
}

func (g *Game) Init(logLvl logging.Level) error {
	h := defaultHuntPool
	rand.Seed(initialSeed)
	initialSeed++
	rand.Shuffle(len(h), func(i, j int) { h[i], h[j] = h[j], h[i] })
	g.HuntPool = h
	g.FellowshipPosition = FspPosition{ModeType: InMordor}
	g.logger = logging.MustGetLogger("game")
	backend := logging.AddModuleLevel(logging.NewLogBackend(os.Stdout, "", 0))
	backend.SetLevel(logLvl, "")
	g.logger.SetBackend(backend)
	// default number + witch king + saruman
	g.numShadowDice = 9
	g.numFreeDice = 6
	g.numTurns = 0

	// reset globals
	numTilesBeforeGollum = NUM_TILES_BEFORE_GOLLUM
	isGollumGuide = false
	return nil
}

func (g *Game) Run() error {
	for !g.isGameOver() {
		if err := g.TakeTurn(); err != nil {
			return errors.Wrap(err, "take turn failed")
		}
	}
	return nil
}

func (g *Game) TakeTurn() error {
	g.logger.Debug(g.StateString())
	allocatedEyes := 1

	// roll shadow dice
	shadowDice := RollShadowDice(g.numShadowDice - allocatedEyes)
	g.numRolledEyes += shadowDice.NumEyes
	totalEyes, movesMade := allocatedEyes+shadowDice.NumEyes, 0
	g.numShadowAttacks += shadowDice.NumAttacks

	// roll free dice
	freeDice := RollFreeDice(g.numFreeDice)
	for i := 0; i < freeDice.NumMoves; i++ {
		if err := g.FspUseDie(totalEyes, movesMade); err != nil {
			return errors.Wrap(err, "fsp failed to make action")
		}
		if g.isGameOver() {
			break
		}
		g.logger.Debug("FSP moved and have to add to hunt pool")
		movesMade++
	}

	// first turn on mordor
	if g.numTurns == 0 {
		// muster mouth (assume that you rolled a muster)
		g.numShadowDice++
	}
	g.numTurns++
	g.logger.Debug(g.StateString())
	return nil
}

func (g *Game) StateString() string {
	_, _, line, _ := runtime.Caller(1)
	return fmt.Sprintf("LINE %d: Fsp='%s@%d' | Corr='%d' | Chars='%d' | Reveals='%d' | PoolSize='%d' | Turns='%d'",
		line, g.FellowshipPosition.ModeType, g.FellowshipPosition.CurrPosition, g.Corruption,
		g.CharacterDieUsed, g.numReveals, len(g.HuntPool), g.numTurns)
}

func (g *Game) FspUseDie(totalEyes, movesMade int) error {
	if err := g.moveFsp(totalEyes, movesMade); err != nil {
		return errors.Wrap(err, "movement error")
	}
	return nil
}

func (g *Game) UpdateResults(res TotalResultContainer) error {
	if err := res.CorruptionTable.RegisterResult(g.Corruption); err != nil {
		return errors.Wrap(err, "failed to register corruption")
	} else if err := res.CharTable.RegisterResult(g.CharacterDieUsed); err != nil {
		return errors.Wrap(err, "failed to register characters")
	} else if err := res.RevealTable.RegisterResult(g.numReveals); err != nil {
		return errors.Wrap(err, "failed to register reveals")
	} else if err := res.RolledEyesTable.RegisterResult(g.numRolledEyes); err != nil {
		return errors.Wrap(err, "failed to register rolled eyes")
	} else if err := res.AttacksTable.RegisterResult(g.numShadowAttacks); err != nil {
		return errors.Wrap(err, "failed to register shadow attacks")
	} else if err := res.TurnsTable.RegisterResult(g.numTurns); err != nil {
		return errors.Wrap(err, "failed to register number turns")
	}
	return nil
}

func (g *Game) moveFsp(numEyes, movesMade int) error {
	var numHits int
	g.CharacterDieUsed++
	if g.isFspRevealed {
		g.isFspRevealed = false
		g.logger.Debug("hiding")
	} else {
		if g.FellowshipPosition.ModeType == InMordor {
			numHits = numEyes + movesMade
		} else {
			panic("out of mordor not implemented")
		}
		drawnTile := g.drawHuntTile(numHits)
		g.logger.Debugf("Tile=%v: doing %d dmg and IsReveal %t\n",
			drawnTile,
			drawnTile.GetDmg(numHits),
			drawnTile.IsReveal())
		if drawnTile.GetDmg(numHits) > 0 {
			numTilesBeforeGollum -= 1
		}
		if numTilesBeforeGollum < 1 {
			isGollumGuide = true
		}
		g.isFspRevealed = drawnTile.IsReveal()
		g.Corruption += drawnTile.GetDmg(numHits)
		g.FellowshipPosition.CurrPosition++
		if g.isFspRevealed {
			g.numReveals++
		}
	}
	return nil
}

func (g *Game) isGameOver() bool {
	for _, checker := range CHECKS {
		if checker.IsVictory(g) {
			g.logger.Debug(checker.VictoryMessage())
			g.logger.Info(g.StateString())
			return true
		}
	}
	return false
}

/// Does the following
/// 1. removes a tile from the hunt pool
/// 2. calculates damage based on the tile type
/// 3. determines if frodo is revealed based on tile type
func (g *Game) drawHuntTile(eyes int) (t *HuntTile) {
	if len(g.HuntPool) == 0 {
		g.logger.Warning("Hunt Pool exhausted!")
	} else {
		// randomize the hunt drawing at the time of drawing the tile
		t, g.HuntPool = g.HuntPool[0], g.HuntPool[1:]
	}
	return
}

// order of victory conditions matters
/// 1. if frodo > 12 corruption, shadow wins
/// 2. if frodo is at mount doom, frodo wins
/// 3. if shadow has 10 victory points, shadow wins
/// 4. if frodo has 4 victory points, frodo wins
type VictoryChecker interface {
	IsVictory(game *Game) bool
	VictoryMessage() string
}

type FrodoRingChecker struct{}
type ShadowRingChecker struct{}

func (*FrodoRingChecker) IsVictory(game *Game) bool {
	return game.FellowshipPosition.CurrPosition >= MOUNT_DOOM_POSITION
}

func (*FrodoRingChecker) VictoryMessage() string {
	return "FRODO DUNKS THE RING"
}

func (*ShadowRingChecker) IsVictory(game *Game) bool {
	return game.Corruption >= MAX_CORRUPTION
}

func (*ShadowRingChecker) VictoryMessage() string {
	return "SAURON GETS THE RING"
}
