package engine

import (
	"log"
	"math"
	"os"
)

const (
	SystemWindowName   = "Wipeout"
	SystemWindowWidth  = 320
	SystemWindowHeight = 240
)

var Logger *log.Logger

// System is the main system of the game
type System struct {
	timeReal   float64
	timeScaled float64
	timeScale  float64
	tickLast   float64
	cycleTime  float64
	platform   *PlatformSdl
	Render     *Render
}

func NewSystem(platform *PlatformSdl) *System {
	Logger = log.New(os.Stderr, "engine |", log.Ldate|log.Ltime)
	Logger.Printf("Init")
	InputInit()

	r := NewRender()
	r.Init(platform.GetScreenSize())

	return &System{
		timeReal:   platform.Now(),
		timeScaled: 0.0,
		timeScale:  1.0,
		tickLast:   0.0,
		cycleTime:  0.0,
		platform:   platform,
		Render:     r,
	}
}

func (s *System) Cleanup() {
	s.Render.Cleanup()
	InputCleanUp()
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

	// TODO: Update game logic here

	s.Render.FrameEnd(s.cycleTime)
	InputClear()
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
