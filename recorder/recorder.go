package recorder

import (
	"errors"

	"github.com/s-urbaniak/grun"
	"github.com/s-urbaniak/gst"
)

type recorder struct {
	pl  *gst.Pipeline
	bus *gst.Bus
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
			err = errors.New("record start failed: state change failed")
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

func (r *recorder) Stop() {
	grun.Run(func() {
		r.pl.SetState(gst.STATE_NULL)
		r.bus.RemoveSignalWatch()
	})

	r.pl = nil
	r.bus = nil
}
