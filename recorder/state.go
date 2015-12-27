package recorder

import "log"

var queue = make(chan Request)

type stateFn func(Recorder) stateFn

// Enqueue enqueues the given request to the recorder state machine.
func Enqueue(r Request) {
	queue <- r
}

// Run runs the recording state machine with the given recorder.
func Run(r Recorder) {
	for state := stopped; state != nil; {
		state = state(r)
	}
}

func stopped(r Recorder) stateFn {
	req := <-queue

	switch req.Value.(type) {
	case RequestStart:
	default:
		req.ResponseChan <- NewResponseErrorf("invalid request")
		return stopped
	}

	msg := make(chan Msg)
	if err := r.Start(msg); err != nil {
		req.ResponseChan <- NewResponseError(err)
		return stopped
	}

	req.ResponseChan <- ResponseOK{}
	return recording(msg)
}

func recording(msgChan chan Msg) stateFn {
	return func(r *recorder) stateFn {
		select {
		case msg := <-msgChan:
			log.Printf("msg %T %+v", msg, msg)

			switch msg.(type) {
			case MsgEOS:
				r.Stop()
				return stopped
			}

		case req := <-queue:
			switch req.Value.(type) {
			case RequestStop:
				grun.Run(func() { r.pl.SendEvent(gst.NewEventEOS()) })

			default:
				req.ResponseChan <- NewResponseErrorf("invalid request")
				return stopped
			}

			req.ResponseChan <- ResponseOK{}
		}

		return recording(msgChan)
	}
}
