Web of Things 
---
 [![GitHub forks](https://img.shields.io/github/forks/dravenk/webthing-go.svg?style=social&label=Fork&maxAge=2592000)](https://GitHub.com/dravenk/webthing-go/network/)
 [![GitHub version](https://badge.fury.io/gh/dravenk%2Fwebthing-go.svg)](https://badge.fury.io/gh/dravenk%2Fwebthing-go)
 [![GoDoc](https://godoc.org/github.com/dravenk/webthing-go?status.png)](https://godoc.org/github.com/dravenk/webthing-go) 
 [![Codacy Badge](https://api.codacy.com/project/badge/Grade/bef38274a3cb4156b374bb76dc1670e5)](https://www.codacy.com/manual/dravenk/webthing-go?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dravenk/webthing-go&amp;utm_campaign=Badge_Grade) 
 [![travis](https://api.travis-ci.org/dravenk/webthing-go.svg?branch=master)](https://travis-ci.com/dravenk/webthing-go) 
 [![Go Report Card](https://goreportcard.com/badge/github.com/dravenk/webthing-go)](https://goreportcard.com/report/github.com/dravenk/webthing-go)
 [![codebeat badge](https://codebeat.co/badges/090b9189-b20c-4910-8ff2-d7c12a28e55f)](https://codebeat.co/projects/github-com-dravenk-webthing-go-master)

### USAGE:  

You can start building your Web of Thing by looking at [single-thing](https://github.com/dravenk/webthing-go/blob/master/examples/single-thing/single-thing.go)

#### Download and import:
```
 go get -u -v github.com/dravenk/webthing-go
```
This package is called webthing. You just need to import this package the way golang normally imports a package.

```go
import (
	"github.com/dravenk/webthing-go"
)
```

#### Create Thing:
```go
// Create a Lamp.
thing := webthing.NewThing("urn:dev:ops:my-thing-1234",
	"Lamp",
	[]string{"OnOffSwitch", "Light"},
	"A web connected thing")

```
Before creating OnOffProperty you need to create the Forwarder method of OnOff Value. The method that updates the actual value on the thing
Example:
```go
func onValueForwarder(i interface{}) {
    fmt.Println("Now on statue: ", i)
}
```
Create an onValue with default value:
```go
onValue := webthing.NewValue(true, onValueForwarder)
```
```go
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
```
Create an action. The methods you have to implement are:
```go
// Custom Action need create a Generator to generate a action.
// The application will invoke the Action created by the Generator method.
// This is very similar to simply constructor.
// See thing.PerformAction()*Action
Generator(thing *Thing) *Action

// Override this with the code necessary to perform the action.
PerformAction() *Action

// Override this with the code necessary to cancel the action.
Cancel()
```
The Action Request can be used to retrieve input data in the following way like the fade.Thing().Input(). Here is an example:
```go
type FadeAction struct {
	*webthing.Action
}

func (fade *FadeAction) Generator(thing *webthing.Thing) *webthing.Action {
	fade.Action = webthing.NewAction(uuid.New().String(), thing, "fade", nil, fade.PerformAction, fade.Cancel)
	return fade.Action
}

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
	return fade.Action
}

func (fade *FadeAction) Cancel() {}
```

You can find an example of creating Thing above in [single-thing](https://github.com/dravenk/webthing-go/blob/master/examples/single-thing/single-thing.go). 
You can find more Examples in the [Examples](https://github.com/dravenk/webthing-go/tree/master/examples) directory

#### Run example:

```
 cd $GOPATH/src/github.com/dravenk/webthing-go
 
 go run examples/single-thing/single-thing.go
 ```
 
 All [Web Thing REST API](https://iot.mozilla.org/wot/#web-thing-rest-api) are currently supported.
 The currently built Server supports lookup using Thing's title or index in Things.
 
 ```
 # The default address is http://localhost:8888/things
 curl --request GET --url http://localhost:8888/things
 
 # Example: Get a description of a Thing
 # use the Thing index.
 curl --request GET --url http://localhost:8888
 # Or use the Thing title.
 curl --request GET --url http://localhost:8888/Lamp
```

 Example: Properties
 ```
 # Example: Get all properties
 curl --request GET --url http://localhost:8888/properties
 # Or
 curl --request GET --url http://localhost:8888/Lamp/properties

 # Example: Get a property
 curl --request GET --url http://localhost:8888/properties/brightness
 
 # Example: Set a property
 curl --request PUT \
   --url http://localhost:8888/properties/brightness \
   --data '{"brightness": 33}'
```

 Example: Actions 
```
 # Example: Action Request
 curl --request POST \
   --url http://localhost:8888/actions \
   --data '{"fade":{"input":{"brightness":55,"duration":2000}}}'

 # Example: Cancel an Action Request
    curl --request DELETE \
      --url http://localhost:8888/actions/fade/{action_id}

 # Example: Action Request
 curl --request POST \
   --url http://localhost:8888/actions \
   --data '{"toggle":{}}'

 # Example: Actions Queue
 curl --request GET \
   --url http://localhost:8888/actions
 ```

 Example: Events
 ```  
 # Example: Events Request
 curl --request GET \
   --url http://localhost:8888/events
   
  # Example: Event Request
  curl --request GET \
    --url http://localhost:8888/events/overheated
```







RESOURCES
* https://github.com/dravenk/webthing-go/
* https://iot.mozilla.org/framework/
