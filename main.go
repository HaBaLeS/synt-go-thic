package main

import (
	"fmt"
	"github.com/HaBaLeS/synt-go-thic/ui"

	"github.com/HaBaLeS/synt-go-thic/engine"
	"github.com/HaBaLeS/synt-go-thic/mpk"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	//_ "gitlab.com/gomidi/midi/v2/drivers/midicatdrv"
	"log"
	"math"
	"time"

	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
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

type Game struct {
	maxbuffer       int //1chan * 48000/s * 16
	player          *oto.Player
	midi            *mpk.MPK3Mini
	mixer           *engine.Mixer
	keyToChannelMap map[int]int
}

func (g *Game) Draw(screen *ebiten.Image) {

	ui.DrawUI(screen)

	if g.midi == nil {
		return
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Key Press additional Gain -> %d", g.midi.KnobVal("K1")), 20, 100)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Waveform Select K2 -> %d", g.midi.KnobVal("K2")), 20, 116)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf(">%s<", g.mixer.WaveFormName), 200, 116)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("# Harmonics -> %d", g.midi.KnobVal("K3")), 20, 132)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K4 -> %d", g.midi.KnobVal("K4")), 20, 148+100)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K5 -> %d", g.midi.KnobVal("K5")), 20, 164+100)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K6 -> %d", g.midi.KnobVal("K6")), 20, 178+100)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K7 -> %d", g.midi.KnobVal("K7")), 20, 192+100)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("K8 -> %d", g.midi.KnobVal("K8")), 20, 208+100)

	/*for i, v := range g.midi.MidiKeys() {
		if v != nil {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Piano Tone %s", v.Name), 200, 100+16*i)
		}
	}*/

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
		maxbuffer:       int(samplerate/1000.0) * int(bufTimeMS) * (bits / 8) * channels,
		keyToChannelMap: make(map[int]int),
	}

	fmt.Printf("Creating Midi Device\n")
	game.midi, err = mpk.NewMPK3Mini()
	if err != nil {
		log.Printf("error %v", err)
	}

	game.mixer = engine.NewMixer(10, game.midi)

	fmt.Printf("Creating Audio Plyer\n")
	game.player = otoCtx.NewPlayer(game.mixer)
	game.player.SetBufferSize(game.maxbuffer * 2)
	game.player.Play()
	fmt.Printf("Global Buffer is %dms with a buffer size of %d (%d bits, %d channels)\n", bufTimeMS, game.maxbuffer, bits, channels)

	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("Synth-GO-thic")
	if deckBuild == "yes" {
		ebiten.SetFullscreen(true)
	}

	ui.InitUI(game.midi)

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

	ui.UpdateUI()

	if g.midi == nil {
		return nil
	}
	midiKeys := g.midi.MidiKeys()
	for _, k := range midiKeys {
		if k.EventType == mpk.EventPress {
			for i := 0; i < g.mixer.MaxChannels; i++ {
				c := g.mixer.Channel(i)
				if c.Frequency == 0.0 {
					g.keyToChannelMap[k.Idx] = i
					bf := math.Pow(2, float64(k.Octave)) * 27.5 //frequ of A
					frequ := bf * math.Pow(toneFrequIncrment, float64(k.Number-9))
					c.Frequency = frequ
					c.Velocity(k.Velocity)
					break
				}
			}
		} else if k.EventType == mpk.EventRelease {
			cn := g.keyToChannelMap[k.Idx]
			g.mixer.Channel(cn).Frequency = 0.0
			delete(g.keyToChannelMap, k.Idx)
			fmt.Printf("KM %v", g.keyToChannelMap)
		}
	}
	g.midi.ClearEvents()

	return nil
}
