package devicemon

import (
	"github.com/s-urbaniak/gst"
	"github.com/s-urbaniak/rec/msg"
)

var mon = gst.NewDeviceMonitor()

func Start() {
	bus := mon.GetBus()
	defer bus.Unref()

	bus.AddSignalWatch()
	mon.Start()
}

func MsgChan(msgChan chan msg.Msg) {
	bus := mon.GetBus()
	defer bus.Unref()

	devs := mon.GetDevices()
	for _, dev := range devs {
		defer dev.Unref()
		msgChan <- msg.NewMsgDeviceAdded(dev)
	}

	bus.Connect("message", msg.NewOnMessageFunc(msgChan), nil)
}

func Stop() {
	bus := mon.GetBus()
	defer bus.Unref()

	bus.RemoveSignalWatch()
}
