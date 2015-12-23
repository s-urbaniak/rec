package main

import (
	"errors"
	"log"

	"github.com/s-urbaniak/glib"
	"github.com/s-urbaniak/gst"
	"github.com/s-urbaniak/rec/webapp"
)

func initRecord() (*gst.Pipeline, error) {
	src := gst.ElementFactoryMake("alsasrc", "alsasrc")
	src.SetProperty("device", "hw:0")

	audioconvert := gst.ElementFactoryMake("audioconvert", "audioconvert")

	level := gst.ElementFactoryMake("level", "level")
	level.SetProperty("post-messages", true)

	audioresample := gst.ElementFactoryMake("audioresample", "audioresample")
	vorbisenc := gst.ElementFactoryMake("vorbisenc", "vorbisenc")
	vorbisenc.SetProperty("quality", 0.7)
	oggmux := gst.ElementFactoryMake("oggmux", "oggmux")

	filesink := gst.ElementFactoryMake("filesink", "filesink")
	filesink.SetProperty("location", "test.ogg")

	pl := gst.NewPipeline("pl")
	if ok := pl.Add(
		src,
		audioconvert,
		level,
		audioresample,
		vorbisenc,
		oggmux,
		filesink,
	); !ok {
		return nil, errors.New("adding elements to pipeline failed")
	}

	if ok := src.Link(
		audioconvert,
		level,
		audioresample,
		vorbisenc,
		oggmux,
		filesink,
	); !ok {
		return nil, errors.New("linking elements failed")
	}

	return pl, nil
}

type LevelMsg struct {
	Peak, Rms                        float64
	Timestamp, Duration, RunningTime uint64
}

func NewLevelMsg(params glib.Params) LevelMsg {
	arr := new(glib.ValueArray)
	arr.SetPtr(params["peak"].(glib.Pointer))
	peak := arr.GetNth(0).Get()

	arr = new(glib.ValueArray)
	arr.SetPtr(params["rms"].(glib.Pointer))
	rms := arr.GetNth(0).Get()

	return LevelMsg{
		Peak:      peak.(float64),
		Rms:       rms.(float64),
		Timestamp: params["timestamp"].(uint64),
	}
}

func onMessage(bus *gst.Bus, msg *gst.Message) {
	t := msg.GetType()
	name, params := msg.GetStructure()
	// println(fmt.Sprintf("msg %v", t.String()))

	switch {
	case t == gst.MESSAGE_ELEMENT && name == "level":
		level := NewLevelMsg(params)
		if level.Peak > 0 {
			level.Peak = 0
		}

		// println(fmt.Sprintf(
		// 	"name %q level %+v peak %f",
		// 	name, level, math.Pow(10.0, level.Peak/20.0),
		// ))
	case t == gst.MESSAGE_STATE_CHANGED:
		// new, old, _ := msg.ParseStateChanged()
		//println(fmt.Sprintf("state change old %v new %v", old, new))
	case t == gst.MESSAGE_STREAM_STATUS:
		//s := msg.ParseStreamStatus()
		//println(fmt.Sprintf("%v", s))
	case t == gst.MESSAGE_STREAM_START:
		//println(fmt.Sprintf("### stream start"))
	case t == gst.MESSAGE_ERROR:
		//err, debug := msg.ParseError()
		//println(fmt.Sprintf("err %v debug %q", err, debug))
		//err.Free()
	}
}

func main() {
	pl, err := initRecord()
	if err != nil {
		log.Fatal(err)
	}

	bus := pl.GetBus()
	bus.AddSignalWatch()
	bus.Connect("message", onMessage, nil)

	if state := pl.SetState(gst.STATE_PLAYING); state == gst.STATE_CHANGE_FAILURE {
		log.Fatal("state change failed")
	}

	go webapp.ListenAndServe()

	glib.NewMainLoop(nil).Run()
}
