package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/s-urbaniak/glib"
	"github.com/s-urbaniak/gst"
	"github.com/s-urbaniak/rec/msg"
	"github.com/s-urbaniak/rec/recorder"
	"github.com/s-urbaniak/rec/webapp"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	go webapp.ListenAndServe()

	msgChan := make(chan msg.Msg)
	go func() {
		println("waiting for device monitor msg")
		for msg := range msgChan {
			println(fmt.Sprintf("device monitor msg %+v", msg))
		}
	}()

	mon := gst.NewDeviceMonitor()
	bus := mon.GetBus()
	bus.AddSignalWatch()
	defer bus.RemoveSignalWatch()
	bus.Connect("message", msg.NewOnMessageFunc(msgChan), nil)
	bus.Unref()

	caps := gst.NewCapsEmptySimple("audio/x-raw")
	mon.AddFilter("Audio/Source", caps)
	caps.Unref()
	mon.Start()

	ds := mon.GetDevices()
	println(len(ds))
	for _, d := range ds {
		defer d.Unref()
		log.Printf("%q\n", d.GetDisplayName())
	}

	go func() {
		r := recorder.NewRecorder()
		r.Start()
		msgChan := make(chan msg.Msg)
		r.MsgChan(msgChan)

		for m := range msgChan {
			println(fmt.Sprintf("%+v\n", m))
		}
	}()

	println("starting main loop")
	glib.NewMainLoop(nil).Run()
}

func init() {
	runtime.LockOSThread()
}
