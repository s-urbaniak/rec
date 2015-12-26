package recorder

import (
	"errors"
	"log"

	"github.com/s-urbaniak/grun"
	"github.com/s-urbaniak/gst"
)

var queue = make(chan Request)

type recorder struct {
	pl      *gst.Pipeline
	msgChan chan Msg
}

type stateFn func(*recorder) stateFn

func Enqueue(r Request) {
	queue <- r
}

func Run() {
	recoder := &recorder{}
	for state := stopped; state != nil; {
		state = state(recoder)
	}
}

func stopped(r *recorder) stateFn {
	req := <-queue

	switch req.Value.(type) {
	case RequestStart:
	default:
		req.ResponseChan <- NewResponseErrorf("invalid request")
		return stopped
	}

	var (
		pl      *gst.Pipeline
		err     error
		msgChan = make(chan Msg)
	)

	grun.Run(func() {
		src := gst.ElementFactoryMake("alsasrc", "alsasrc")
		src.SetProperty("device", "hw:0")

		audioconvert := gst.ElementFactoryMake("audioconvert", "audioconvert")

		level := gst.ElementFactoryMake("level", "level")
		level.SetProperty("post-messages", true)

		audioresample := gst.ElementFactoryMake("audioresample", "audioresample")
		vorbisenc := gst.ElementFactoryMake("vorbisenc", "vorbisenc")
		vorbisenc.SetProperty("quality", 0.7)
		oggmux := gst.ElementFactoryMake("oggmux", "oggmux")

		filesink := gst.ElementFactoryMake("filesink", "filesink")
		filesink.SetProperty("location", "test.ogg")

		pl = gst.NewPipeline("pl")
		if ok := pl.Add(
			src,
			audioconvert,
			level,
			audioresample,
			vorbisenc,
			oggmux,
			filesink,
		); !ok {
			err = errors.New("adding elements to pipeline failed")
			return
		}

		if ok := src.Link(
			audioconvert,
			level,
			audioresample,
			vorbisenc,
			oggmux,
			filesink,
		); !ok {
			err = errors.New("linking elements failed")
			return
		}

		if state := pl.SetState(gst.STATE_PLAYING); state == gst.STATE_CHANGE_FAILURE {
			err = errors.New("record start failed: state change failed")
			return
		}

		bus := pl.GetBus()
		bus.AddSignalWatch()
		bus.Connect("message", NewOnMessageFunc(msgChan), nil)
	})

	if err != nil {
		req.ResponseChan <- NewResponseError(err)
		return stopped
	}

	req.ResponseChan <- ResponseOK{}
	r.pl = pl
	r.msgChan = msgChan

	return recording
}

func recording(r *recorder) stateFn {
	select {
	case msg := <-r.msgChan:
		log.Printf("msg %T %+v", msg, msg)

	case req := <-queue:
		switch req.Value.(type) {
		case RequestStop:
		default:
			req.ResponseChan <- NewResponseErrorf("invalid request")
			return stopped
		}

		r.pl = nil
		close(r.msgChan)
	}

	return recording
}
