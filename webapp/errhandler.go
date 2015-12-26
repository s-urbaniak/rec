package webapp

import (
	"log"
	"net/http"
)

// A Handler that responds to an HTTP request
// optionally returning an error.
type HandlerErr interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

// The HandlerErrFunc type is an adapter to allow the use of
// ordinary functions as HTTP error handlers.
type HandlerErrFunc func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerErrFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// A HandlerDecorator is a function that takes an http handler
// and returns an http handler.
type HandlerDecorator func(HandlerErr) HandlerErr

// Decorate decorates the given handler with all given decorators
// and returns a decorated handler.
func DecorateHandler(h HandlerErr, ds ...HandlerDecorator) HandlerErr {
	decorated := h
	for _, d := range ds {
		decorated = d(decorated)
	}
	return decorated
}

// logger returns a handler logging all errors in the given handler.
func logger(h HandlerErr) HandlerErr {
	return HandlerErrFunc(func(res http.ResponseWriter, req *http.Request) error {
		err := h.ServeHTTP(res, req)
		if err != nil {
			log.Println(err)
		}
		return err
	})
}

// Handler converts the given http error handler to an ordinary http handler
// ignoring any errors.
func HandlerIgnoreErr(h HandlerErr) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(res, req)
	})
}
