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
