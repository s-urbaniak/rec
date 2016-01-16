package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/s-urbaniak/glib"
	"github.com/s-urbaniak/rec/devicemon"
	"github.com/s-urbaniak/rec/msg"
	"github.com/s-urbaniak/rec/recorder"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.Ltime | log.Lshortfile)

	msgChan := make(chan msg.Msg)
	go func() {
		for msg := range msgChan {
			println(fmt.Sprintf("msg %+v", msg))
		}
	}()

	devicemon.Start()
	devicemon.MsgChan(msgChan)

	go func() {
		for {
			println("### starting new recorder")
			r := recorder.NewRecorder(RandString(10) + ".ogg")
			if err := r.Start(); err != nil {
				panic(err)
			}
			r.MsgChan(msgChan)
			time.Sleep(time.Second)
			println("### done, stopping")
			if err := r.Stop(); err != nil {
				panic(err)
			}
			println("### stopped")
		}
	}()

	glib.NewMainLoop(nil).Run()
}

func init() {
	runtime.LockOSThread()
}
