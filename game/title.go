package wipeout

import "github.com/adsozuan/wipeout-rw-go/engine"

type Title struct {
	titleImage      uint16
	startTime       float32
	hasShownAttract bool
	render          *engine.Render
}

func NewTitle(startTime float32) *Title {

	return &Title{}

}
