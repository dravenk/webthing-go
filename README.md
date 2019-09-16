Web of Things 
---
 
 [![GoDoc](https://godoc.org/github.com/dravenk/webthing-go?status.png)](https://godoc.org/github.com/dravenk/webthing-go) 
 [![travis](https://travis-ci.org/dravenk/webthing-go.svg?branch=master)](https://travis-ci.org/dravenk/webthing-go) 
 [![coveralls](https://coveralls.io/repos/dravenk/webthing-go/badge.svg?branch=master&service=github)](https://coveralls.io/github/dravenk/webthing-go?branch=master)
 [![Go Report Card](https://goreportcard.com/badge/github.com/dravenk/webthing-go)](https://goreportcard.com/report/github.com/dravenk/webthing-go)
 
### USAGE:  
To get started look at [examples](https://github.com/dravenk/webthing-go/tree/master/examples) directory:
```
 go get -u -v github.com/dravenk/webthing-go
 cd $GOPATH/src/github.com/dravenk/webthing-go
 
 go run examples/single-thing/single-thing.go
 
 # The default address is http://localhost:8888/things
 # Example: Get a description of a Thing
 curl --request GET --url http://localhost:8888/things
 # Or
 curl --request GET --url http://localhost:8888/things/0
 
 # Example: Get all properties
 curl --request GET --url http://localhost:8888/things/0/properties
 
 # Example: Get a property
 curl --request GET --url http://localhost:8888/things/0/properties/brightness
 
 # Example: Set a property
 curl --request PUT \
   --url http://localhost:8888/things/0/properties/brightness \
   --data '{"brightness": 33}'
  
 # Example: Action Request
 curl --request POST \
   --url http://localhost:8888/things/0/actions \
   --data '{"fade":{"input":{"brightness":55,"duration":2000}}}'
 
 # Example: Action Request
 curl --request POST \
   --url http://localhost:8888/things/0/actions \
   --data '{"toggle":{}}'

 # Example: Actions Queue
 curl --request GET \
   --url http://localhost:8888/things/0/actions
```







RESOURCES
* https://github.com/dravenk/webthing-go/
* https://iot.mozilla.org/framework/
