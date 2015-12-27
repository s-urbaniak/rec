package recorder

import (
	"errors"

	"github.com/s-urbaniak/grun"
	"github.com/s-urbaniak/gst"
)

// Recorder is the interface that defines the behavior of a recorder.
//
// Start starts the recorder. Pipeline messages occuring during recording
// will be transmitted over the given Msg channel.
//
// Stop stops the recorder.
// Note that this method does not stop the recorder immediately
// but a MsgEOS message will be sent via the Msg channel.
//
// Reset resets the recorder.
// After reset is invoked no messages via the Msg channel are about to happen.
type Recorder interface {
	Start(chan Msg) error
	Stop() error
	Reset() error
}

var _ Recorder = (*recorder)(nil)

type recorder struct {
	pl  *gst.Pipeline
	bus *gst.Bus
}

// NewRecorder returns a recorder which can be started.
func NewRecorder() Recorder {
	return &recorder{}
}

func (r *recorder) Start(msgChan chan Msg) (err error) {
	var (
		pl  *gst.Pipeline
		bus *gst.Bus
	)

	grun.Run(func() {
		src := gst.ElementFactoryMake("audiotestsrc", "src")
		//src.SetProperty("num-buffers", 100)
		//src.SetProperty("device", "hw:0")

		audioconvert := gst.ElementFactoryMake("audioconvert", "audioconvert")

		level := gst.ElementFactoryMake("level", "level")
		level.SetProperty("post-messages", true)
		level.SetProperty("interval", 100000000)

		audioresample := gst.ElementFactoryMake("audioresample", "audioresample")
		vorbisenc := gst.ElementFactoryMake("vorbisenc", "vorbisenc")
		vorbisenc.SetProperty("quality", 0.7)
		oggmux := gst.ElementFactoryMake("oggmux", "oggmux")

		filesink := gst.ElementFactoryMake("filesink", "filesink")
		filesink.SetProperty("sync", true)
		filesink.SetProperty("location", "test.ogg")

		pl = gst.NewPipeline("pl")
		if ok := pl.Add(
			src,
			audioconvert,
			level,
			audioresample,
			vorbisenc,
			oggmux,
			filesink,
		); !ok {
			err = errors.New("adding elements to pipeline failed")
			return
		}

		if ok := src.Link(
			audioconvert,
			level,
			audioresample,
			vorbisenc,
			oggmux,
			filesink,
		); !ok {
			err = errors.New("linking elements failed")
			return
		}

		if state := pl.SetState(gst.STATE_PLAYING); state == gst.STATE_CHANGE_FAILURE {
			err = errors.New("start failed: state change failed")
			return
		}

		bus = pl.GetBus()
		bus.AddSignalWatch()
		bus.Connect("message", NewOnMessageFunc(msgChan), nil)
	})

	if err != nil {
		return
	}

	r.pl = pl
	r.bus = bus

	return
}

func (r *recorder) Stop() error {
	var result bool

	grun.Run(func() {
		result = r.pl.SendEvent(gst.NewEventEOS())
	})

	if !result {
		return errors.New("stop failed (EOS event was not handled)")
	}

	return nil
}

func (r *recorder) Reset() (err error) {
	grun.Run(func() {
		if state := r.pl.SetState(gst.STATE_NULL); state == gst.STATE_CHANGE_FAILURE {
			err = errors.New("reset failed: state change failed")
		}
		r.bus.RemoveSignalWatch()
	})

	if err != nil {
		return
	}

	r.pl = nil
	r.bus = nil
	return
}
