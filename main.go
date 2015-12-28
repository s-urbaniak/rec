package main

import (
	"log"
	"runtime"
	"time"

	"github.com/s-urbaniak/glib"
	"github.com/s-urbaniak/rec/recorder"
	"github.com/s-urbaniak/rec/webapp"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	go webapp.ListenAndServe()

	go func() {
		rec := recorder.NewRecorder()
		msgChan := make(chan recorder.Msg)

		if err := rec.Start(msgChan); err != nil {
			log.Fatal(err)
		}

		timeout := time.After(5 * time.Second)

	loop:
		for {
			select {
			case msg := <-msgChan:
				log.Printf("%T %+v", msg, msg)
			case <-timeout:
				log.Printf("### TIMEOUT")
				break loop
			}
		}

		if err := rec.Stop(); err != nil {
			log.Fatal(err)
		}

	loop2:
		for {
			select {
			case msg := <-msgChan:
				if _, ok := msg.(recorder.MsgEOS); ok {
					break loop2
				}
			}
		}

		if err := rec.Reset(); err != nil {
			log.Fatal(err)
		}
		println("stopped")
	}()

	glib.NewMainLoop(nil).Run()
}

func init() {
	runtime.LockOSThread()
}
