package ui

import "github.com/hajimehoshi/ebiten/v2"

var knob1 *Knob

func InitUI() {
	knob1 = &Knob{}
}

func UpdateUI() {

}

func DrawUI(screen *ebiten.Image) {
	iop := &ebiten.DrawImageOptions{}
	iop.GeoM.Translate(100, 100)
	screen.DrawImage(knob1.GetKnob(), iop)
}
