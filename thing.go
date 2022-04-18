package webthing

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/xeipuuv/gojsonschema"
)

// Thing A Web Thing struct.
type Thing struct {
	id               string
	context          string
	atType           []string
	title            string
	description      string
	properties       map[string]*Property
	availableActions map[string]*AvailableAction
	availableEvents  map[string]*AvailableEvent
	actions          map[string][]*Action
	events           []*Event
	subscribers      map[string]*websocket.Conn
	hrefPrefix       string
	uiHref           string
}

// ThingMember thingmember
type ThingMember struct {
	ID          string                     `json:"id"`
	Context     string                     `json:"@context"`
	AtType      []string                   `json:"@type"`
	Title       string                     `json:"title"`
	Description string                     `json:"description,omitempty"`
	Properties  json.RawMessage            `json:"properties,omitempty"`
	Actions     map[string]json.RawMessage `json:"actions,omitempty"`
	Events      map[string]json.RawMessage `json:"events,omitempty"`
	Links       []Link                     `json:"links"`
}

func NewThingMember(thing *Thing) *ThingMember {
	th := &ThingMember{
		ID:          thing.id,
		Title:       thing.title,
		Context:     thing.context,
		AtType:      thing.atType,
		Description: thing.description,
		Properties:  json.RawMessage{},
		Actions:     make(map[string]json.RawMessage),
		Events:      make(map[string]json.RawMessage),
	}
	return th
}

// NewThing create a thing.
func NewThing(id, title string, atType []string, description string) *Thing {
	thing := &Thing{}
	thing.id = id
	thing.title = title
	thing.context = "https://webthings.io/schemas"
	thing.atType = atType
	thing.description = description
	thing.properties = make(map[string]*Property)
	thing.availableActions = make(map[string]*AvailableAction)
	thing.availableEvents = make(map[string]*AvailableEvent)
	thing.actions = make(map[string][]*Action)
	thing.events = []*Event{}
	thing.subscribers = map[string]*websocket.Conn{}
	thing.hrefPrefix = ""
	thing.uiHref = ""
	return thing
}

