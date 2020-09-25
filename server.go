package webthing

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ThingServer Web Thing Server.
type ThingServer struct {
	*http.Server
	Things   []*Thing
	Name     string
	BasePath string
}

// NewWebThingServer Initialize the WebThingServer.
//
// @param thingType        List of Things managed by this server
// @param basePath         Base URL path to use, rather than '/'
//
func NewWebThingServer(thingType ThingsType, httpServer *http.Server, basePath string) *ThingServer {
	server := &ThingServer{httpServer, thingType.Things(), thingType.Name(), basePath}
	thingsNum := len(server.Things)

	thingsHandle := &ThingsHandle{server.Things, basePath}
	http.HandleFunc("/", thingsHandle.Handle)

	if thingsNum == 1 {
		thing := server.Things[0]
		prePath := strings.TrimRight(server.BasePath+"/"+thing.Title(), "/")
		preIdx := strings.TrimRight(server.BasePath, "/")
		thing.SetHrefPrefix(preIdx)
		thingHandle := &ThingHandle{thing}
		propertiesHandle := &PropertiesHandle{thingHandle}
		actionsHandle := &ActionsHandle{thingHandle}
		eventsHandle := &EventsHandle{thingHandle}

		http.HandleFunc(prePath, thingHandle.Handle)
		http.HandleFunc(prePath+"/properties", propertiesHandle.Handle)
		http.HandleFunc(prePath+"/properties/", propertiesHandle.Handle)
		http.HandleFunc(prePath+"/actions", actionsHandle.Handle)
		http.HandleFunc(prePath+"/actions/", actionsHandle.Handle)
		http.HandleFunc(prePath+"/events", eventsHandle.Handle)
		http.HandleFunc(prePath+"/events/", eventsHandle.Handle)

		http.HandleFunc(preIdx+"/properties", propertiesHandle.Handle)
		http.HandleFunc(preIdx+"/properties/", propertiesHandle.Handle)
		http.HandleFunc(preIdx+"/actions", actionsHandle.Handle)
		http.HandleFunc(preIdx+"/actions/", actionsHandle.Handle)
		http.HandleFunc(preIdx+"/events", eventsHandle.Handle)
		http.HandleFunc(preIdx+"/events/", eventsHandle.Handle)

		if preIdx != "" {
			http.HandleFunc(preIdx, thingHandle.Handle)
		}
		return server
	}

	for id, thing := range server.Things {
		prePath := strings.TrimRight(server.BasePath+"/"+thing.Title(), "/")
		preIdx := "/" + strconv.Itoa(id)
		thing.SetHrefPrefix(preIdx)
		thingHandle := &ThingHandle{thing}
		propertiesHandle := &PropertiesHandle{thingHandle}
		actionsHandle := &ActionsHandle{thingHandle}
		eventsHandle := &EventsHandle{thingHandle}

		http.HandleFunc(prePath, thingHandle.Handle)
		http.HandleFunc(prePath+"/properties", propertiesHandle.Handle)
		http.HandleFunc(prePath+"/properties/", propertiesHandle.Handle)
		http.HandleFunc(prePath+"/actions", actionsHandle.Handle)
		http.HandleFunc(prePath+"/actions/", actionsHandle.Handle)
		http.HandleFunc(prePath+"/events", eventsHandle.Handle)
		http.HandleFunc(prePath+"/events/", eventsHandle.Handle)

		http.HandleFunc(preIdx+"/properties", propertiesHandle.Handle)
		http.HandleFunc(preIdx+"/properties/", propertiesHandle.Handle)
		http.HandleFunc(preIdx+"/actions", actionsHandle.Handle)
		http.HandleFunc(preIdx+"/actions/", actionsHandle.Handle)
		http.HandleFunc(preIdx+"/events", eventsHandle.Handle)
		http.HandleFunc(preIdx+"/events/", eventsHandle.Handle)

		http.HandleFunc(preIdx, thingHandle.Handle)
	}

	return server
}

// Start Start listening for incoming connections.
//
// @return Error on failure to listen on port
func (server *ThingServer) Start() error {
	return server.ListenAndServe()
}

