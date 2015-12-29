package recorder

import (
	"errors"

	"github.com/s-urbaniak/gst"
)

// Recorder is the interface that defines the behavior of a recorder.
//
// Start starts the recorder. Pipeline messages occuring during recording
// will be transmitted over the given Msg channel.
//
// Stop stops the recorder.
// Note that this method does not stop the recorder immediately
// but a MsgEOS message will be sent eventually via the Msg channel.
//
// Reset resets the recorder.
// After reset is invoked no messages via the Msg channel are about to happen.
type Recorder interface {
	Start() error
	Stop() error
}

var _ Recorder = (*recorder)(nil)

type recorder struct {
	pl *gst.Pipeline
}

// NewRecorder returns a recorder which can be started.
func NewRecorder() Recorder {
	return &recorder{}
}

func (r *recorder) Start() (err error) {
	var pl *gst.Pipeline

	src := gst.ElementFactoryMake("pulsesrc", "src")
	//src := gst.ElementFactoryMake("audiotestsrc", "src")
	//src.SetProperty("num-buffers", 100)
	//src.SetProperty("device", "hw:2")

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

	if state := pl.SetState(gst.STATE_PLAYING); state == gst.STATE_CHANGE_FAILURE {
		return errors.New("start failed: state change failed")
	}

	pl.GetBus().AddSignalWatch()

	r.pl = pl
	return
}

func (r *recorder) MsgChan(msgChan chan Msg) {
	r.pl.GetBus().Connect("message", NewOnMessageFunc(msgChan), nil)
}

func (r *recorder) Stop() error {
	if ok := r.pl.SendEvent(gst.NewEventEOS()); !ok {
		return errors.New("stop failed (EOS event was not handled)")
	}

	// wait for EOS
	eosChan := make(chan Msg)
	r.MsgChan(eosChan)
	for ok := false; !ok; _, ok = (<-eosChan).(MsgEOS) {
	}

	var err error
	if state := r.pl.SetState(gst.STATE_NULL); state == gst.STATE_CHANGE_FAILURE {
		err = errors.New("reset failed: state change failed")
	}

	r.pl.GetBus().RemoveSignalWatch()
	r.pl = nil

	return err
}
