package recorder_test

import (
	"testing"

	"github.com/s-urbaniak/rec/recorder"
)

func TestRequestStart(t *testing.T) {
	req := recorder.NewRequestStart()
	req.Type = recorder.RequestStop

	go recorder.Run()

	recorder.Enqueue(req)
	res := <-req.ResponseChan

	if res.Err == nil {
		t.Error(res.Err)
	}
}
