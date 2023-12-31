package main

import (
	"fmt"
	"github.com/HaBaLeS/synt-go-thic/mpk"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	//_ "gitlab.com/gomidi/midi/v2/drivers/midicatdrv"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	"log"
	"math"
	"time"
)

var pos = 0.0

/*
27.5 -> A0
55 -> A1
110 -> A2
220 -> A3
440 -> A4

Base Frequ is
--> 27,5 *2^n



*/

var version string
var buildtime string
var deckBuild string

const samplerate float64 = 48000.0

const octaveBaseFreq = 110.0

var toneFrequIncrment float64 = math.Pow(2, 1.0/12.0)

type OscillatorType int

const (
	SINE OscillatorType = iota
	RECT
)

type Game struct {
	playSound bool
	osc       *Oscillator
	maxbuffer int //1chan * 48000/s * 16
	player    *oto.Player
	midi      *mpk.MPK3Mini
}

type Oscillator struct {
	ot        OscillatorType
	game      *Game
	frequency float64
}

func (o *Oscillator) Read(p []byte) (n int, err error) {

	if !o.game.playSound {
		for i := 0; i < len(p)/2; i++ {
			p[2*i] = 0x00
			p[2*i+1] = 0x00
		}
		return len(p), nil
	}

	switch o.ot {
	case SINE:

		return o.genSine(p)
	case RECT:
		return o.genRect(p)
	default:
		panic(fmt.Errorf("not implemented %d", o.ot))
	}

}

func (o *Oscillator) genRect(p []byte) (n int, err error) {
	sr := 44800.0            //Hz
	maxAmpl := 32767.0 * 0.5 //num
	length := sr / o.frequency
	ampl := 0

	for i := 0; i < len(p)/2; i++ {
		val := math.Sin(2 * math.Pi * (float64(pos) / length))
		if val > 0 {
			ampl = int(maxAmpl)
		} else {
			ampl = int(-maxAmpl)
		}
		p[2*i] = byte(ampl)
		p[2*i+1] = byte(ampl >> 8)
		pos++
	}
	return len(p), nil
}

func (o *Oscillator) genSine(p []byte) (n int, err error) {

	sr := 44800.0            //Hz
	maxAmpl := 32767.0 * 0.5 //num
	length := sr / o.frequency

	for i := 0; i < len(p)/2; i++ {
		val := math.Sin(2 * math.Pi * (float64(pos) / length))
		ampl := int(val * maxAmpl)
		p[2*i] = byte(ampl)
		p[2*i+1] = byte(ampl >> 8)
		pos++
	}
	return len(p), nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.playSound {
		ebitenutil.DebugPrintAt(screen, "Playing", 320/2, 240/2)
	} else {
		ebitenutil.DebugPrintAt(screen, "Pause", 320/2, 240/2)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K1 -> %d", g.midi.KnobVal("K1")), 20, 100)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K2 -> %d", g.midi.KnobVal("K2")), 20, 116)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K3 -> %d", g.midi.KnobVal("K3")), 20, 132)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K4 -> %d", g.midi.KnobVal("K4")), 20, 148)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K5 -> %d", g.midi.KnobVal("K5")), 20, 164)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K6 -> %d", g.midi.KnobVal("K6")), 20, 178)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K7 -> %d", g.midi.KnobVal("K7")), 20, 192)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K8 -> %d", g.midi.KnobVal("K8")), 20, 208)

	for i, v := range g.midi.MidiKeys() {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Piano Tone %s", v.Name), 200, 100+16*i)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	fmt.Println("synt-go-thic starting ....")

	var bufTimeMS time.Duration = 5
	bits := 16
	channels := 1

	op := &oto.NewContextOptions{}

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	op.SampleRate = int(samplerate)

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	op.ChannelCount = channels

	// Format of the source.
	op.Format = oto.FormatSignedInt16LE
	op.BufferSize = bufTimeMS * time.Millisecond

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	game := &Game{
		maxbuffer: int(samplerate/1000.0) * int(bufTimeMS) * (bits / 8) * channels,
	}

	fmt.Printf("Creating Midi Device\n")
	game.midi = mpk.NewMPK3Mini()

	game.osc = &Oscillator{
		ot:        RECT,
		game:      game,
		frequency: 440,
	}

	fmt.Printf("Creating Audio Plyer\n")
	game.player = otoCtx.NewPlayer(game.osc)
	game.player.SetBufferSize(game.maxbuffer * 2)
	game.player.Play()
	fmt.Printf("Global Buffer is %dms with a buffer size of %d (%d bits, %d channels)\n", bufTimeMS, game.maxbuffer, bits, channels)

	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("Synth-GO-thic")
	if deckBuild == "yes" {
		ebiten.SetFullscreen(true)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
	game.midi.Stop()
	err = game.player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}

func (g *Game) Update() error {
	if inpututil.KeyPressDuration(ebiten.KeySpace) > 0 {
		g.playSound = true
	} else {
		g.playSound = false
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.osc.ot = SINE
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.osc.ot = RECT
	}

	/*if inpututil.KeyPressDuration(ebiten.KeyA) > 0 {
		g.osc.frequency = octaveBaseFreq * math.Pow(toneFrequIncrment, 0)
		g.playSound = true
	} else if inpututil.KeyPressDuration(ebiten.KeyS) > 0 {
		g.osc.frequency = octaveBaseFreq * math.Pow(toneFrequIncrment, 1)
		g.playSound = true
	} else if inpututil.KeyPressDuration(ebiten.KeyD) > 0 {
		g.osc.frequency = octaveBaseFreq * math.Pow(toneFrequIncrment, 2)
		g.playSound = true
	} else if inpututil.KeyPressDuration(ebiten.KeyF) > 0 {
		g.osc.frequency = octaveBaseFreq * math.Pow(toneFrequIncrment, 3)
		g.playSound = true
	} else {
		g.playSound = false
	}*/

	midiKeys := g.midi.MidiKeys()
	if len(midiKeys) > 0 {
		k := midiKeys[0]
		bf := math.Pow(2, float64(k.Octave)) * 27.5 //frequ of A
		g.osc.frequency = bf * math.Pow(toneFrequIncrment, float64(k.Number-9))
		g.playSound = true
	} else {
		g.playSound = false
	}

	return nil
}
