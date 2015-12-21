package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/s-urbaniak/glib"
	"github.com/s-urbaniak/gst"
)

func initRecord() (*gst.Pipeline, error) {
	src := gst.ElementFactoryMake("pulsesrc", "pulsesrc")
	// src.SetProperty("device", "hw:1")

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
		return nil, errors.New("elements not accepted")
	}

	src.Link(
		audioconvert,
		level,
		audioresample,
		vorbisenc,
		oggmux,
		filesink,
	)

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
	println(fmt.Sprintf("msg %v", t.String()))

	switch {
	case t == gst.MESSAGE_ELEMENT && name == "level":
		level := NewLevelMsg(params)

		println(fmt.Sprintf(
			"name %q level msg %+v",
			name,
			level,
		))

	case t == gst.MESSAGE_STATE_CHANGED:
		new, old, _ := msg.ParseStateChanged()
		println(fmt.Sprintf("state change old %v new %v", old, new))
	case t == gst.MESSAGE_STREAM_STATUS:
		s := msg.ParseStreamStatus()
		println(fmt.Sprintf("%v", s))
	case t == gst.MESSAGE_STREAM_START:
		println(fmt.Sprintf("### stream start"))
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

	glib.NewMainLoop(nil).Run()
}
