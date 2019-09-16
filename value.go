package webthing

import "fmt"

// Value A property value.
//
// This is used for communicating between the Thing representation and the
// actual physical thing implementation.
//
// Notifies all observers when the underlying value changes through an external
// update (command to turn the light off) or if the underlying sensor reports a
// new value.
type Value struct {
	lastValue      interface{}
	valueForwarder []func(interface{})
}

// NewValue Initialize the object.
//
// @param {*} initialValue The initial value
// @param {function?} valueForwarder The method that updates the actual value
//                                   on the thing
func NewValue(initialValue interface{}, valueForwarder ...func(interface{})) Value {
	return Value{initialValue, valueForwarder}
}

// Set a new value for this thing.
//
// @param {*} value Value to set
func (v *Value) Set(value interface{}) {
	if v.valueForwarder != nil {
		for _, valueForwarder := range v.valueForwarder {
			valueForwarder(value)
		}
	}

	v.NotifyOfExternalUpdate(value)
}

// Get Return the last known value from the underlying thing.
//
// @returns the value.
func (v *Value) Get() interface{} {
	return v.lastValue
}

// NotifyOfExternalUpdate Notify observers of a new value.
//
// @param {*} value New value
func (v *Value) NotifyOfExternalUpdate(value interface{}) {
	if value != nil && value != v.lastValue {
		v.lastValue = value
		fmt.Println("Value update: ", value)
		//v.emit('update', value);
	}
}
