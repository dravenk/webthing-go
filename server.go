package webthing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
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

		handlerfuncs(prePath, preIdx, thingHandle, propertiesHandle, actionsHandle, eventsHandle)
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

		handlerfuncs(prePath, preIdx, thingHandle, propertiesHandle, actionsHandle, eventsHandle)
	}

	return server
}

func handlerfuncs(prePath, preIdx string,
	thingHandle *ThingHandle,
	propertiesHandle *PropertiesHandle,
	actionsHandle *ActionsHandle,
	eventsHandle *EventsHandle,
) {

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

// // BaseHandler Base handler that is initialized with a list of things.
// type BaseHandler interface {
// 	Get(w http.ResponseWriter, r *http.Request)
// 	Post(w http.ResponseWriter, r *http.Request)
// 	Put(w http.ResponseWriter, r *http.Request)
// 	Delete(w http.ResponseWriter, r *http.Request)
// }

// GetInterface Implementation of http Get menthod.
type GetInterface interface {
	Get(w http.ResponseWriter, r *http.Request)
}

// PostInterface Implementation of http Post menthod.
type PostInterface interface {
	Post(w http.ResponseWriter, r *http.Request)
}

// PutInterface Implementation of http Put menthod.
type PutInterface interface {
	Put(w http.ResponseWriter, r *http.Request)
}

// DeleteInterface Implementation of http Delete menthod.
type DeleteInterface interface {
	Delete(w http.ResponseWriter, r *http.Request)
}

// BaseHandle Base handler that is initialized with a list of things.
// func BaseHandle(h BaseHandler, w http.ResponseWriter, r *http.Request) {
func BaseHandle(h interface{}, w http.ResponseWriter, r *http.Request) {
	corsResponse(w)
	jsonResponse(w)
	switch r.Method {
	case http.MethodGet:
		if base, ok := h.(GetInterface); ok {
			base.Get(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case http.MethodPost:
		if base, ok := h.(PostInterface); ok {
			base.Post(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case http.MethodPut:
		if base, ok := h.(PutInterface); ok {
			base.Put(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case http.MethodDelete:
		if base, ok := h.(DeleteInterface); ok {
			base.Delete(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// ThingsHandle things struct.
type ThingsHandle struct {
	Things   []*Thing
	basePath string
}

// Handle handle request.
func (h *ThingsHandle) Handle(w http.ResponseWriter, r *http.Request) {
	if len(h.Things) == 1 {
		thingHandle := &ThingHandle{h.Things[0]}
		thingHandle.Handle(w, r)
		return
	}

	BaseHandle(h, w, r)
}

// Get Handle a Get request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *ThingsHandle) Get(w http.ResponseWriter, r *http.Request) {

	var things []json.RawMessage
	for _, thing := range h.Things {
		things = append(things, thing.AsThingDescription())
	}
	content, _ := json.Marshal(things)

	if _, err := w.Write(content); err != nil {
		fmt.Println(err)
	}
}

// ThingHandle Handle a request to thing.
type ThingHandle struct {
	*Thing
}

// Handle a request to /thing.
func (h *ThingHandle) Handle(w http.ResponseWriter, r *http.Request) {
	BaseHandle(h, w, r)
}

// Get Handle a Get request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *ThingHandle) Get(w http.ResponseWriter, r *http.Request) {
	base := h.Thing.AsThingDescription()

	var ls map[string][]Link
	json.Unmarshal(base, &ls)

	scheme := "ws"
	// if r.URL.Scheme != "" {
	// 	scheme = r.URL.Scheme
	// }
	wsHref := fmt.Sprintf("%s://%s%s", scheme, r.Host, h.Href())
	ls["links"] = append(ls["links"], Link{
		Rel:  "alternate",
		Href: filepath.Clean(strings.TrimRight(wsHref+"/"+h.Href(), "/")),
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
}

// PropertiesHandle Handle a request to /properties.
type PropertiesHandle struct {
	*ThingHandle
}

// Handle Handle a request to /properties.
func (h *PropertiesHandle) Handle(w http.ResponseWriter, r *http.Request) {
	if name, err := resource(trimSlash(r.RequestURI)); err == nil {
		propertyHandle := &PropertyHandle{h, h.properties[name]}
		propertyHandle.Handle(w, r)
		return
	}

	BaseHandle(h, w, r)
}

// Get Handle a Get request.
//
// @param {Object} r The request object
// @param {Object} w The response object
func (h *PropertiesHandle) Get(w http.ResponseWriter, r *http.Request) {
	content, err := json.Marshal(h.Thing.Properties())
	if err != nil {
		fmt.Println(err)
	}
	if _, err := w.Write(content); err != nil {
		fmt.Println(err)
	}
}
