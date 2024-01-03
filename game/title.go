package game

import "github.com/adsozuan/wipeout-rw-go/engine"

type TitleScene struct {
	titleImage      uint16
	startTime       float64
	hasShownAttract bool
	render          *engine.Render
	ui              *UI
}

func NewTitle(startTime float64, render *engine.Render) *TitleScene {

	return &TitleScene{
		startTime:       startTime,
		render:          render,
		hasShownAttract: false,
		ui:              NewUI(render),
	}
}


func (t *TitleScene) Init() error {
	texture := ImageGetTexture("data/textures/wiptitle.tim")
	t.titleImage = texture

	return nil
}

func (t *TitleScene) Update() error {
	t.render.SetView2d()
	err := t.render.Push2d(engine.NewVec2i(0, 0), t.render.Size(), engine.NewRGBA(128, 128, 128, 25), int(t.titleImage))
	if err != nil {
		return err
	}
	t.ui.DrawText("PRESS ENTER", t.ui.ScaledPos(UIPosBottom|UIPosCenter, engine.NewVec2i(0, -40)), UITextSize16, UIColorDefault)
	// t.ui.DrawImage(engine.NewVec2i(0, 0), int(t.titleImage))

	return nil
}
