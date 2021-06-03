package webthing

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/types"
)

// Property Initialize the object.
//
// @param thing    Thing this property belongs to
// @param name     Name of the property
// @param value    Value object to hold the property value
// @param metadata Property metadata, i.e. type, description, unit, etc., as
//                 a Map
type Property struct {
	thing      *Thing
	name       string
	value      *Value
	hrefPrefix string
	href       string
	metadata   json.RawMessage
}

// PropertyObject A property object describes an attribute of a Thing and is indexed by a property id.
// See https://iot.mozilla.org/wot/#property-object
type PropertyObject struct {
	AtType      string      `json:"@type,omitempty"`
	Title       string      `json:"title,omitempty"`
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Unit        string      `json:"unit,omitempty"`
	ReadOnly    bool        `json:"readOnly,omitempty"`
	Minimum     json.Number `json:"minimum,omitempty"`
	Maximum     json.Number `json:"maximum,omitempty"`
	Links       []Link      `json:"links,omitempty"`
}

// NewProperty Initialize the object.
//
// @param {Object} thing Thing this property belongs to
// @param {String} name Name of the property
// @param {Value} value Value object to hold the property value
// @param {Object} metadata Property metadata, i.e. type, description, unit,
//                          etc., as an object.
func NewProperty(thing *Thing, name string, value Value, metadata json.RawMessage) *Property {
	property := &Property{
		thing:      thing,
		name:       name,
		value:      &value,
		hrefPrefix: "",
		href:       `/properties/` + name,
		metadata:   metadata,
	}

	// Add the property change observer to notify the Thing about a property
	// change.
	//property.Value.on("update", () => property.thing.PropertyNotify(*property));
	thing.PropertyNotify(*property)

	return property
}

//
// ValidateValue Validate new property value before setting it.
//
// @param {*} value - New value
func (property *Property) ValidateValue(value interface{}) error {
	prop := &PropertyObject{}
	meta, err := property.Metadata().MarshalJSON()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(meta, &prop); err != nil {
		return err
	}
	if prop.ReadOnly {
		return errors.New(" Read-only property. ")
	}
	if !validate(prop.Type, value) {
		return errors.New(" Invalid property value. ")
	}

	return nil
}

// AsPropertyDescription Get the property description.
//
// @returns {Object} Description of the property as an object.
func (property *Property) AsPropertyDescription() []byte {

	link := Link{
		Rel:  "property",
		Href: property.hrefPrefix + property.href,
	}
	base := &PropertyObject{Links: []Link{link}}

	meta, _ := property.Metadata().MarshalJSON()
	if err := json.Unmarshal(meta, base); err != nil {
		fmt.Println("Error: ", err.Error())
		return []byte(err.Error())
	}

	propertyBase, _ := json.Marshal(base)

	return propertyBase
}

// SetHrefPrefix Set the prefix of any hrefs associated with this property.
//
// @param {String} prefix The prefix
func (property *Property) SetHrefPrefix(prefix string) {
	property.hrefPrefix = prefix
}

// Href Get the href of this property.
//
// @returns {String} The href
func (property *Property) Href() string {
	return property.hrefPrefix + property.href
}

// Value Get the current property value.
//
// @returns {*} The current value
func (property *Property) Value() *Value {
	return property.value
}

// SetValue Set the current value of the property.
//
// @param {*} value The value to set
func (property *Property) SetValue(value *Value) error {
	if err := property.ValidateValue(value); err != nil {
		fmt.Print("Ser property value failure. Err: ", err)
		return err
	}
	property.value = value
	return nil
}

// Name Get the name of this property.
//
// @returns {String} The property name.
func (property *Property) Name() string {
	return property.name
}

// Thing Get the thing associated with this property.
//
// @returns {Object} The thing.
func (property *Property) Thing() *Thing {
	return property.thing
}

// Metadata Get the metadata associated with this property
//
// @returns {Object} The metadata
//
func (property *Property) Metadata() json.RawMessage {
	return property.metadata
}

// validate Custom a simply validate. It is not yet possible to verify all types.
//
// A primitive type (one of null, boolean, object, array, number, integer or string as per [json-schema])
// See:  deihttps://iot.mozilla.org/wot/#property-object
//
// See: https://tools.ietf.org/html/draft-zyp-json-schema-04#section-3.5
//
// 3.5.  JSON Schema primitive types
// 	JSON Schema defines seven primitive types for JSON values:
// 	array  A JSON array.
// 	boolean  A JSON boolean.
// 	integer  A JSON number without a fraction or exponent part.
// 	number  Any JSON number.  Number includes integer.
// 	null  The JSON null value.
// 	object  A JSON object.
// 	string  A JSON string.
//
func validate(primitive string, v interface{}) bool {
	switch v.(type) {
	case types.Array:
		if primitive == "array" {
			return true
		}
	case bool:
		if primitive == "boolean" {
			return true
		}
	case string:
		if primitive == "string" {
			return true
		}
	case int, int8, int16, int32, int64:
		if primitive == "integer" || primitive == "number" {
			return true
		}
	case float32, float64:
		if primitive == "number" {
			return true
		}
	}

	return false
}
