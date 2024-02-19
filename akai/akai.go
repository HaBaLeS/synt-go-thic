package akai

import (
	"fmt"
	mpkminimk32 "github.com/HaBaLeS/synt-go-thic/akai/mpkminimk3"
	"gitlab.com/gomidi/midi/v2"
	"os"
)

func AkaiSendSettingsWithDriver() error {

	outPorts := midi.GetOutPorts()
	fmt.Printf("Found OUT Midi Device: %v", outPorts)
	midiOut := outPorts[1]
	err := midiOut.Open()
	defer midiOut.Close()
	if err != nil {
		return err
	}

	msg, err := mpkminimk32.CustomSettings(
		"fotze",
		10, // channel for pads TODO: this is -1 for some reason. we actually mean 9
		mpkminimk32.AutopopulatePads(36),
		mpkminimk32.AutopopulatePads(36), // TODO
		1,                                // channel for keys and knobs TODO: is this 1 or 0?
		mpkminimk32.KnobAbsolute0to127("KNOB 1", 70), // "K1"
		mpkminimk32.KnobAbsolute0to127("KNOB 2", 71),
		mpkminimk32.KnobRelative("i3 panel width", 72), // "K3"
		mpkminimk32.KnobAbsolute0to127("KNOB 4", 73),
		mpkminimk32.KnobAbsolute0to127("KNOB 5", 74),
		mpkminimk32.KnobAbsolute0to127("Sublime cols", 75),
		mpkminimk32.KnobRelative("i3 panel height", 76),
		mpkminimk32.KnobRelative("Volume", 77), // "K8"
	).SysExStore(mpkminimk32.Program1)

	f, _ := os.Create("/tmp/them.hex")
	f.Write(msg)
	f.Close()

	if err := midiOut.Send(msg); err != nil {
		return fmt.Errorf("error sending msg: %w", err)
	}

	return err
}

/*
func akaiEntry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "akai",
		Short: "",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			osutil.ExitIfError(akai(
				osutil.CancelOnInterruptOrTerminate(nil)))
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "send-settings",
		Short: "Send Hautomo settings to AKAI's RAM",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			osutil.ExitIfError(akaiSendSettings(
				osutil.CancelOnInterruptOrTerminate(nil)))
		},
	})

	return cmd
}*/

/*
func akai(ctx context.Context) error {
	drv, err := driver.New()
	if err != nil {
		return err
	}

	// make sure to close all open ports at the end
	defer drv.Close()

	in, err := midi.OpenIn(drv, -1, dev)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := midi.OpenOut(drv, -1, dev)
	if err != nil {
		return err
	}
	defer out.Close()

	//	if err := in.SetListener(func(data []byte, deltaMicroseconds int64) {
	//		log.Printf("MIDI inbound: %x", data)
	//	}); err != nil {
	//		return fmt.Errorf("SetListener: %w", err)
	//	}

	rd := reader.New(
		reader.NoLogger(),
		//	reader.SysEx(func(_ *reader.Position, data []byte) {
		//		log.Printf("MIDI inbound: %x", data)
		//	}),
		reader.Each(func(pos *reader.Position, msgGeneric midi.Message) {
			switch msg := msgGeneric.(type) {
			case sysex.SysEx:
				// dada does not contain header and footer, but does contain manufacturer byte!
				log.Printf("sysex %x", msg.Data())
			default:
				log.Printf("other MIDI message: %s", msgGeneric.String())
			}
		}),
	)

	if err := rd.ListenTo(in); err != nil {
		return err
	}

	// msg := requestConfigMsg()
	msg, err := mpkminimk3.ExampleSettings().SysExStore(mpkminimk3.Program8)
	if err != nil {
		return err
	}

	if _, err := out.Write(msg); err != nil {
		return fmt.Errorf("error sending msg: %w", err)
	}

	<-ctx.Done()

	return nil
}*/

/*
const (
	SNDRV_SEQ_IOCTL_PVERSION 0x40045300
	SNDRV_SEQ_IOCTL_CLIENT_ID 0x40045301
)

func akaiReceive(ctx context.Context)error{
	seq,err:=os.Open("/dev/snd/seq")
	if err!=nil{
		return
	}

	unix.IoctlSetInt(seq.Fd(),SNDRV_SEQ_IOCTL_PVERSION,)
openat(AT_FDCWD, "/dev/snd/seq", O_RDWR|O_CLOEXEC) = 3
ioctl(3, SNDRV_SEQ_IOCTL_PVERSION, 0x7fff9ca90fc8) = 0
ioctl(3, SNDRV_SEQ_IOCTL_CLIENT_ID, 0x7fff9ca90fcc) = 0
ioctl(3, SNDRV_SEQ_IOCTL_RUNNING_MODE, 0x7fff9ca90fd0) = 0
ioctl(3, SNDRV_SEQ_IOCTL_GET_CLIENT_INFO, 0x7fff9ca91290) = 0
ioctl(3, SNDRV_SEQ_IOCTL_SET_CLIENT_INFO, 0x7fff9ca91290) = 0
ioctl(3, SNDRV_SEQ_IOCTL_CREATE_PORT, 0x7fff9ca912a0) = 0
ioctl(3, SNDRV_SEQ_IOCTL_SUBSCRIBE_PORT, 0x7fff9ca91300) = 0
fcntl(3, F_GETFL)                       = 0x8002 (flags O_RDWR|O_LARGEFILE)
fcntl(3, F_SETFL, O_RDWR|O_NONBLOCK|O_LARGEFILE) = 0
}
*/