// Stop Stop listening.
func (server *ThingServer) Stop() error {
	return server.Close()
}

// corsResponse Add necessary CORS headers to response.
//
// @param response Response to add headers to
// @return The Response object.
func corsResponse(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, PUT, POST, DELETE")
	return w
}

//jsonResponse Add json headers to response.
func jsonResponse(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	return w
}

// ThingsType Container of Things Type
type ThingsType interface {

	// Thing Get the thing at the given index.
	//
	// @param idx Index of thing.
	// @return The thing, or null.
	Thing(idx int) *Thing

	// Things Get the list of things.
	//
	// @return The list of things.
	Things() []*Thing

	// Name Get the mDNS server name.
	//
	// @return The server name.
	Name() string
}

// SingleThing A container for a single thing.
type SingleThing struct {
	thing *Thing
}

// NewSingleThing Initialize the container.
//
// @param {Object} thing The thing to store
func NewSingleThing(thing *Thing) *SingleThing {
	return &SingleThing{thing}
}

// Thing Get the thing at the given index.
func (st *SingleThing) Thing(idx int) *Thing {
	return st.thing
}

// Things Get the list of things.
func (st *SingleThing) Things() []*Thing {
	return []*Thing{st.thing}
}

// Name Get the mDNS server name.
func (st *SingleThing) Name() string {
	return st.thing.title
}

// MultipleThings  A container for multiple things.
type MultipleThings struct {
	things []*Thing
	name   string
}

// NewMultipleThings Initialize the container.
//
// @param {Object} things The things to store
// @param {String} name The mDNS server name
func NewMultipleThings(things []*Thing, name string) *MultipleThings {
	mt := &MultipleThings{
		things: things,
		name:   name,
	}
	return mt
}

// Thing Get the thing at the given index.
//
// @param {Number|String} idx The index
func (mt *MultipleThings) Thing(idx int) *Thing {
	return mt.things[idx]
}

// Things Get the list of things.
func (mt *MultipleThings) Things() []*Thing {
	return mt.things
}

// Name Get the mDNS server name.
func (mt *MultipleThings) Name() string {
	return mt.name
}

// BaseHandler Base handler that is initialized with a list of things.
type BaseHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// BaseHandle Base handler that is initialized with a list of things.
func BaseHandle(h BaseHandler, w http.ResponseWriter, r *http.Request) {
	corsResponse(w)
	jsonResponse(w)
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPost:
		h.Post(w, r)
	case http.MethodPut:
		h.Put(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	}
}

type ThingsHandle struct {
	Things   []*Thing
	basePath string
}

func (h *ThingsHandle) Handle(w http.ResponseWriter, r *http.Request) {
	if len(h.Things) == 1 {
		thingHandle := &ThingHandle{h.Things[0]}
		thingHandle.Handle(w, r)
		return
	}

	BaseHandle(h, w, r)
}

func (h *ThingsHandle) Get(w http.ResponseWriter, r *http.Request) {

	var things []json.RawMessage
	for _, thing := range h.Things {
		things = append(things, thing.AsThingDescription())
	}
	content, _ := json.Marshal(things)

	if _, err := w.Write(content); err != nil {
		fmt.Println(err)
	}
	return
}

func (h *ThingsHandle) Post(w http.ResponseWriter, r *http.Request)   {}
func (h *ThingsHandle) Put(w http.ResponseWriter, r *http.Request)    {}
func (h *ThingsHandle) Delete(w http.ResponseWriter, r *http.Request) {}

type ThingHandle struct {
	*Thing
}

func (h *ThingHandle) Handle(w http.ResponseWriter, r *http.Request) {
	BaseHandle(h, w, r)
}

