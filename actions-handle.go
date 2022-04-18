package webthing

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// ActionsHandle Handle a request to /actions.
type ActionsHandle struct {
	*ThingHandle
}

// Handle a request to /actions.
func (h *ActionsHandle) Handle(w http.ResponseWriter, r *http.Request) {
	if name, actionID, err := h.matchActionOrID(trimSlash(r.RequestURI)); err == nil {
		action := h.Thing.Action(name, actionID)
		actionHandle := &ActionHandle{h, name}
		if actionID != "" {
			actionIDHandle := &ActionIDHandle{actionHandle, action}
			actionIDHandle.Handle(w, r)
			return
		}
		actionHandle.Handle(w, r)
		return
	}
	BaseHandle(h, w, r)
}

func (h *ActionsHandle) matchActionOrID(path string) (actionName, actionID string, err error) {
	re := regexp.MustCompile(`^/actions/(.*)/(.*)`)
	name := re.FindStringSubmatch(path)
	if name != nil {
		return name[1], name[2], nil
	}
	m := validPath().FindStringSubmatch(path)
	if m == nil {
		return "", "", errors.New(" Invalid! ")
	}
	return m[2], "", nil // The resource is the second subexpression.
}

// Get Handle a GET request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *ActionsHandle) Get(w http.ResponseWriter, r *http.Request) {
	var description []json.RawMessage

	for name := range h.Thing.actions {
		if content := h.Thing.ActionDescriptions(name); content != nil {
			description = append(description, content...)
		}
	}
	if len(description) == 0 {
		if _, err := w.Write([]byte(`{}`)); err != nil {
			fmt.Println(err)
		}
		return
	}

	content, _ := json.Marshal(description)
	if _, err := w.Write(content); err != nil {
		fmt.Println(err)
	}
}

// Post Handle a POST request.
//
// @param {Object} req The request object
// @param {Object} res The response object
func (h *ActionsHandle) Post(w http.ResponseWriter, r *http.Request) {
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
		if _, ok := h.Thing.actions[name]; ok {
			input := params["input"]
			action, err := h.Thing.PerformAction(name, input)

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
	w.Write(content)
}
