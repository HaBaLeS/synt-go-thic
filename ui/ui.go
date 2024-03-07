package ui

import (
	"github.com/HaBaLeS/synt-go-thic/mpk"
	"github.com/hajimehoshi/ebiten/v2"
)

var knobs []*UiKnob
var waveformUI *WaveFormUI

func InitUI(mpk3m *mpk.MPK3Mini) {
	knobs = make([]*UiKnob, 8)

	knobs[0] = NewUiKnob(mpk3m.MidiKnob("k1"))
	knobs[1] = NewUiKnob(mpk3m.MidiKnob("k2"))
	knobs[2] = NewUiKnob(mpk3m.MidiKnob("k3"))
	knobs[3] = NewUiKnob(mpk3m.MidiKnob("k4"))
	knobs[4] = NewUiKnob(mpk3m.MidiKnob("k5"))
	knobs[5] = NewUiKnob(mpk3m.MidiKnob("k6"))
	knobs[6] = NewUiKnob(mpk3m.MidiKnob("k7"))
	knobs[7] = NewUiKnob(mpk3m.MidiKnob("k8"))

	waveformUI = &WaveFormUI{}
	waveformUI.Create()
}

func UpdateUI() {

}

func DrawUI(screen *ebiten.Image) {

	for i := 0; i < 8; i++ {
		iop := &ebiten.DrawImageOptions{}
		if i < 4 {
			iop.GeoM.Translate(300+float64(i)*128, 100)
		} else {
			iop.GeoM.Translate(300+float64(i-4)*128, 228)
		}
		screen.DrawImage(knobs[i].GetKnob(), iop)
	}

	iop := &ebiten.DrawImageOptions{}
	iop.GeoM.Translate(100, 500)
	screen.DrawImage(waveformUI.GetWaveForm(), iop)
}
