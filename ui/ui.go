package ui

import "github.com/hajimehoshi/ebiten/v2"

var knob1 *Knob
var waveformUI *WaveFormUI

func InitUI() {
	knob1 = &Knob{}
	waveformUI = &WaveFormUI{}
	waveformUI.Create()
}

func UpdateUI() {

}

func DrawUI(screen *ebiten.Image) {
	iop := &ebiten.DrawImageOptions{}
	iop.GeoM.Translate(100, 100)
	screen.DrawImage(knob1.GetKnob(), iop)

	iop = &ebiten.DrawImageOptions{}
	iop.GeoM.Translate(100, 300)
	screen.DrawImage(waveformUI.GetWaveForm(), iop)
}
