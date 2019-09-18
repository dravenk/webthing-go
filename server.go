package webthing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

// ThingServer Web Thing Server.
type ThingServer struct {
	*http.Server
	Things []*Thing
	Name   string
	//BasePath string
}

// NewWebThingServer Initialize the WebThingServer.
//
// @param thingType        List of Things managed by this server
// @param basePath         Base URL path to use, rather than '/'
//
func NewWebThingServer(thingType ThingsType, httpServer *http.Server) *ThingServer {
	server := &ThingServer{httpServer, thingType.Things(), thingType.Name()}

	thingsHandle := &ThingsHandle{server.Things}

	http.HandleFunc("/", thingsHandle.Handle)

	basePath := "/things"

	http.HandleFunc(basePath, thingsHandle.Handle)

	for id, thing := range server.Things {
		thingName := thing.Title()
		thingIdx := strconv.Itoa(id)

		thing.SetHrefPrefix(fmt.Sprintf("%s/%s", basePath, thingIdx))

		thingHandle := &ThingHandle{thing}
		http.HandleFunc(basePath+"/"+thingIdx, thingHandle.Handle)
		http.HandleFunc(basePath+"/"+thingName, thingHandle.Handle)

		propertiesHandle := &PropertiesHandle{thingHandle}
		http.HandleFunc(basePath+"/"+thingIdx+"/properties", propertiesHandle.Handle)
		http.HandleFunc(basePath+"/"+thingName+"/properties", propertiesHandle.Handle)

		actionsHandle := &ActionsHandle{thingHandle}
		http.HandleFunc(basePath+"/"+thingIdx+"/actions", actionsHandle.Handle)
		http.HandleFunc(basePath+"/"+thingName+"/actions", actionsHandle.Handle)

		eventsHandle := &EventsHandle{thingHandle}
		http.HandleFunc(basePath+"/"+thingIdx+"/events", eventsHandle.Handle)
		http.HandleFunc(basePath+"/"+thingName+"/events", eventsHandle.Handle)

		for name, property := range thing.properties {
			propertyHandle := PropertyHandle{propertiesHandle, property}
			http.HandleFunc(basePath+"/"+thingIdx+"/properties/"+name, propertyHandle.Handle)
			http.HandleFunc(basePath+"/"+thingName+"/properties/"+name, propertyHandle.Handle)
		}

		for actionName := range thing.actions {
			actionHandle := &ActionHandle{actionsHandle, actionName}
			http.HandleFunc(basePath+"/"+thingIdx+"/actions/"+actionName, actionHandle.Handle)
			http.HandleFunc(basePath+"/"+thingName+"/actions/"+actionName, actionHandle.Handle)

			actionIDHandle := &ActionIDHandle{actionHandle}
			http.HandleFunc(basePath+"/"+thingIdx+"/actions/"+actionName+"/", actionIDHandle.Handle)
			http.HandleFunc(basePath+"/"+thingName+"/actions/"+actionName+"/", actionIDHandle.Handle)

		}

		eventHandle := &EventHandle{eventsHandle}
		http.HandleFunc(basePath+"/"+thingIdx+"/events/", eventHandle.Handle)
		http.HandleFunc(basePath+"/"+thingName+"/events/", eventHandle.Handle)

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
	Things []*Thing
}

func (h *ThingsHandle) Handle(w http.ResponseWriter, r *http.Request) {
	BaseHandle(h, w, r)
}

func (h *ThingsHandle) Get(w http.ResponseWriter, r *http.Request) {

	if len(h.Things) == 1 {
		thingHandle := &ThingHandle{h.Things[0]}
		thingHandle.Handle(w, r)
		return
	}

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
	content := h.Thing.AsThingDescription()
	if _, err := w.Write(content); err != nil {
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
	BaseHandle(h, w, r)
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
		w.WriteHeader(400)
		fmt.Println(err)
		return
	}

	name := h.Property.Name()
	h.Property.Value().Set(obj[name])

	description := make(map[string]interface{})
	description[name] = h.Property.Value().Get()
	content, err := json.Marshal(description)

	_, err = w.Write(content)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
	}
}
func (h *PropertyHandle) Delete(w http.ResponseWriter, r *http.Request) {}

// Handle a request to /actions.
type ActionsHandle struct {
	*ThingHandle
}

func (h *ActionsHandle) Handle(w http.ResponseWriter, r *http.Request) {
	BaseHandle(h, w, r)
}

