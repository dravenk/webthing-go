Web of Things 
---
 [![GitHub forks](https://img.shields.io/github/forks/dravenk/webthing-go.svg?style=social&label=Fork&maxAge=2592000)](https://GitHub.com/dravenk/webthing-go/network/)
 [![GitHub version](https://badge.fury.io/gh/dravenk%2Fwebthing-go.svg)](https://badge.fury.io/gh/dravenk%2Fwebthing-go)
 [![GoDoc](https://godoc.org/github.com/dravenk/webthing-go?status.png)](https://godoc.org/github.com/dravenk/webthing-go) 
 [![Codacy Badge](https://api.codacy.com/project/badge/Grade/bef38274a3cb4156b374bb76dc1670e5)](https://www.codacy.com/manual/dravenk/webthing-go?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dravenk/webthing-go&amp;utm_campaign=Badge_Grade) 
 [![travis](https://api.travis-ci.org/dravenk/webthing-go.svg?branch=master)](https://travis-ci.com/dravenk/webthing-go) 
 [![coveralls](https://coveralls.io/repos/dravenk/webthing-go/badge.svg?branch=master&service=github)](https://coveralls.io/github/dravenk/webthing-go?branch=master)
 [![Go Report Card](https://goreportcard.com/badge/github.com/dravenk/webthing-go)](https://goreportcard.com/report/github.com/dravenk/webthing-go)
 [![codebeat badge](https://codebeat.co/badges/090b9189-b20c-4910-8ff2-d7c12a28e55f)](https://codebeat.co/projects/github-com-dravenk-webthing-go-master)

### USAGE:  
To get started look at [examples](https://github.com/dravenk/webthing-go/tree/master/examples) directory:
```
 go get -u -v github.com/dravenk/webthing-go
 cd $GOPATH/src/github.com/dravenk/webthing-go
 
 go run examples/single-thing/single-thing.go
 
 # The default address is http://localhost:8888/things
 curl --request GET --url http://localhost:8888/things
 
 # Example: Get a description of a Thing
 # use the Thing index.
 curl --request GET --url http://localhost:8888/things/0
 # Or use the Thing title.
 curl --request GET --url http://localhost:8888/things/Lamp

 # Example: Get all properties
 curl --request GET --url http://localhost:8888/things/0/properties
 # Or
 curl --request GET --url http://localhost:8888/things/Lamp/properties

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

 # Example: Cancel an Action Request
    curl --request DELETE \
      --url http://localhost:8888/things/0/actions/fade/{action_id}

 # Example: Action Request
 curl --request POST \
   --url http://localhost:8888/things/0/actions \
   --data '{"toggle":{}}'

 # Example: Actions Queue
 curl --request GET \
   --url http://localhost:8888/things/0/actions
   
 # Example: Events Request
 curl --request GET \
   --url http://localhost:8888/things/0/events
   
  # Example: Event Request
  curl --request GET \
    --url http://localhost:8888/things/0/events/overheated

```







RESOURCES
* https://github.com/dravenk/webthing-go/
* https://iot.mozilla.org/framework/
