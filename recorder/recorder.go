package recorder

import (
	"errors"

	"github.com/s-urbaniak/grun"
	"github.com/s-urbaniak/gst"
	"github.com/s-urbaniak/rec/msg"
)

type Recorder struct {
	pl       *gst.Pipeline
	bus      *gst.Bus
	location string
}

// NewRecorder returns a Recorder which can be started.
func NewRecorder(location string) *Recorder {
	return &Recorder{location: location}
}

// Start starts the Recorder. Pipeline messages occuring during recording
// will be transmitted over the given Msg channel.
func (r *Recorder) Start() error {
	var pl *gst.Pipeline

	src := gst.ElementFactoryMake("pulsesrc", "src")
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
	filesink.SetProperty("location", r.location)

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
		return errors.New("adding elements to pipeline failed")
	}

	if ok := src.Link(
		audioconvert,
		level,
		audioresample,
		vorbisenc,
		oggmux,
		filesink,
	); !ok {
		return errors.New("linking elements failed")
	}

	r.pl = pl

	if err := r.setState(gst.STATE_PLAYING); err != nil {
		return err
	}

	grun.Run(func() {
		r.bus = pl.GetBus()
		r.bus.AddSignalWatch()
	})

	return nil
}

func (r *Recorder) MsgChan(msgChan chan msg.Msg) {
	grun.Run(func() {
		r.bus.Connect("message", msg.NewOnMessageFunc(msgChan), nil)
	})
}

func (r *Recorder) setState(state gst.State) error {
	var err error
	grun.Run(func() {
		if state := r.pl.SetState(state); state == gst.STATE_CHANGE_FAILURE {
			err = errors.New("state change failed")
		}
	})
	return err
}

func (r *Recorder) sendEvent(evt *gst.Event) error {
	var err error
	grun.Run(func() {
		if ok := r.pl.SendEvent(evt); !ok {
			err = errors.New("stop failed (EOS event was not handled)")
		}
	})
	return err
}

// Stop stops the Recorder.
func (r *Recorder) Stop() error {
	r.sendEvent(gst.NewEventEOS())

	// wait for EOS
	eosChan := make(chan msg.Msg)
	r.MsgChan(eosChan)
	for ok := false; !ok; _, ok = (<-eosChan).(msg.MsgEOS) {
	}

	if err := r.setState(gst.STATE_NULL); err != nil {
		return err
	}

	grun.Run(func() {
		r.bus.RemoveSignalWatch()
		r.bus.Unref()
		r.pl.Unref()
		r.bus = nil
		r.pl = nil
	})

	return nil
}
