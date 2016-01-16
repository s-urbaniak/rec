package msg

import (
	"github.com/s-urbaniak/glib"
	"github.com/s-urbaniak/gst"
)

type Msg interface{}

type MsgLevel struct {
	Peak, Rms                        float64
	Timestamp, Duration, RunningTime uint64
}

type MsgDeviceAdded struct {
	Class       string
	Name        string
	DisplayName string
	Properties  map[string]interface{}
}

func NewMsgDeviceAdded(dev *gst.Device) MsgDeviceAdded {
	_, props := dev.GetProperties()

	return MsgDeviceAdded{
		Class:       dev.GetDeviceClass(),
		Name:        dev.GetName(),
		DisplayName: dev.GetDisplayName(),
		Properties:  props,
	}
}

type MsgUnknown string

type MsgEOS struct{}

type onMessageFunc func(bus *gst.Bus, msg *gst.Message)

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

func NewOnMessageFunc(msgChan chan Msg) onMessageFunc {
	return func(bus *gst.Bus, msg *gst.Message) {
		typ := msg.GetType()
		name, params := msg.GetStructure()

		switch {
		case typ == gst.MESSAGE_ELEMENT && name == "level":
			msgChan <- NewMsgLevel(params)
		case typ == gst.MESSAGE_DEVICE_ADDED:
			dev := msg.ParseDeviceAdded()
			defer dev.Unref()
			msgChan <- NewMsgDeviceAdded(dev)
		case typ == gst.MESSAGE_EOS:
			msgChan <- MsgEOS{}
		default:
			msgChan <- MsgUnknown(typ.String())
		}
	}
}
