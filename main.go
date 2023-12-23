package main

import (
	"fmt"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"gitlab.com/gomidi/midi/v2"
	"log"
	"math"
	"time"
)

var pos = 0.0

const samplerate float64 = 48000.0
const freq float64 = 220.0

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
}

type Oscillator struct {
	ot   OscillatorType
	game *Game
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
	t := time.Now()
	fmt.Printf("WaveGen: %d\n", t.UnixMilli())
	sr := 44800.0            //Hz
	maxAmpl := 32767.0 * 0.5 //num
	noteA := 110.0           //Hz
	length := sr / noteA
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
	noteA := 110.0           //Hz
	length := sr / noteA

	for i := 0; i < len(p)/2; i++ {
		val := math.Sin(2 * math.Pi * (float64(pos) / length))
		ampl := int(val * maxAmpl)
		//fmt.Printf("(%d)%f ", pos, ampl)
		p[2*i] = byte(ampl)
		p[2*i+1] = byte(ampl >> 8)
		pos++
	}
	return len(p), nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	if g.playSound {
		ebitenutil.DebugPrintAt(screen, "Playing", 320/2, 240/2)
	} else {
		ebitenutil.DebugPrintAt(screen, "Pause", 320/2, 240/2)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	fmt.Println("synt-go-thic starting ....")

	//initMidi()
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

	game.osc = &Oscillator{
		ot:   RECT,
		game: game,
	}
	game.player = otoCtx.NewPlayer(game.osc)
	game.player.SetBufferSize(game.maxbuffer * 2)
	game.player.Play()
	fmt.Printf("Global Buffer is %dms with a buffer size of %d (%d bits, %d channels)\n", bufTimeMS, game.maxbuffer, bits, channels)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
	err = game.player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}

func initMidi() {

	defer midi.CloseDriver()

	in, err := midi.FindInPort("VMPK")
	if err != nil {
		fmt.Println("can't find VMPK")
		return
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, vel uint8
		switch {
		case msg.GetSysEx(&bt):
			fmt.Printf("got sysex: % X\n", bt)
		case msg.GetNoteStart(&ch, &key, &vel):
			fmt.Printf("starting note %s on channel %v with velocity %v\n", midi.Note(key), ch, vel)
		case msg.GetNoteEnd(&ch, &key):
			fmt.Printf("ending note %s on channel %v\n", midi.Note(key), ch)
		default:
			// ignore
		}
	}, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	time.Sleep(time.Second * 5)

	stop()
}

func (g *Game) Update() error {
	if inpututil.KeyPressDuration(ebiten.KeySpace) > 0 {
		g.playSound = true
	} else {
		g.playSound = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.osc.ot = SINE
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.osc.ot = RECT
	}
	return nil
}
