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

	go recorder.Run()
	go webapp.ListenAndServe()

	go func() {
		start := recorder.NewRequestStart()
		recorder.Enqueue(start)
		log.Printf("start response %T\n", <-start.ResponseChan)
	}()

	go func() {
		time.Sleep(5 * time.Second)
		log.Println("stopping")
		stop := recorder.NewRequestStop()
		recorder.Enqueue(stop)
		log.Printf("stop response %T\n", <-stop.ResponseChan)
	}()

	glib.NewMainLoop(nil).Run()
}

func init() {
	runtime.LockOSThread()
}
