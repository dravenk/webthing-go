package main

import (
	"encoding/json"
	"fmt"
	"github.com/dravenk/webthing-go"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	thing := MakeThing()

	singleThing := webthing.NewSingleThing(thing)
	httpServer := &http.Server{Addr: "0.0.0.0:8888"}

	server := webthing.NewWebThingServer(singleThing, httpServer, "")
	log.Fatal(server.Start())
}

func MakeThing() *webthing.Thing {
	// Create a Lamp.
	thing := webthing.NewThing("urn:dev:ops:my-thing-1234",
		"My Lamp",
		[]string{"OnOffSwitch", "Light"},
		"A web connected lamp")

	// Adding an OnOffProperty to thing.
	onDescription := []byte(`{
    "@type": "OnOffProperty",
    "type": "boolean",
    "title": "On/Off",
    "description": "Whether the lamp is turned on"
	}`)
	onValue := webthing.NewValue(true, onValueForwarder)
	on := webthing.NewProperty(thing,
		"on",
		onValue,
		onDescription)
	thing.AddProperty(on)

	// Adding an BrightnessProperty to this Lamp.
	brightnessDescription := []byte(`{
    "@type": "BrightnessProperty",
    "type": "integer",
    "title": "Brightness",
    "description": "The level of light from 0-100",
    "minimum": 0,
    "maximum": 100,
	"unit": "percent"
	}`)

	thing.AddProperty(webthing.NewProperty(thing,
		"brightness",
		webthing.NewValue(50),
		brightnessDescription))

	//Adding a Fade action to this Lamp.
	fadeMeta := []byte(`{
    "title": "Fade",
    "description": "Fade the lamp to a given level",
    "input": {
        "@type": "FadeAction",
        "type": "object",
        "properties": {
            "brightness": {
                "type": "integer",
                "minimum": 0,
                "maximum": 100,
				"unit": "percent"
            },
            "duration": {
                "type": "integer",
                "minimum": 1,
                "unit": "milliseconds"
            }
        }
    }
	}`)
	fade := &FadeAction{}
	thing.AddAvailableAction("fade", fadeMeta, fade)

	thing.AddAvailableEvent("overheated",
		[]byte(`{
            "description":
            "The lamp has exceeded its safe operating temperature",
            "type": "number",
            "unit": "degree celsius"
        }`))

	toggleMeta := []byte(`{
	"title": "Toggle",
	"description": "Toggles a boolean state on and off."
	}`)
	toggle := &ToggleAction{}
	thing.AddAvailableAction("toggle", toggleMeta, toggle)

	return thing
}

// Fade the lamp to a given brightness
type FadeAction struct {
	*webthing.Action
}

func (fade *FadeAction) Generator(thing *webthing.Thing) *webthing.Action {
	fade.Action = webthing.NewAction(uuid.New().String(), thing, "fade", nil, fade.PerformAction, fade.Cancel)
	return fade.Action
}

func (fade *FadeAction) PerformAction() *webthing.Action {
	fmt.Println("Perform fade action…...: ", fade.Name(), " | UUID: ", fade.ID())
	thing := fade.Thing()
	params, _ := fade.Input().MarshalJSON()

	input := make(map[string]interface{})
	if err := json.Unmarshal(params, &input); err != nil {
		fmt.Println(err)
	}
	if brightness, ok := input["brightness"]; ok {
		fmt.Println("Set brightness value: ", brightness)
		thing.Property("brightness").Set(brightness)
	}
	if duration, ok := input["duration"]; ok {
		fmt.Println("Fade duration: ", duration)
		time.Sleep(time.Duration(int64(duration.(float64))) * time.Millisecond)
	}

	event := webthing.NewEvent(thing, "overheated", []byte(fmt.Sprintln(102)))
	thing.AddEvent(event)

	fmt.Println("Fade action Done...", fade.Name())
	return fade.Action
}

func (fade *FadeAction) Cancel() {
	fmt.Println("Cancel fade action...")
}

// Toggles a boolean state on and off.
// Customize a toggles to control the on-off state
type ToggleAction struct {
	*webthing.Action
}

func (toggle *ToggleAction) Generator(thing *webthing.Thing) *webthing.Action {
	toggle.Action = webthing.NewAction(uuid.New().String(), thing, "toggle", nil, toggle.PerformAction, toggle.Cancel)
	return toggle.Action
}

func (toggle *ToggleAction) PerformAction() *webthing.Action {
	fmt.Println("Perform toggle action...: ", toggle.Name(), " | UUID: ", toggle.ID())

	thing := toggle.Thing()
	property := thing.Property("on")
	on := property.Get().(bool)
	property.Set(!on)

	event := webthing.NewEvent(thing, "overheated", []byte(fmt.Sprintln(rand.Intn(100))))
	thing.AddEvent(event)

	fmt.Println("Toggle action done...")
	return toggle.Action
}

func (toggle *ToggleAction) Cancel() {
	fmt.Println("Cancel toggle action...")
}

func onValueForwarder(i interface{}) {
	fmt.Println("Now on statue: ", i)
}
