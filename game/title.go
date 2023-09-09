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

	return &Title{
		startTime:       startTime,
		render:          render,
		hasShownAttract: false,
		ui:              NewUI(render),
	}
}

func (t *Title) Update() {
	t.render.SetView2d()
	t.ui.DrawText("Wi", engine.NewVec2i(100, 100), UITextSize16, UIColorAccent)
}
