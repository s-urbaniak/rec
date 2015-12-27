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

	go recorder.Run(recorder.NewRecorder())
	go webapp.ListenAndServe()

	go func() {
		log.Printf("start err %v", recorder.Start())

		lr := recorder.NewRequest(recorder.RequestLevel{})
		for i := 0; i < 50; i++ {
			recorder.Enqueue(lr)
			log.Printf("level response %+v\n", <-lr.ResponseChan)
			time.Sleep(100 * time.Millisecond)
		}

		log.Printf("stop err %v", recorder.Stop())
	}()

	glib.NewMainLoop(nil).Run()
}

func init() {
	runtime.LockOSThread()
}
