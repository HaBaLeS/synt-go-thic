package ui

//https://github.com/cettoana/go-waveform
import (
	"fmt"
	"github.com/HaBaLeS/synt-go-thic/engine"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

const (
	SR int = 48000
	W  int = 1000
	H  int = 300
)

type WaveFormUI struct {
	img        *ebiten.Image
	durationMs int
	sampleRate int
	width      int
	height     int
	buf        []float64
}

func (wf *WaveFormUI) Create() {

	//1sec -> 48000 lines
	//steam deck max:  1280x800
	// compression

	osc := engine.NewOscillator()
	osc.Frequency = 440
	osc.Velocity(100)
	osc.Harmonics(100)
	//var lastVal float64
	wf.buf = make([]float64, 0)
	for i := 0; i < W; i++ {
		lastVal := osc.AdvanceOscillatorSine()
		wf.buf = append(wf.buf, lastVal)
	}

}

func (wf *WaveFormUI) GetWaveForm() *ebiten.Image {
	if wf.img == nil {
		wf.genWaveForm()
	}
	return wf.img
}

func (wf *WaveFormUI) genWaveForm() {
	fmt.Print("Creating WaveForm\ns")
	dc := gg.NewContext(W, H)
	//gg.NewRadialGradient()
	dc.SetColor(colornames.Lightgray)
	dc.DrawRectangle(0, 0, float64(W), float64(H))
	dc.Fill()
	dc.SetColor(colornames.Orange)
	dc.SetLineWidth(1)
	for i, v := range wf.buf {
		//sy1 := 150.0
		sy2 := 150.0 + 150*v/2
		dc.DrawLine(float64(i), float64(sy2-1.0), float64(i), float64(sy2))
	}
	//dc.DrawCircle(64, 64, 32)
	dc.Stroke() //stroke or fill to commit drawing

	wf.img = ebiten.NewImageFromImageWithOptions(dc.Image(), nil)
}
