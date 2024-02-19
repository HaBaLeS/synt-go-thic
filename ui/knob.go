package ui

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math"
)

// Show What the knob is
// Type, linear, Log? (volume), Value A,B,C...
// Relative or fixed
// Attack, D, S ,R .. Base Freq (min/max),
// UI -- Where are we in % or A,B,C,D,...
// Scale
// Knob with shadow (gradient)

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
	//gg.NewRadialGradient()
	dc.SetRGB(100, 0, 0)
	dc.SetColor(colornames.Lime)
	dc.SetLineWidth(4)
	dc.DrawArc(64, 64, 32, math.Pi/1, 0)
	//dc.DrawCircle(64, 64, 32)
	dc.Stroke() //stroke or fill to commit drawing
	//dc.Fill()
	k.img = ebiten.NewImageFromImageWithOptions(dc.Image(), nil)
}
