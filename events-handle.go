package webthing

import (
	"fmt"
	"net/http"
)

// EventsHandle a request to /events.
type EventsHandle struct {
	*ThingHandle
}

// Handle a request to /events.
func (h *EventsHandle) Handle(w http.ResponseWriter, r *http.Request) {
	if name, err := resource(r.RequestURI); err == nil {
		eventHandle := &EventHandle{h, name}
		eventHandle.Handle(w, r)
		return
	}
	BaseHandle(h, w, r)
}

// Get Handle a GET request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *EventsHandle) Get(w http.ResponseWriter, r *http.Request) {
	if content := h.Thing.EventDescriptions(""); content != nil {
		if _, err := w.Write(content); err != nil {
			fmt.Println(err)
		}
	}
}
