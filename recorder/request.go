package recorder

type Request struct {
	Value        interface{}
	ResponseChan chan Response
}

type RequestStart struct{}

type RequestStop struct{}

func NewRequestStart() Request {
	return Request{
		Value:        RequestStart{},
		ResponseChan: make(chan Response),
	}
}

func NewRequestStop() Request {
	return Request{
		Value:        RequestStop{},
		ResponseChan: make(chan Response),
	}
}
