package ui

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type Knob struct {
	img *ebiten.Image
}

func (k *Knob) GetKnob() *ebiten.Image {
	if k.img == nil {
		k.genKnob()
	}
	return k.img
}

func (k *Knob) genKnob() {
	fmt.Print("Creating Knob\ns")
	dc := gg.NewContext(128, 128)
	gg.NewRadialGradient()
	dc.SetRGB(100, 0, 0)
	dc.SetLineWidth(4)
	dc.DrawCircle(64, 64, 32)

	//dc.Fill()
	k.img = ebiten.NewImageFromImageWithOptions(dc.Image(), nil)
}
