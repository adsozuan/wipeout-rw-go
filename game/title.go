package wipeout

import "github.com/adsozuan/wipeout-rw-go/engine"

type Title struct {
	titleImage      uint16
	startTime       float32
	hasShownAttract bool
	render          *engine.Render
	ui              *UI
}

func NewTitle(startTime float32, render *engine.Render) *Title {

	t := ImageGetTexture("data/textures/wiptitle.tim")
	Logger.Printf("Title texture index: %d", t)

	return &Title{
		titleImage:      t,
		startTime:       startTime,
		render:          render,
		hasShownAttract: false,
		ui:              NewUI(render),
	}
}

func (t *Title) Update() {
	t.render.SetView2d()
	t.render.Push2d(engine.NewVec2i(0, 0), t.render.Size(), engine.NewRGBA(128, 128, 128, 25), int(t.titleImage))
	t.ui.DrawText("Wi", engine.NewVec2i(100, 100), UITextSize16, UIColorAccent)
}
