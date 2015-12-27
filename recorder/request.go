package recorder

// Request holds a message to the recorder state machine.
// Value can be one of the Request* types.
// Responses will be transmitted over the given channel.
type Request struct {
	Value        interface{}
	ResponseChan chan Response
}

// RequestStart is the message to start recording.
type RequestStart struct{}

// RequestStop is the message to stop recording.
type RequestStop struct{}

// RequestLevel requests the current recording level.
type RequestLevel struct{}

// NewRequestStart returns a new start request
// and initializes the response channel.
func NewRequestStart() Request {
	return Request{
		Value:        RequestStart{},
		ResponseChan: make(chan Response),
	}
}

// NewRequestStop returns a new stop request
// and initializes the response channel.
func NewRequestStop() Request {
	return Request{
		Value:        RequestStop{},
		ResponseChan: make(chan Response),
	}
}

// NewRequest returns a new request with the given value
// and initializes the response channel.
func NewRequest(value interface{}) Request {
	return Request{
		Value:        value,
		ResponseChan: make(chan Response),
	}
}
