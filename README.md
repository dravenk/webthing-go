Web of Things
---
 [![GitHub forks](https://img.shields.io/github/forks/dravenk/webthing-go.svg?style=social&label=Fork&maxAge=2592000)](https://GitHub.com/dravenk/webthing-go/network/)
 [![GitHub version](https://badge.fury.io/gh/dravenk%2Fwebthing-go.svg)](https://badge.fury.io/gh/dravenk%2Fwebthing-go)
 [![GoDoc](https://godoc.org/github.com/dravenk/webthing-go?status.png)](https://godoc.org/github.com/dravenk/webthing-go) 
 [![Codacy Badge](https://api.codacy.com/project/badge/Grade/bef38274a3cb4156b374bb76dc1670e5)](https://www.codacy.com/manual/dravenk/webthing-go?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dravenk/webthing-go&amp;utm_campaign=Badge_Grade) 
 [![travis](https://api.travis-ci.org/dravenk/webthing-go.svg?branch=master)](https://travis-ci.com/dravenk/webthing-go) 
 [![Go Report Card](https://goreportcard.com/badge/github.com/dravenk/webthing-go)](https://goreportcard.com/report/github.com/dravenk/webthing-go)
 [![codebeat badge](https://codebeat.co/badges/090b9189-b20c-4910-8ff2-d7c12a28e55f)](https://codebeat.co/projects/github-com-dravenk-webthing-go-master)
 [![Build Status](https://img.shields.io/docker/cloud/build/dravenk/webthing.svg)](https://cloud.docker.com/repository/docker/dravenk/webthing/builds)

### USAGE
This library fully supports [Web Thing REST API](https://iot.mozilla.org/wot/).You can start building your Web of Thing by looking at [single-thing](https://github.com/dravenk/webthing-go/blob/master/examples/single-thing/single-thing.go). 

#### Download and import
This package name is called `webthing`. This project is called webthing-go to keep the naming consistent with the implementation of other languages. You just need to import this package the way golang normally imports a package.
```
go get -u -v github.com/dravenk/webthing-go
```
```go
import (
	"github.com/dravenk/webthing-go"
)
```

#### Create Thing
```go
// Create a Lamp.
thing := webthing.NewThing("urn:dev:ops:my-thing-1234",
	"Lamp",
	[]string{"OnOffSwitch", "Light"},
	"A web connected thing")
```
For more information on Creating Webthing, please check the wiki [Create-Thing](https://github.com/dravenk/webthing-go/wiki/Create-Thing)

#### Example
```
cd $GOPATH/src/github.com/dravenk/webthing-go
go run examples/single-thing/single-thing.go
```
You can also run a sample with [docker](https://hub.docker.com/r/dravenk/webthing):
```
docker run -ti --name single-thing -p 8888:8888 dravenk/webthing
```
For more information on Run Example, please check the wiki [Run-Example](https://github.com/dravenk/webthing-go/wiki/Run-example)

RESOURCES
  * https://github.com/dravenk/webthing-go/
  * https://iot.mozilla.org/framework/
