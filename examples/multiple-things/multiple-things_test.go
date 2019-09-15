package main

import (
	"testing"
)

func TestMultipleThings(t *testing.T) {

}

// multiple-things
//[
//    {
//        "id": "urn:dev:ops:my-lamp-1234",
//        "title": "My Lamp",
//        "@context": "https://iot.mozilla.org/schemas",
//        "@type": [
//            "OnOffSwitch",
//            "Light"
//        ],
//        "properties": {
//            "on": {
//                "@type": "OnOffProperty",
//                "title": "On/Off",
//                "type": "boolean",
//                "description": "Whether the lamp is turned on",
//                "links": [
//                    {
//                        "rel": "property",
//                        "href": "/0/properties/on"
//                    }
//                ]
//            },
//            "brightness": {
//                "@type": "BrightnessProperty",
//                "title": "Brightness",
//                "type": "integer",
//                "description": "The level of light from 0-100",
//                "minimum": 0,
//                "maximum": 100,
//                "unit": "percent",
//                "links": [
//                    {
//                        "rel": "property",
//                        "href": "/0/properties/brightness"
//                    }
//                ]
//            }
//        },
//        "actions": {
//            "fade": {
//                "title": "Fade",
//                "description": "Fade the lamp to a given level",
//                "input": {
//                    "type": "object",
//                    "required": [
//                        "brightness",
//                        "duration"
//                    ],
//                    "properties": {
//                        "brightness": {
//                            "type": "integer",
//                            "minimum": 0,
//                            "maximum": 100,
//                            "unit": "percent"
//                        },
//                        "duration": {
//                            "type": "integer",
//                            "minimum": 1,
//                            "unit": "milliseconds"
//                        }
//                    }
//                },
//                "links": [
//                    {
//                        "rel": "action",
//                        "href": "/0/actions/fade"
//                    }
//                ]
//            }
//        },
//        "events": {
//            "overheated": {
//                "description": "The lamp has exceeded its safe operating temperature",
//                "type": "number",
//                "unit": "degree celsius",
//                "links": [
//                    {
//                        "rel": "event",
//                        "href": "/0/events/overheated"
//                    }
//                ]
//            }
//        },
//        "links": [
//            {
//                "rel": "properties",
//                "href": "/0/properties"
//            },
//            {
//                "rel": "actions",
//                "href": "/0/actions"
//            },
//            {
//                "rel": "events",
//                "href": "/0/events"
//            },
//            {
//                "rel": "alternate",
//                "href": "ws://127.0.0.1:8888/0"
//            }
//        ],
//        "description": "A web connected lamp",
//        "href": "/0",
//        "base": "http://127.0.0.1:8888/0",
//        "securityDefinitions": {
//            "nosec_sc": {
//                "scheme": "nosec"
//            }
//        },
//        "security": "nosec_sc"
//    },
//    {
//        "id": "urn:dev:ops:my-humidity-sensor-1234",
//        "title": "My Humidity Sensor",
//        "@context": "https://iot.mozilla.org/schemas",
//        "@type": [
//            "MultiLevelSensor"
//        ],
//        "properties": {
//            "level": {
//                "@type": "LevelProperty",
//                "title": "Humidity",
//                "type": "number",
//                "description": "The current humidity in %",
//                "minimum": 0,
//                "maximum": 100,
//                "unit": "percent",
//                "readOnly": true,
//                "links": [
//                    {
//                        "rel": "property",
//                        "href": "/1/properties/level"
//                    }
//                ]
//            }
//        },
//        "actions": {},
//        "events": {},
//        "links": [
//            {
//                "rel": "properties",
//                "href": "/1/properties"
//            },
//            {
//                "rel": "actions",
//                "href": "/1/actions"
//            },
//            {
//                "rel": "events",
//                "href": "/1/events"
//            },
//            {
//                "rel": "alternate",
//                "href": "ws://127.0.0.1:8888/1"
//            }
//        ],
//        "description": "A web connected humidity sensor",
//        "href": "/1",
//        "base": "http://127.0.0.1:8888/1",
//        "securityDefinitions": {
//            "nosec_sc": {
//                "scheme": "nosec"
//            }
//        },
//        "security": "nosec_sc"
//    }
//]