func (h *ThingHandle) Get(w http.ResponseWriter, r *http.Request) {
	base := h.Thing.AsThingDescription()

	var ls map[string][]Link
	json.Unmarshal(base, &ls)

	scheme := "ws"
	// if r.URL.Scheme != "" {
	// 	scheme = r.URL.Scheme
	// }
	wsHref := fmt.Sprintf("%s://%s%s",scheme, r.Host, h.Href());
	ls["links"] = append(ls["links"], Link{
		Rel:       "alternate",
		Href: strings.TrimRight(wsHref+"/"+h.Href(), "/"),
	})
	var desc map[string]interface{}
	if err := json.Unmarshal(base, &desc); err != nil {
		fmt.Print(err)
	}
	desc["links"] = ls["links"]

	type securityDefinitions struct {
		NosecSc struct {
			Scheme string `json:"scheme"`
		} `json:"nosec_sc"`
	}
	sec := &securityDefinitions{}
	sec.NosecSc.Scheme = "nosec"
	desc["securityDefinitions"] = sec
	desc["security"] = "nosec_sc"

	re, _ := json.Marshal(desc)
	if _, err := w.Write(re); err != nil {
		fmt.Println(err)
	}

	return
}

func (h *ThingHandle) Post(w http.ResponseWriter, r *http.Request)   {}
func (h *ThingHandle) Put(w http.ResponseWriter, r *http.Request)    {}
func (h *ThingHandle) Delete(w http.ResponseWriter, r *http.Request) {}

/**
 * Handle a request to /properties.
 */
type PropertiesHandle struct {
	*ThingHandle
}

func (h *PropertiesHandle) Handle(w http.ResponseWriter, r *http.Request) {
	if name, err := resource(trimSlash(r.RequestURI)); err == nil {
		propertyHandle := &PropertyHandle{h, h.properties[name]}
		propertyHandle.Handle(w, r)
		return
	}

	BaseHandle(h, w, r)
}

func trimSlash(path string) string {
	l := len(path)
	if l != 1 && path[l-1:] == "/" {
		return path[:l-1]
	}
	return path
}
func resource(path string) (string, error) {
	m := validPath().FindStringSubmatch(path)
	if m == nil {
		return "", errors.New("Invalid! ")
	}
	return m[2], nil // The resource is the second subexpression.
}

func validPath() *regexp.Regexp {
	return regexp.MustCompile("^/(properties|actions|events)/([a-zA-Z0-9]+)$")
}

func (h *PropertiesHandle) Get(w http.ResponseWriter, r *http.Request) {
	content, err := json.Marshal(h.Thing.Properties())
	if err != nil {
		fmt.Println(err)
	}
	if _, err := w.Write(content); err != nil {
		fmt.Println(err)
	}
}

func (h *PropertiesHandle) Post(w http.ResponseWriter, r *http.Request)   {}
func (h *PropertiesHandle) Put(w http.ResponseWriter, r *http.Request)    {}
func (h *PropertiesHandle) Delete(w http.ResponseWriter, r *http.Request) {}

/**
 * Handle a request to /properties/<property>.
 */
type PropertyHandle struct {
	*PropertiesHandle
	*Property
}

func (h *PropertyHandle) Handle(w http.ResponseWriter, r *http.Request) {
	name, err := resource(trimSlash(r.RequestURI))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.Property = h.properties[name]
	BaseHandle(h, w, r)
}

// Handle a GET request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *PropertyHandle) Get(w http.ResponseWriter, r *http.Request) {

	name := h.Property.Name()
	value := h.Property.Value().Get()
	description := make(map[string]interface{})
	description[name] = value

	content, err := json.Marshal(description)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := w.Write(content); err != nil {
		fmt.Println(err)
	}
}

