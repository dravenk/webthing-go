package webthing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// PropertyHandle a request to /properties/<property>.
type PropertyHandle struct {
	*PropertiesHandle
	*Property
}

// Handle a request to /properties/<property>.
func (h *PropertyHandle) Handle(w http.ResponseWriter, r *http.Request) {
	name, err := resource(trimSlash(r.RequestURI))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.Property = h.properties[name]
	BaseHandle(h, w, r)
}

// Get Handle a GET request.
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

// Put Handle a PUT request.
//
// @param {Object} r The request object
// @param {Object} w The response object
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
	content, _ := json.Marshal(description)

	if _, err = w.Write(content); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
