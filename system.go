package wipeout

import "math"

const (
	SystemWindowName   = "Wipeout"
	SystemWindowWidth  = 1280
	SystemWindowHeight = 720
)

// System is the main system of the game
type System struct {
	timeReal   float64
	timeScaled float64
	timeScale  float64
	tickLast   float64
	cycleTime  float64
	plaform    *PlatformSdl
	render     *Render
}

func NewSystem(platform *PlatformSdl, render *Render) *System {
	InputInit()

	r := NewRender()
	r.Init(platform.GetScreenSize())

	return &System{
		timeReal:   platform.Now(),
		timeScaled: 0.0,
		timeScale:  1.0,
		tickLast:   0.0,
		cycleTime:  0.0,
		plaform:    platform,
		render:     render,
	}
}

func (s *System) Cleanup() {
	s.render.Cleanup()
	InputCleanUp()
}

func (s *System) Exit() {
	s.plaform.Exit()
}

func (s *System) Update() {
	timeRealNow := s.plaform.Now()
	realDelta := timeRealNow - s.timeReal
	s.timeReal = timeRealNow
	s.tickLast = math.Min(realDelta, 0.1) * s.timeScale
	s.timeScaled += s.tickLast

	// FIXME: This is a hack to prevent the cycleTime from growing too large, must be a better way
	s.cycleTime = s.timeScaled
	if s.cycleTime > 3600*math.Pi {
		s.cycleTime -= 3600 * math.Pi
	}
	s.render.FramePrepare()

	// TODO: Update game logic here

	s.render.FrameEnd(s.cycleTime)
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
