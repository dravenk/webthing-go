package webthing

import (
	"fmt"
	"net/http"
)

// EventHandle handle a request to /events.
type EventHandle struct {
	*EventsHandle
	eventName string
}

// Handle a request to /events.
func (h *EventHandle) Handle(w http.ResponseWriter, r *http.Request) {
	BaseHandle(h, w, r)
}

// Get Handle a GET request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *EventHandle) Get(w http.ResponseWriter, r *http.Request) {
	if content := h.Thing.EventDescriptions(h.eventName); content != nil {
		if _, err := w.Write(content); err != nil {
			fmt.Print(err)
		}
	}
}
