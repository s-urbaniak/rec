package main

import (
	"log"
	"runtime"

	"github.com/s-urbaniak/glib"
	"github.com/s-urbaniak/rec/webapp"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	go webapp.ListenAndServe()
	glib.NewMainLoop(nil).Run()
}

func init() {
	runtime.LockOSThread()
}
