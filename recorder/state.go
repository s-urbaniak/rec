package recorder

import (
	"errors"
	"log"
)

var queue = make(chan Request)

type stateFn func(Recorder) stateFn

// Start enqueues a start request and blocks until the recording starts
// or the request failed.
func Start() error {
	return request(RequestStart{})
}

// Stop enqueues a stop request and blocks until the recording stops
// or the request failed.
func Stop() error {
	return request(RequestStop{})
}

func request(v interface{}) error {
	req := Request{
		Value:        v,
		ResponseChan: make(chan Response),
	}
	Enqueue(req)

	res := <-req.ResponseChan
	switch res.(type) {
	case ResponseOK:
		return nil
	case ResponseError:
		return res.(error)
	}

	return errors.New("unknown response")
}

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
		msg := make(chan Msg)
		if err := r.Start(msg); err != nil {
			req.ResponseChan <- NewResponseError(err)
			break
		}
		req.ResponseChan <- ResponseOK{}
		return recording(msg)

	default:
		req.ResponseChan <- NewResponseErrorf("invalid request")
	}

	return stopped
}

func recording(msgChan chan Msg) stateFn {
	return func(r Recorder) stateFn {
		var level MsgLevel

		for {
			select {
			case msg := <-msgChan:
				switch msg.(type) {
				case MsgLevel:
					level = msg.(MsgLevel)

				case MsgEOS:
					if err := r.Reset(); err != nil {
						log.Fatal(err) // must not happen
					}
					return stopped
				}

			case req := <-queue:
				switch req.Value.(type) {
				case RequestStop:
					if err := r.Stop(); err != nil {
						req.ResponseChan <- NewResponseError(err)
						return stopped
					}
					req.ResponseChan <- ResponseOK{}

				case RequestLevel:
					req.ResponseChan <- Response(level)

				default:
					req.ResponseChan <- NewResponseErrorf("invalid request")
					return stopped
				}
			}
		}
	}
}
