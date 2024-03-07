package ui

import (
	"fmt"
	"github.com/HaBaLeS/synt-go-thic/mpk"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
)

// Show What the knob is
// Type, linear, Log? (volume), Value A,B,C...
// Relative or fixed
// Attack, D, S ,R .. Base Freq (min/max),
// UI -- Where are we in % or A,B,C,D,...
// Scale
// Knob with shadow (gradient)

type UiKnob struct {
	img  *ebiten.Image
	data *mpk.MidiKnob
}

func NewUiKnob(midiKnib *mpk.MidiKnob) *UiKnob {
	return &UiKnob{
		data: midiKnib,
	}
}

func (k *UiKnob) GetKnob() *ebiten.Image {
	//if k.img == nil {
	k.genKnob() //check if val changed
	//}
	return k.img
}

func (k *UiKnob) genKnob() {

	null := math.Pi - math.Pi/4
	max := 2*math.Pi + math.Pi/4

	//fmt.Print("Creating Knob\ns")
	dc := gg.NewContext(128, 128)
	//gg.NewRadialGradient()
	//dc.SetRGB(100, 0, 0)
	dc.SetColor(colornames.Purple)
	txt := fmt.Sprintf("%d", k.data.Value())
	dc.DrawStringAnchored(txt, 64, 64, 0.5, 0.5)
	dc.Stroke()

	dc.SetColor(color.White)
	dc.DrawEllipticalArc(64, 64, 48, 48, max, null)
	dc.Stroke()

	dc.SetColor(colornames.Lime)
	dc.SetLineWidth(4)

	rng := max - math.Pi*3/4
	cur := (rng * k.data.Percent() / 100) + null
	//fmt.Printf("PErc %f - pos: %f\n", k.data.Percent(), cur)

	dc.DrawEllipticalArc(64, 64, 50, 50, cur+0.1, cur-0.1)
	//dc.DrawCircle(64, 64, 32)
	dc.Stroke() //stroke or fill to commit drawing
	//dc.Fill()
	k.img = ebiten.NewImageFromImageWithOptions(dc.Image(), nil)
}
