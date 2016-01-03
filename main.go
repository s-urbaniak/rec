package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/s-urbaniak/glib"
	"github.com/s-urbaniak/rec/devicemon"
	"github.com/s-urbaniak/rec/msg"
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

	devicemon.Start()
	devicemon.MsgChan(msgChan)

	glib.NewMainLoop(nil).Run()
}

func init() {
	runtime.LockOSThread()
}