// Link base link struct
type Link struct {
	Href      string `json:"href,omitempty"`
	Rel       string `json:"rel,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
}

func (th *ThingMember) availableActionsDesc(thing *Thing) {
	for name := range thing.availableActions {
		meta := thing.availableActions[name].Metadata()
		var m map[string]interface{}
		json.Unmarshal(meta, &m)
		m["links"] = []Link{{
			Rel:  "action",
			Href: filepath.Clean(fmt.Sprintf("/%s/actions/%s", thing.Href(), name)),
		}}
		obj, _ := json.Marshal(m)
		th.Actions[name] = obj
	}
}

func (th *ThingMember) availableEventsDesc(thing *Thing) {
	for name := range thing.availableEvents {
		meta, _ := thing.availableEvents[name].Metadata().MarshalJSON()
		var m map[string]interface{}
		json.Unmarshal(meta, &m)
		m["links"] = []Link{{
			Rel:  "events",
			Href: filepath.Clean(fmt.Sprintf("/%s/events/%s", thing.Href(), name)),
		}}
		obj, _ := json.Marshal(m)
		th.Events[name] = obj
	}
}

func (th *ThingMember) links(thing *Thing) {
	for _, name := range []string{"properties", "actions", "events"} {
		th.Links = append(th.Links, Link{
			Rel:  name,
			Href: fmt.Sprintf("%s/%s", thing.hrefPrefix, name),
		})
	}

	if thing.UIHref() != "" {
		th.Links = append(th.Links, Link{
			Rel:       "alternate",
			MediaType: "text/html",
			Href:      filepath.Clean(thing.UIHref()),
		})
	}
}

// AsThingDescription retrun []byte data of thing struct.
// Return the thing state as a Thing Description.
// @returns {Object} Current thing state
func (thing *Thing) AsThingDescription() []byte {

	th := NewThingMember(thing)

	th.Properties = []byte(thing.PropertyDescriptions())
	th.availableActionsDesc(thing)
	th.availableEventsDesc(thing)
	th.links(thing)

	thingDescription, err := json.Marshal(th)
	if err != nil {
		fmt.Println(err.Error())
	}

	return thingDescription
}

// Href Get this thing's href.
//
// @returns {String} The href.
func (thing *Thing) Href() string {
	if thing.hrefPrefix != "" {
		return thing.hrefPrefix
	}

	return "/"
}

// UIHref Get this thing's UI href.
//
// @returns {String|null} The href.
func (thing *Thing) UIHref() string {
	return thing.uiHref
}

// SetHrefPrefix Set the prefix of any hrefs associated with this thing.
//
// @param {String} prefix The prefix
func (thing *Thing) SetHrefPrefix(prefix string) {
	thing.hrefPrefix = prefix
	for name := range thing.properties {
		thing.properties[name].SetHrefPrefix(prefix)
	}
	for name := range thing.actions {
		for key := range thing.actions[name] {
			thing.actions[name][key].SetHrefPrefix(prefix)
		}
	}
}

// SetUIHref Set the href of this thing's custom UI.
//
// @param {String} href The href
func (thing *Thing) SetUIHref(href string) {
	thing.uiHref = href
}

// ID Get the ID of the thing.
//
// @returns {String} The ID.
func (thing *Thing) ID() string {
	return thing.id
}

// Title Get the title of the thing.
//
// @returns {String} The title.
func (thing *Thing) Title() string {
	return thing.title
}

// Context Get the type context of the thing.
//
// @returns {String} The contexthing.
func (thing *Thing) Context() string {
	return thing.context
}

// Type Get the type(s) of the thing.
//
// @returns {String[]} The type(s).
func (thing *Thing) Type() []string {
	return thing.atType
}

// Description Get the description of the thing.
//
// @returns {String} The description.
func (thing *Thing) Description() string {
	return thing.description
}

// PropertyDescriptions Get the thing's properties as an object.
//
// @returns {Object} Properties, i.e. name -> description
func (thing *Thing) PropertyDescriptions() string {
	descriptions := make(map[string]json.RawMessage)
	for name, property := range thing.properties {
		descriptions[name] = []byte(property.AsPropertyDescription())
	}

	str, _ := json.Marshal(descriptions)
	return string(str)
}

// ActionDescriptions Get the thing's actions as an array.
//
// @param {String?} actionName Optional action name to get descriptions for
// @returns {Object} Action descriptions.
func (thing *Thing) ActionDescriptions(actionName string) []json.RawMessage {
	var descriptions []json.RawMessage
	if actionName == "" {
		for name := range thing.actions {
			for _, action := range thing.actions[name] {
				descriptions = append(descriptions, []byte(action.AsActionDescription()))
			}
		}
	} else {
		if actions, ok := thing.actions[actionName]; ok {
			for _, action := range actions {
				if action != nil {
					descriptions = append(descriptions, []byte(action.AsActionDescription()))
				}
			}
		}
	}

	return descriptions
}

// EventDescriptions Get the thing's events as an array.
//
//@param {String?} eventName Optional event name to get descriptions for
//
//@returns {Object} Event descriptions.
func (thing *Thing) EventDescriptions(eventName string) []byte {
	var descriptions []json.RawMessage
	if len(thing.events) == 0 {
		return []byte(`{}`)
	}
	for _, event := range thing.events {
		if eventName == "" || strings.EqualFold(event.Name(), eventName) {
			descriptions = append(descriptions, event.AsEventDescription())
		}
	}

	content, _ := json.Marshal(descriptions)
	return content
}

// AddProperty Add a property to this thing.
//
// @param property Property to add.
func (thing *Thing) AddProperty(property *Property) {
	property.SetHrefPrefix(thing.hrefPrefix)
	thing.properties[property.Name()] = property
}

// RemoveProperty Remove a property from this thing.
//
// @param property Property to remove.
func (thing *Thing) RemoveProperty(property Property) {
	if p, ok := thing.properties[property.Name()]; ok {
		delete(thing.properties, p.Name())
	}
}

// Find a property by name.
//
// @param propertyName Name of the property to find
// @return Property if found, else null.
func (thing *Thing) findProperty(propertyName string) (*Property, bool) {
	if p, ok := thing.properties[propertyName]; ok {
		return p, true
	}
	return &Property{}, false
}

// Property Get a property's value.
//
// @param propertyName Name of the property to get the value of
// @param <T>          Type of the property value
// @return Current property value if found, else null.
func (thing *Thing) Property(propertyName string) *Value {
	if prop, ok := thing.findProperty(propertyName); ok {
		return prop.Value()
	}
	return &Value{}
}

// Properties et a mapping of all properties and their values.
//
// @return JSON object of propertyName -&gt; value.
func (thing *Thing) Properties() map[string]interface{} {
	properties := make(map[string]interface{})
	for name, property := range thing.properties {
		properties[name] = property.Value().Get()
	}
	return properties
}

// Determine whether or not this thing has a given property.
//
// @param propertyName The property to look for
// @return Indication of property presence.
func (thing *Thing) HasProperty(propertyName string) bool {
	if _, ok := thing.properties[propertyName]; ok {
		return true
	}
	return false
}

// SetProperty Set a property value.
//
// @param propertyName Name of the property to set
// @param value        Value to set
// @param <T>          Type of the property value
// @throws PropertyError If value could not be set.
func (thing *Thing) SetProperty(propertyName string, value *Value) error {
	if _, ok := thing.findProperty(propertyName); !ok {
		return errors.New(`"General property error"`)
	}
	property := thing.properties[propertyName]
	return property.SetValue(value)
}

// Action Get an action.
//
// @param actionName Name of the action
// @param actionId   ID of the action
// @return The requested action if found, else null.
func (thing *Thing) Action(actionName, actionID string) (action *Action) {
	if _, ok := thing.actions[actionName]; !ok {
		return nil
	}
	for _, ac := range thing.actions[actionName] {
		// Each newly created action must contain a new uuid,
		// otherwise a random action with the same uuid will be found and returned.
		if ac != nil && ac.ID() == actionID {
			action = ac
		}
	}
	return action
}

// AddEvent Add a new event and notify subscribers.
//
// @param event The event that occurred.
func (thing *Thing) AddEvent(event *Event) {
	thing.events = append(thing.events, event)
	thing.EventNotify(event)
}

// AddAvailableEvent Add an available event.
//
// @param name     Name of the event
// @param metadata Event metadata, i.e. type, description, etc., as a
//                 JSONObject
func (thing *Thing) AddAvailableEvent(name string, metadata json.RawMessage) {
	thing.availableEvents[name] = NewAvailableEvent(metadata)
}

// PerformAction Perform an action on the thing.
//
// @param actionName Name of the action
// @param input      Any action inputs
// @return The action that was created.
func (thing *Thing) PerformAction(actionName string, input *json.RawMessage) (*Action, error) {
	if _, ok := thing.availableActions[actionName]; !ok {
		fmt.Print("Not found action: ", actionName)
		return nil, errors.New("Not found action: " + actionName)
	}

	actionType := thing.availableActions[actionName]

	schemaLoader := gojsonschema.NewGoLoader(actionType.schema)
	documentLoader := gojsonschema.NewGoLoader(input)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// if result.Valid() {
	// 	fmt.Printf("The document is valid\n")
	// }
	if !result.Valid() {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return nil, err
	}
	// if !actionType.ValidateActionInput(input) {
	// 	return nil
	// }

	cls := actionType.getCls()

	// The Generator is called to create an action.
	action := cls.Generator(thing)
	action.SetInput(input)
	action.SetHrefPrefix(thing.hrefPrefix)

	thing.ActionNotify(action)
	thing.actions[actionName] = append(thing.actions[actionName], action)

	return action, nil
}

// RemoveAction Remove an existing action.
//
// @param actionName name of the action
// @param actionId   ID of the action
// @return Boolean indicating the presence of the action.
func (thing *Thing) RemoveAction(actionName, actionID string) bool {
	action := thing.Action(actionName, actionID)
	if action.ID() == "" {
		return false
	}

	defer action.Cancel() // Cancel action after delete from origin.

	actions := thing.actions[actionName]
	for k, ac := range actions {
		if ac != nil && ac.ID() == actionID {
			actions[k] = nil
		}
	}

	return true
}

// AddAvailableAction Add an available action.
//
// @param name     Name of the action
// @param metadata Action metadata, i.e. type, description, etc., as a
//                 JSONObject
// @param action   Instantiate for this action
func (thing *Thing) AddAvailableAction(name string, metadata json.RawMessage, action Actioner) {
	thing.availableActions[name] = NewAvailableAction(metadata, action)
	thing.actions[name] = []*Action{}
}

// AddSubscriber Add a new websocket subscriber.
//
// @param ws The websocket
func (thing *Thing) AddSubscriber(wsID string, ws *websocket.Conn) {
	thing.subscribers[wsID] = ws
}

// RemoveSubscriber Remove a websocket subscriber.
//
// @param ws The websocket
func (thing *Thing) RemoveSubscriber(name string, ws *websocket.Conn) {

	delete(thing.subscribers, name)

	for name := range thing.availableEvents {
		thing.RemoveEventSubscriber(name, ws)
	}
}

// AddEventSubscriber Add a new websocket subscriber to an event.
//
// @param name Name of the event
// @param ws   The websocket
func (thing *Thing) AddEventSubscriber() {}

// RemoveEventSubscriber Remove a websocket subscriber from an event.
//
// @param name Name of the event
// @param ws   The websocket
func (thing *Thing) RemoveEventSubscriber(name string, ws *websocket.Conn) error {

	delete(thing.availableEvents, name)

	if _, ok := thing.availableEvents[name]; ok {
		for _, eventWS := range thing.availableEvents[name].subscribers {
			if err := eventWS.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

type message struct {
	MessageType string          `json:"messageType"`
	Data        json.RawMessage `json:"data"`
}

// PropertyNotify Notify all subscribers of a property change.
//
// @param property The property that changed
func (thing *Thing) PropertyNotify(property Property) error {
	str := message{
		MessageType: "propertyStatus",
		Data:        property.AsPropertyDescription(),
	}
	msg, err := json.Marshal(str)
	if err != nil {
		return err
	}
	for _, sub := range thing.subscribers {
		if err := sub.WriteJSON(msg); err != nil {
			return err
		}
	}
	return nil
}

// ActionNotify Notify all subscribers of an action status change.
//
// @param action The action whose status changed
func (thing *Thing) ActionNotify(action *Action) error {
	str := message{
		MessageType: "actionStatus",
		Data:        action.AsActionDescription(),
	}
	msg, err := json.Marshal(str)
	if err != nil {
		return err
	}
	for _, sub := range thing.subscribers {
		if err := sub.WriteJSON(msg); err != nil {
			return err
		}
	}
	return nil
}

// EventNotify Notify all subscribers of an event.
//
// @param event The event that occurred
func (thing *Thing) EventNotify(event *Event) error {
	eventName := event.Name()
	if _, ok := thing.availableEvents[eventName]; !ok {
		return errors.New("Event not found. ")
	}
	str := message{
		MessageType: "event",
		Data:        event.AsEventDescription(),
	}
	msg, err := json.Marshal(str)
	if err != nil {
		return err
	}
	for _, sub := range thing.subscribers {
		if err := sub.WriteJSON(msg); err != nil {
			return err
		}
	}
	return nil
}

// AvailableEvent Class to describe an event available for subscription.
type AvailableEvent struct {
	metadata    json.RawMessage
	subscribers map[string]*websocket.Conn
}

// NewAvailableEvent Initialize the object.
//
// @param metadata The event metadata
func NewAvailableEvent(metadata json.RawMessage) *AvailableEvent {
	return &AvailableEvent{metadata: metadata, subscribers: make(map[string]*websocket.Conn)}
}

// Metadata Get the event metadata.
//
// @return The metadata.
func (ae *AvailableEvent) Metadata() json.RawMessage {
	return ae.metadata
}

// AvailableAction Class to describe an action available to be taken.
type AvailableAction struct {
	metadata json.RawMessage
	action   *Action
	schema   interface{}
	cls      Actioner
}

// NewAvailableAction Initialize the object.
//
// @param metadata The action metadata
// @param action   Instance for the action
func NewAvailableAction(metadata json.RawMessage, cls Actioner) *AvailableAction {
	ac := &AvailableAction{}
	ac.metadata = metadata
	ac.cls = cls

	// Creating the map for input
	m := map[string]json.RawMessage{}
	json.Unmarshal(metadata, &m)
	if _, ok := m["input"]; ok {
		ac.schema = m["input"]
	}

	return ac
}

// Get the class to instantiate for the action.
//
// @return The class.
func (ac *AvailableAction) getCls() Actioner {
	return ac.cls
}

// Metadata Get the action metadata.
//
// @return The metadata.
func (ac *AvailableAction) Metadata() []byte {
	metaData, _ := json.Marshal(ac.metadata)
	return metaData
}

// Action Get the class to instantiate for the action.
//
// @return The class.
func (ac *AvailableAction) Action() *Action {
	return ac.action
}

// ValidateActionInput Validate the input for a new action.
//
// @param actionInput The input to validate
// @return Boolean indicating validation success.
func (ac *AvailableAction) ValidateActionInput(actionInput interface{}) bool {
	if ac.schema == nil {
		return true
	}
	_, err := json.Marshal(actionInput)

	return err == nil
}
