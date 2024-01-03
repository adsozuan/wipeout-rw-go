package system

import (
	"log"
	"math"
	"os"

	"github.com/adsozuan/wipeout-rw-go/engine"
	"github.com/adsozuan/wipeout-rw-go/game"
)

const (
	WindowName   = "Wipeout"
	WindowWidth  = 320
	WindowHeight = 240
)

var Logger *log.Logger

// System is the main system of the game
type System struct {
	timeReal   float64
	timeScaled float64
	timeScale  float64
	tickLast   float64
	cycleTime  float64
	platform   *engine.PlatformSdl
	Render     *engine.Render
	Game       *game.Game
}

func New(platform *engine.PlatformSdl) (*System, error) {

	Logger = log.New(os.Stderr, "system |", log.Ldate|log.Ltime)
	Logger.Printf("Init")

	engine.InputInit()

	r := engine.NewRender()
	r.Init(platform.GetScreenSize())

	g, err := game.NewGame(r, platform)
	if err != nil {
		return nil, err
	}

	return &System{
		timeReal:   platform.Now(),
		timeScaled: 0.0,
		timeScale:  1.0,
		tickLast:   0.0,
		cycleTime:  0.0,
		platform:   platform,
		Render:     r,
		Game: g,
	}, err
}

func (s *System) Cleanup() {
	s.Render.Cleanup()
	engine.InputCleanUp()
}

func (s *System) Exit() {
	s.platform.Exit()
}

func (s *System) Update() {
	timeRealNow := s.platform.Now()
	realDelta := timeRealNow - s.timeReal
	s.timeReal = timeRealNow
	s.tickLast = math.Min(realDelta, 0.1) * s.timeScale
	s.timeScaled += s.tickLast

	// FIXME: This is a hack to prevent the cycleTime from growing too large, must be a better way
	s.cycleTime = s.timeScaled
	if s.cycleTime > 3600*math.Pi {
		s.cycleTime -= 3600 * math.Pi
	}
	s.Render.FramePrepare()

	resetCycleTime := s.Game.Update()
	if resetCycleTime {
		s.ResetCycleTime()
	}

	s.Render.FrameEnd(s.cycleTime)
	engine.InputClear()
}

func (s *System) ResetCycleTime() {
	s.cycleTime = 0.0
}

func (s *System) Resize() {

}

func (s *System) TimeScale() float64 {
	return s.timeScale
}

func (s *System) TickLast() float64 {
	return s.tickLast
}

func (s *System) CycleTime() float64 {
	return s.cycleTime
}
func (s *System) Time() float64 {
	return s.timeScaled
}
