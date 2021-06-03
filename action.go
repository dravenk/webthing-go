package webthing

import (
	"encoding/json"
	"fmt"
)

// Action An Action represents an individual action on a thing.
type Action struct {
	id            string
	thing         *Thing
	name          string
	input         *json.RawMessage
	hrefPrefix    string
	href          string
	status        string
	timeRequested string
	timeCompleted string

	// Override this with the code necessary to perform the action.
	PerformAction func() *Action

	// Override this with the code necessary to cancel the action.
	Cancel func()
}

// Actioner Customize the methods that the action must implement.
type Actioner interface {
	// Custom Action need create a Generator to generate a action.
	// The application will invoke the Action created by the Generator method.
	// This is very similar to simply constructor.
	// See thing.PerformAction()*Action
	Generator(thing *Thing) *Action

	// Override this with the code necessary to perform the action.
	PerformAction() *Action

	// Override this with the code necessary to cancel the action.
	Cancel()

	Start() *Action
	Finish() *Action

	AsActionDescription() []byte
	SetHrefPrefix(prefix string)
	ID() string
	Name() string
	Href() string
	Status() string
	Thing() *Thing
	TimeRequested() string
	TimeCompleted() string
	Input() *json.RawMessage
	SetInput(input *json.RawMessage)
}

// NewAction Initialize the object.
// @param id    ID of this action
// @param thing Thing this action belongs to
// @param name  Name of the action
// @param input Any action inputs
func NewAction(id string, thing *Thing, name string, input *json.RawMessage, PerformAction func() *Action, Cancel func()) *Action {
	action := &Action{
		id:            id,
		thing:         thing,
		name:          name,
		hrefPrefix:    "",
		href:          fmt.Sprintf("/actions/%s/%s", name, id),
		status:        "created",
		timeRequested: Timestamp(),
		PerformAction: PerformAction,
		Cancel:        Cancel,
	}
	if input != nil {
		action.input = input
	}

	return action
}

// AsActionDescription Get the action description.
// @return Description of the action as a JSONObject.
func (action *Action) AsActionDescription() []byte {
	actionName := action.Name()
	obj := make(map[string]interface{})
	actionObj := make(map[string]interface{})

	if input := action.Input(); input != nil {
		actionObj["input"] = input
	}
	if timeCompleted := action.TimeCompleted(); timeCompleted != "" {
		actionObj["timeCompleted"] = timeCompleted
	}
	actionObj["href"] = action.Href()
	actionObj["status"] = action.Status()
	actionObj["timeRequested"] = action.TimeRequested()
	obj[actionName] = actionObj

	description, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return []byte(err.Error())
	}

	return description
}

// SetHrefPrefix Set the prefix of any hrefs associated with this action.
// @param prefix The prefix
func (action *Action) SetHrefPrefix(prefix string) {
	action.hrefPrefix = prefix
}

// ID Get this action's ID.
// @returns {String} The ID.
func (action *Action) ID() string {
	return action.id
}

// Name Get this action's name.
// @returns {String} The name.
func (action *Action) Name() string {
	return action.name
}

// Href Get this action's href.
// @returns {String} The href.
func (action *Action) Href() string {
	return action.hrefPrefix + action.href
}

// Status Get this action's status.
// Get this action's status.
// @returns {String} The status.
func (action *Action) Status() string {
	return action.status
}

// Thing Get the thing associated with this action.
// @returns {Object} The thing.
func (action *Action) Thing() *Thing {
	return action.thing
}

// TimeRequested Get the time the action was requested.
// @returns {String} The time.
func (action *Action) TimeRequested() string {
	return action.timeRequested
}

// TimeCompleted Get the time the action was completed.
// @returns {String} The time.
func (action *Action) TimeCompleted() string {
	return action.timeCompleted
}

// Input Get the inputs for this action.
// @returns {Object} The inputs.
func (action *Action) Input() *json.RawMessage {
	return action.input
}

// SetInput Set any input to this action.
// @param input The input
func (action *Action) SetInput(input *json.RawMessage) {
	if input != nil {
		action.input = input
	}
}

// Start performing the action.
func (action *Action) Start() *Action {
	// Todo
	// Handle an error performing the action.
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Perform Action encountered an error")
		}
	}()

	action.status = "pending"
	action.thing.ActionNotify(action)
	action.PerformAction()
	action.Finish()

	return action
}

// Finish performing the action.
func (action *Action) Finish() *Action {
	action.status = "completed"
	action.timeCompleted = Timestamp()
	action.thing.ActionNotify(action)
	return action
}
