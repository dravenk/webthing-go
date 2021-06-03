package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/dravenk/webthing-go"
	"github.com/google/uuid"
)

func main() {

	light := MakeDimmableLight()

	sensor := FakeGpioHumiditySensor()

	multiple := webthing.NewMultipleThings([]*webthing.Thing{light, sensor}, "LightAndTempDevice")

	httpServer := &http.Server{Addr: "0.0.0.0:8888"}
	server := webthing.NewWebThingServer(multiple, httpServer, "")
	log.Fatal(server.Start())
}

// MakeDimmableLight A dimmable light that logs received commands to stdout.
func MakeDimmableLight() *webthing.Thing {
	// Create a Lamp.
	thing := webthing.NewThing("urn:dev:ops:my-lamp-1234",
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
	onValue := webthing.NewValue(true, func(i interface{}) {
		fmt.Println("On-State is now ", i)
	})
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
		webthing.NewValue(50, func(i interface{}) {
			fmt.Print("Brightness is now ", i)
		}),
		brightnessDescription))

	//Adding a Fade action to this Lamp.
	fadeMeta := []byte(`{
    "title": "Fade",
    "description": "Fade the lamp to a given level",
    "input": {
      "@type": "FadeAction",
      "type": "object",
      "properties": {
        "required": [
          "brightness",
          "duration"
        ],
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

	return thing
}

// Fade the lamp to a given brightness
type FadeAction struct {
	*webthing.Action
}

// Generator  Action generate.
func (fade *FadeAction) Generator(thing *webthing.Thing) *webthing.Action {
	fade.Action = webthing.NewAction(uuid.New().String(), thing, "fade", nil, fade.PerformAction, fade.Cancel)
	return fade.Action
}

// PerformAction Perform an action.
func (fade *FadeAction) PerformAction() *webthing.Action {
	thing := fade.Thing()
	params, _ := fade.Input().MarshalJSON()

	input := make(map[string]interface{})
	if err := json.Unmarshal(params, &input); err != nil {
		fmt.Println(err)
	}
	if brightness, ok := input["brightness"]; ok {
		thing.Property("brightness").Set(brightness)
	}
	if duration, ok := input["duration"]; ok {
		time.Sleep(time.Duration(int64(duration.(float64))) * time.Millisecond)
	}

	event := webthing.NewEvent(thing, "overheated", []byte(fmt.Sprintln(102)))
	thing.AddEvent(event)

	return fade.Action
}

func (fade *FadeAction) Cancel() {}

// FakeGpioHumiditySensor A humidity sensor which updates its measurement every few seconds.
func FakeGpioHumiditySensor() *webthing.Thing {
	thing := webthing.NewThing(
		"urn:dev:ops:my-humidity-sensor-1234",
		"My Humidity Sensor",
		[]string{"MultiLevelSensor"},
		"A web connected humidity sensor")

	level := webthing.NewValue(0.0)
	levelDescription := []byte(`{
        "@type": "LevelProperty",
        "title": "Humidity",
        "type": "number",
        "description": "The current humidity in %",
        "minimum": 0,
        "maximum": 100,
        "unit": "percent",
        "readOnly": true
  }`)
	thing.AddProperty(webthing.NewProperty(
		thing,
		"level",
		level,
		levelDescription))

	go func(level webthing.Value) {
		for {
			time.Sleep(3000 * time.Millisecond)
			newLevel := readFromGPIO()
			fmt.Println("setting new humidity level:", newLevel)
			level.NotifyOfExternalUpdate(newLevel)
		}
	}(level)

	return thing
}

// Mimic an actual sensor updating its reading every couple seconds.
func readFromGPIO() float64 {
	return math.Abs(rand.Float64() * 70.0 * 0.5 * rand.NormFloat64())
}