// Handle a PUT request.
//
// @param {Object} R The request object
// @param {Object} W The response object
func (h *PropertyHandle) Put(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var obj map[string]interface{}
	err := json.Unmarshal(body, &obj)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := h.Property.Name()
	h.Property.Value().Set(obj[name])

	description := make(map[string]interface{})
	description[name] = h.Property.Value().Get()
	content, err := json.Marshal(description)

	if _, err = w.Write(content); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
func (h *PropertyHandle) Delete(w http.ResponseWriter, r *http.Request) {}

// Handle a request to /actions.
type ActionsHandle struct {
	*ThingHandle
}

func (h *ActionsHandle) Handle(w http.ResponseWriter, r *http.Request) {
	if name, actionID, err := h.MatchActionOrID(trimSlash(r.RequestURI)); err == nil {
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

func (h *ActionsHandle) MatchActionOrID(path string) (actionName, actionID string, err error) {
	re := regexp.MustCompile(`^/actions/(.*)/(.*)`)
	name := re.FindStringSubmatch(path)
	if name != nil {
		return name[1], name[2], nil
	}
	m := validPath().FindStringSubmatch(path)
	if m == nil {
		return "", "", errors.New("Invalid! ")
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
	return
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
			action := h.Thing.PerformAction(name, input)

			if action == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Perform an Action in a goroutine.
			go action.Start()

			if len(obj) == 1 {
				w.WriteHeader(http.StatusCreated)
				_, err = w.Write(action.AsActionDescription())
				return
			}
			description = append(description, action.AsActionDescription())
		}
	}
	content, _ := json.Marshal(description)
	_, err = w.Write(content)
	return
}

func (h *ActionsHandle) Put(w http.ResponseWriter, r *http.Request)    {}
func (h *ActionsHandle) Delete(w http.ResponseWriter, r *http.Request) {}

// ActionHandle Handle a request to /actions/<action_name>.
type ActionHandle struct {
	*ActionsHandle
	ActionName string
}

func (h *ActionHandle) Handle(w http.ResponseWriter, r *http.Request) {
	BaseHandle(h, w, r)
}

func (h *ActionHandle) Get(w http.ResponseWriter, r *http.Request) {
	if descriptions := h.Thing.ActionDescriptions(h.ActionName); len(descriptions) > 0 {
		content, _ := json.Marshal(descriptions)
		if _, err := w.Write(content); err != nil {
			fmt.Println(err)
		}
	}
	return
}

func (h *ActionHandle) Post(w http.ResponseWriter, r *http.Request) {
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
		if name == h.ActionName {
			if _, ok := h.Thing.actions[name]; ok {
				input := params["input"]
				action := h.Thing.PerformAction(name, input)

				if action == nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				// Perform an Action in a goroutine.
				go action.Start()

				if len(obj) == 1 {
					w.WriteHeader(http.StatusCreated)
					_, err = w.Write(action.AsActionDescription())
					return
				}

				description = append(description, action.AsActionDescription())
			}
		}
	}
	content, _ := json.Marshal(description)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(content)
	return

}

func (h *ActionHandle) Put(w http.ResponseWriter, r *http.Request)    {}
func (h *ActionHandle) Delete(w http.ResponseWriter, r *http.Request) {}

// ActionIDHandle Handle a request to /actions/<action_name>/<action_id>.
type ActionIDHandle struct {
	*ActionHandle
	*Action
}

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

func (h *ActionIDHandle) Get(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write(h.Action.AsActionDescription()); err != nil {
		fmt.Println(err)
	}
	return

}

func (h *ActionIDHandle) Post(w http.ResponseWriter, r *http.Request) {}
func (h *ActionIDHandle) Put(w http.ResponseWriter, r *http.Request)  {}
func (h *ActionIDHandle) Delete(w http.ResponseWriter, r *http.Request) {
	if h.RemoveAction(h.ActionName, h.Action.ID()) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// Handle a request to /actions.
type EventsHandle struct {
	*ThingHandle
}

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
	return
}
func (h *EventsHandle) Post(w http.ResponseWriter, r *http.Request)   {}
func (h *EventsHandle) Put(w http.ResponseWriter, r *http.Request)    {}
func (h *EventsHandle) Delete(w http.ResponseWriter, r *http.Request) {}

// Handle a request to /actions.
type EventHandle struct {
	*EventsHandle
	eventName string
}

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
	return
}

func (h *EventHandle) Post(w http.ResponseWriter, r *http.Request)   {}
func (h *EventHandle) Put(w http.ResponseWriter, r *http.Request)    {}
func (h *EventHandle) Delete(w http.ResponseWriter, r *http.Request) {}
