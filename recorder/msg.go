package recorder

import (
	"github.com/s-urbaniak/glib"
	"github.com/s-urbaniak/gst"
)

type Msg interface{}

type MsgLevel struct {
	Peak, Rms                        float64
	Timestamp, Duration, RunningTime uint64
}

func NewMsgLevel(params glib.Params) MsgLevel {
	arr := new(glib.ValueArray)
	arr.SetPtr(params["peak"].(glib.Pointer))
	peak := arr.GetNth(0).Get()

	arr = new(glib.ValueArray)
	arr.SetPtr(params["rms"].(glib.Pointer))
	rms := arr.GetNth(0).Get()

	return MsgLevel{
		Peak:      peak.(float64),
		Rms:       rms.(float64),
		Timestamp: params["timestamp"].(uint64),
	}
}

type onMessageFunc func(bus *gst.Bus, msg *gst.Message)

func NewOnMessageFunc(msgChan chan Msg) onMessageFunc {
	return func(bus *gst.Bus, msg *gst.Message) {
		typ := msg.GetType()
		name, params := msg.GetStructure()

		switch {
		case typ == gst.MESSAGE_ELEMENT && name == "level":
			msgChan <- NewMsgLevel(params)
		}
	}
}
