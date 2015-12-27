package recorder

// Request holds a message to the recorder state machine.
// Value can be one of the Request* types.
// Responses will be transmitted over the given channel.
type Request struct {
	Value        interface{}
	ResponseChan chan Response
}

// RequestStart is the message to start recording.
// The response will be ResponseOK or ResponseError in case an error happened.
type RequestStart struct{}

// RequestStop is the message to stop recording.
// The response will be ResponseOK or ResponseError in case an error happened.
type RequestStop struct{}

// RequestLevel requests the current recording level.
// The response will be of type MsgLevel.
type RequestLevel struct{}

// NewRequest returns a new request with the given value
// and initializes the response channel.
func NewRequest(value interface{}) Request {
	return Request{
		Value:        value,
		ResponseChan: make(chan Response),
	}
}
