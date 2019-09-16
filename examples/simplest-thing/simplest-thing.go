package main

import (
	"fmt"
	"github.com/dravenk/webthing-go"
	"log"
	"net/http"
)

func main() {
	runServer()
}

func makeThing() *webthing.Thing {
	thing := webthing.NewThing("urn:dev:ops:my-actuator-1234",
		"ActuatorExample",
		[]string{"OnOffSwitch"},
		"An actuator example that just log")

	value := webthing.NewValue(true, func(i interface{}) {
		fmt.Println("Change: ", i)
	})
	meta := []byte(`{
	'@type': 'OnOffProperty',
	title: 'On/Off',
	type: 'boolean',
	description: 'Whether the output is changed',
	}`)
	property := webthing.NewProperty(thing, "on", value, meta)
	thing.AddProperty(property)

	return thing

}

func runServer() {
	fmt.Println("Usage:\n" +
		"Try: \n " +
		`curl -X PUT -H 'Content-Type: application/json' --data '{"on": true }' ` + ` http://localhost:8888/things/0/properties/on`)

	thing := makeThing()
	serve := &http.Server{Addr: ":8888"}
	server := webthing.NewWebThingServer(webthing.NewSingleThing(thing), serve)

	log.Fatal(server.Start())
}
