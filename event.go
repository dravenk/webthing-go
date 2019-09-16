package webthing

import (
	"encoding/json"
	"fmt"
)

// Event An Event represents an individual event from a thing.
type Event struct {
	thing *Thing
	name  string
	data  json.RawMessage
	time  string
}

// EventObject An event object describes a kind of event which may be emitted by a device.
// See https://iot.mozilla.org/wot/#event-object
type EventObject struct {
	AtType      string `json:"@type,omitempty"`
	Title       string `json:"title,omitempty"`
	ObjectType  string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Unit        string `json:"unit,omitempty"`
	Links       []Link `json:"links,omitempty"`
}

// NewEvent Initialize the object.
// @param thing Thing this event belongs to
// @param name  Name of the event
// @param data  Data associated with the event
func NewEvent(thing *Thing, name string, data json.RawMessage) *Event {
	return &Event{
		thing: thing,
		name:  name,
		data:  data,
		time:  Timestamp(),
	}
}

// AsEventDescription Get the event description.
// @return Description of the event as a JSONObject.
func (event *Event) AsEventDescription() []byte {
	eve := struct {
		Timestamp string          `json:"timestamp"`
		Data      json.RawMessage `json:"data,omitempty"`
	}{
		Timestamp: event.Time(),
		Data:      event.Data(),
	}
	base := make(map[string]interface{})
	base[event.Name()] = eve

	description, err := json.Marshal(base)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}

	return description

}

// Thing Get the thing associated with this event.
// @returns {Object} The thing.
func (event *Event) Thing() *Thing {
	return event.thing
}

// Name Get the event's name.
// @returns {String} The name.
func (event *Event) Name() string {
	return event.name
}

// Data Get the event's data.
// @returns {*} The data.
func (event *Event) Data() json.RawMessage {
	return event.data
}

// Time Get the event's timestamp.
// @returns {String} The time.
func (event *Event) Time() string {
	return event.time
}
