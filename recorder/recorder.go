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
	if r.pl != nil {
		return errors.New("pipeline already configured")
	}

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

	if state := r.pl.SetState(gst.STATE_PLAYING); state == gst.STATE_CHANGE_FAILURE {
		return errors.New("state change failed")
	}

	r.bus = pl.GetBus()
	r.bus.AddSignalWatch()

	return nil
}

// MsgChan will advertise recorder events to the given channel.
// This method must be called after Start otherwise it will block forever.
func (r *Recorder) MsgChan(msgChan chan msg.Msg) {
	if r.bus == nil {
		return
	}

	r.bus.Connect("message", msg.NewOnMessageFunc(msgChan), nil)
}

// Stop stops the Recorder.
func (r *Recorder) Stop() error {
	if r.pl == nil {
		return nil
	}

	if ok := r.pl.SendEvent(gst.NewEventEOS()); !ok {
		return errors.New("stop failed (EOS event was not handled)")
	}

	// wait for EOS
	eosChan := make(chan msg.Msg)
	r.MsgChan(eosChan)
	for ok := false; !ok; _, ok = (<-eosChan).(msg.MsgEOS) {
	}

	// TODO(sur): not sure why this has to run in the main loop thread,
	// needs further investigation.
	// When not run in the main loop thread, the app seems to deadlock.
	return grunErr(func() error {
		if state := r.pl.SetState(gst.STATE_NULL); state == gst.STATE_CHANGE_FAILURE {
			return errors.New("state change failed")
		}
		return nil
	})

	r.bus.RemoveSignalWatch()
	r.bus.Unref()
	r.pl.Unref()
	r.bus = nil
	r.pl = nil
	return nil
}

func grunErr(f func() error) error {
	var err error
	grun.Run(func() {
		err = f()
	})
	return err
}
