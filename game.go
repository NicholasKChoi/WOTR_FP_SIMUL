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
	MORDOR_POSITION = 5
	MAX_CORRUPTION  = 12
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
	NumReveals         int
	HuntPool           []*HuntTile

	isFspRevealed bool
	logger        *logging.Logger
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
	if err := g.FspUseDie(); err != nil {
		return errors.Wrap(err, "fsp failed to make action")
	}
	g.logger.Debug(g.StateString())
	return nil
}

func (g *Game) StateString() string {
	_, _, line, _ := runtime.Caller(1)
	return fmt.Sprintf("LINE %d: Fsp='%s@%d' | Corr='%d' | Chars='%d' | Reveals='%d' | PoolSize='%d'",
		line, g.FellowshipPosition.ModeType, g.FellowshipPosition.CurrPosition, g.Corruption,
		g.CharacterDieUsed, g.NumReveals, len(g.HuntPool))
}

func (g *Game) FspUseDie() error {
	numEyes := 3
	if err := g.moveFsp(numEyes); err != nil {
		return errors.Wrap(err, "movement error")
	}
	return nil
}

func (g *Game) UpdateResults(res TotalResultContainer) error {
	if err := res.CorruptionTable.RegisterResult(g.Corruption); err != nil {
		return errors.Wrap(err, "failed to register corruption")
	} else if err := res.CharTable.RegisterResult(g.CharacterDieUsed); err != nil {
		return errors.Wrap(err, "failed to register corruption")
	} else if err := res.RevealTable.RegisterResult(g.NumReveals); err != nil {
		return errors.Wrap(err, "failed to register corruption")
	}
	return nil
}

func (g *Game) moveFsp(numEyes int) error {
	g.CharacterDieUsed++
	if g.isFspRevealed {
		g.isFspRevealed = false
		g.logger.Debug("hiding")
	} else {
		drawnTile := g.doHunt(numEyes)
		g.logger.Debugf("Tile=%v: doing %d dmg and IsReveal %t\n",
			drawnTile,
			drawnTile.GetDmg(numEyes),
			drawnTile.IsReveal())
		g.isFspRevealed = drawnTile.IsReveal()
		g.Corruption += drawnTile.GetDmg(numEyes)
		g.FellowshipPosition.CurrPosition++
		if g.isFspRevealed {
			g.NumReveals++
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
func (g *Game) doHunt(eyes int) (t *HuntTile) {
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
	return game.FellowshipPosition.CurrPosition >= MORDOR_POSITION
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
