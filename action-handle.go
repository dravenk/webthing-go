package webthing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ActionHandle Handle a request to /actions/<action_name>.
type ActionHandle struct {
	*ActionsHandle
	ActionName string
}

// Handle a request to /actions/<action_name>.
func (h *ActionHandle) Handle(w http.ResponseWriter, r *http.Request) {
	BaseHandle(h, w, r)
}

// Get Handle a GET request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *ActionHandle) Get(w http.ResponseWriter, r *http.Request) {
	if descriptions := h.Thing.ActionDescriptions(h.ActionName); len(descriptions) > 0 {
		content, _ := json.Marshal(descriptions)
		if _, err := w.Write(content); err != nil {
			fmt.Println(err)
		}
	}
}

// Post Handle a Post request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *ActionHandle) Post(w http.ResponseWriter, r *http.Request) {
	handleActionPost(h.Thing, w, r)
}

func handleActionPost(th *Thing, w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var obj map[string]map[string]*json.RawMessage
	err := json.Unmarshal(body, &obj)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var description []json.RawMessage
	for name, params := range obj {
		if _, ok := th.actions[name]; ok {
			input := params["input"]
			action, err := th.PerformAction(name, input)

			if action == nil || err != nil {
				fmt.Println("PerformAction failed! The action name is: ", name)
				w.WriteHeader(http.StatusBadRequest)
				// w.Write([]byte(err.Error()))
				return
			}

			// Perform an Action in a goroutine.
			go action.Start()

			if len(obj) == 1 {
				w.WriteHeader(http.StatusCreated)
				w.Write(action.AsActionDescription())
				return
			}

			description = append(description, action.AsActionDescription())
		}
	}

	content, _ := json.Marshal(description)
	w.WriteHeader(http.StatusCreated)
	w.Write(content)
}

// ActionIDHandle Handle a request to /actions/<action_name>/<action_id>.
type ActionIDHandle struct {
	*ActionHandle
	*Action
}

// Handle a request to /actions/<action_name>/<action_id>.
func (h *ActionIDHandle) Handle(w http.ResponseWriter, r *http.Request) {
	if h.Action == nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`Bad request. Action not found.`)); err != nil {
			fmt.Print(err)
		}
		return
	}
	BaseHandle(h, w, r)
}

// Get Handle a GET request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *ActionIDHandle) Get(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write(h.Action.AsActionDescription()); err != nil {
		fmt.Println(err)
	}
}

// Delete Handle a Delete request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *ActionIDHandle) Delete(w http.ResponseWriter, r *http.Request) {
	if h.RemoveAction(h.ActionName, h.Action.ID()) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