// Get Handle a GET request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *ActionsHandle) Get(w http.ResponseWriter, r *http.Request) {
	var description []json.RawMessage

	for name := range h.Thing.actions {
		description = append(description, h.Thing.ActionDescriptions(name)...)
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
		w.WriteHeader(400)
		fmt.Println(err)
		return
	}

	for name, params := range obj {
		if _, ok := h.Thing.actions[name]; ok {
			input := params["input"]
			action := h.Thing.PerformAction(name, input)
			// Perform an Action in a goroutine.
			go action.Start()

			_, err = w.Write(action.AsActionDescription())
			if err != nil {
				fmt.Println(err)
			}
		}
	}
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
	content, _ := json.Marshal(h.Thing.ActionDescriptions(h.ActionName))
	if _, err := w.Write(content); err != nil {
		fmt.Println(err)
	}
}
func (h *ActionHandle) Post(w http.ResponseWriter, r *http.Request)   {}
func (h *ActionHandle) Put(w http.ResponseWriter, r *http.Request)    {}
func (h *ActionHandle) Delete(w http.ResponseWriter, r *http.Request) {}

// ActionIDHandle Handle a request to /actions/<action_name>/<action_id>.
type ActionIDHandle struct {
	*ActionHandle
}

func (h *ActionIDHandle) Handle(w http.ResponseWriter, r *http.Request) {
	BaseHandle(h, w, r)
}

func (h *ActionIDHandle) MatchActionID(path string) (actionID string, bool bool) {
	re := regexp.MustCompile(`/things/(.*)/actions/(.*)` + `/(.*)`)
	if re.MatchString(path) {
		name := re.FindStringSubmatch(path)
		if len(path) >= 4 {
			return name[3], true
		}
	}
	return
}

func (h *ActionIDHandle) Get(w http.ResponseWriter, r *http.Request) {
	if actionID, ok := h.MatchActionID(r.RequestURI); ok {
		action := h.Thing.Action(h.ActionName, actionID)
		if action != nil {
			if _, err := w.Write(action.AsActionDescription()); err != nil {
				fmt.Println(err)
			}
			return
		}
	}
	w.WriteHeader(400)
}

func (h *ActionIDHandle) Post(w http.ResponseWriter, r *http.Request) {}
func (h *ActionIDHandle) Put(w http.ResponseWriter, r *http.Request)  {}
func (h *ActionIDHandle) Delete(w http.ResponseWriter, r *http.Request) {
	if actionID, ok := h.MatchActionID(r.RequestURI); ok {
		if h.Thing.RemoveAction(h.ActionName, actionID) {
			w.Write([]byte(`204 No Content`))
			w.WriteHeader(204)
			return
		}
	}
}

// Handle a request to /actions.
type EventsHandle struct {
	*ThingHandle
}

func (h *EventsHandle) Handle(w http.ResponseWriter, r *http.Request) {
	BaseHandle(h, w, r)
}

// Get Handle a GET request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *EventsHandle) Get(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(h.Thing.EventDescriptions(""))); err != nil {
		fmt.Println(err)
	}
}
func (h *EventsHandle) Post(w http.ResponseWriter, r *http.Request)   {}
func (h *EventsHandle) Put(w http.ResponseWriter, r *http.Request)    {}
func (h *EventsHandle) Delete(w http.ResponseWriter, r *http.Request) {}

// Handle a request to /actions.
type EventHandle struct {
	*EventsHandle
}

func (h *EventHandle) Handle(w http.ResponseWriter, r *http.Request) {
	BaseHandle(h, w, r)
}

// Get Handle a GET request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *EventHandle) Get(w http.ResponseWriter, r *http.Request) {
	if eventName, ok := h.MatchEventName(r.RequestURI); ok {
		w.Write([]byte(h.Thing.EventDescriptions(eventName)))
	}
}

func (h *EventHandle) MatchEventName(path string) (eventName string, bool bool) {
	re := regexp.MustCompile(`/things/(.*)/events/(.*)`)
	if re.MatchString(path) {
		name := re.FindStringSubmatch(path)
		if len(path) >= 3 {
			return name[2], true
		}
	}
	return
}

func (h *EventHandle) Post(w http.ResponseWriter, r *http.Request)   {}
func (h *EventHandle) Put(w http.ResponseWriter, r *http.Request)    {}
func (h *EventHandle) Delete(w http.ResponseWriter, r *http.Request) {}
