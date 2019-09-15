package main

import (
	"encoding/json"
	"github.com/dravenk/webthing-go"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSingleThing(t *testing.T) {

	thing := MakeThing()
	singleThing := webthing.NewSingleThing(thing)

	basePath := "/things"
	server := webthing.NewWebThingServer(singleThing, &http.Server{})

	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	// Test api: GET /things
	res, err := http.Get(ts.URL + basePath)
	if err != nil || res.StatusCode != http.StatusOK {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	st := &SingleThingStruct{}
	if err = json.Unmarshal(body, &st); err != nil {
		log.Fatal("Unmarshal body error: GET /things ", err.Error())
	}

	if st.ID == "" {
		log.Fatal("ID is empty. ")
	}

	if st.Properties.Brightness.Title != "Brightness" {
		log.Fatal("Not found Brightness property: ", st.Properties.Brightness)
	}
	if st.Properties.Brightness.Links == nil {
		log.Fatal("Not found links of Brightness property: ", st.Properties.Brightness)
	}
	if st.Properties.On.Links == nil {
		log.Fatal("Not found links of On property. ", st.Properties.On)
	}

	if st.Properties.On.Type == "" {
		log.Fatal("Not found type of On property. ", st.Properties.On)
	}

	if st.Actions.Fade.Title == "" {
		log.Fatal("Not found Fade of action. ", st.Actions)
	}

	if st.Actions.Fade.Input.Properties.Brightness.Type != "integer" {
		log.Fatal("Type of brightness in FadeAction input is wrong. ", st.Actions)
	}

	if st.Links == nil {
		log.Fatal("Links object is nil in thing. ")
	}

	// Test api: GET /things/0/properties
	res, err = http.Get(ts.URL + basePath + "/0/properties")
	if err != nil || res.StatusCode != http.StatusOK {
		log.Fatal(err)
	}
	properties := make(map[string]interface{})
	body, err = ioutil.ReadAll(res.Body)
	if err = json.Unmarshal(body, &properties); err != nil {
		log.Fatal(err)
	}
	if on, ok := properties["on"]; !ok {
		log.Fatal("Not found On properties: GET /things/0/properties")
	} else {
		if !on.(bool) {
			log.Fatal("On properties is not boolean: GET /things/0/properties ", on)
		}
	}
	if _, ok := properties["brightness"]; !ok {
		log.Fatal("Not found brightness properties: GET /things/0/properties")
	}

	// Test api: GET /things/0/properties/brightness
	res, err = http.Get(ts.URL + basePath + "/0/properties/brightness")
	if err != nil || res.StatusCode != http.StatusOK {
		log.Fatal(err)
	}
	brightnessRes := make(map[string]json.Number)
	body, _ = ioutil.ReadAll(res.Body)
	if err = json.Unmarshal(body, &brightnessRes); err != nil {
		log.Fatal("Not found brightness properties: GET /things/0/properties/brightness", string(body), " Err: ", err)
	}

	// Test api: PUT /things/0/properties/brightness
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/things/0/properties/brightness", strings.NewReader(`{"brightness": 66 }`))
	res, _ = http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		log.Fatal(err)
	}
	brightnessRes = make(map[string]json.Number)
	body, _ = ioutil.ReadAll(res.Body)
	if err = json.Unmarshal(body, &brightnessRes); err != nil {
		log.Fatal("Not found brightness properties: GET /things/0/properties")
	}

	// Check brightness value
	brightness := thing.Property("brightness").Get()
	// Last test was set brightness property is 66
	if brightness.(float64) != 66 {
		log.Fatal("Thing brightness is : ", brightness)
	}

	contentType := `application/json`
	// Test api: POST /things/0/actions
	res, err = http.Post(ts.URL+basePath+"/0/actions", contentType, strings.NewReader(`{
    "fade": {
        "input": {
            "brightness": 33,
            "duration": 2000
        }
    }
	}`))
	if err != nil || res.StatusCode != http.StatusOK {
		log.Fatal(err)
	}
	actionsRes := struct {
		Fade struct {
			Input struct {
				Brightness float64 `json:"brightness"`
				Duration   int     `json:"duration"`
			} `json:"input"`
		} `json:"fade"`
		Status string `json:"status"`
	}{}
	body, _ = ioutil.ReadAll(res.Body)
	if err = json.Unmarshal(body, &actionsRes); err != nil {
		log.Fatal("Create fade action wrong: POST /things/0/actions", err)
	}
	brightnessV := actionsRes.Fade.Input.Brightness
	// Check response.
	if brightnessV == brightness.(float64) || brightnessV != 33 {
		log.Fatal("Brightness value not change. POST /things/0/actions. Brightness is : ", brightnessV)
	}
	// Check brightness value if changed after perform action.
	// Verify that the values are checked when the program is complete
	// Sleep "duration": 2000 *time.Millisecond
	time.Sleep(2000 * time.Millisecond)
	brightnessAfterChange := thing.Property("brightness").Get()
	// Last test was set brightness property is 33
	if brightnessAfterChange.(float64) != 33 {
		log.Fatal("Thing brightness is : ", brightnessAfterChange)
	}

	// Perform this action 100 times and count the number of executions
	for i := 0; i < 100; i++ {
		res, err = http.Post(ts.URL+basePath+"/0/actions", contentType, strings.NewReader(`{"toggle": {}}`))
		if err != nil || res.StatusCode != http.StatusOK {
			log.Fatal("Perform action wrong. POST /things/0/actions ", err)
		}
	}
	// Test api: GET /things/0/actions//toggle
	res, err = http.Get(ts.URL + basePath + "/0/actions/toggle")
	if err != nil || res.StatusCode != http.StatusOK {
		log.Fatal("Get Action resources wrong.", res, " Error: ", err)
	}
	var countActions []interface{}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Get Action resources wrong.", res, " Error: ", err)
	}
	if err = json.Unmarshal(body, &countActions); err != nil {
		log.Fatal("Get Action resources wrong. GE/things/0/actions/toggle.  Count of actions :", len(countActions), " Error: ", err)
	}
	if len(countActions) != 100 {
		log.Fatal("Count of actions :", len(countActions))
	}
}

// See example: https://github.com/mozilla-iot/webthing-node
//{
//    "id": "urn:dev:ops:my-lamp-1234",
//    "title": "My Lamp",
//    "@context": "https://iot.mozilla.org/schemas",
//    "@type": [
//        "OnOffSwitch",
//        "Light"
//    ],
//    "properties": {
//        "on": {
//            "@type": "OnOffProperty",
//            "title": "On/Off",
//            "type": "boolean",
//            "description": "Whether the lamp is turned on",
//            "links": [
//                {
//                    "rel": "property",
//                    "href": "/0/properties/on"
//                }
//            ]
//        },
//        "brightness": {
//            "@type": "BrightnessProperty",
//            "title": "Brightness",
//            "type": "integer",
//            "description": "The level of light from 0-100",
//            "minimum": 0,
//            "maximum": 100,
//            "unit": "percent",
//            "links": [
//                {
//                    "rel": "property",
//                    "href": "/0/properties/brightness"
//                }
//            ]
//        }
//    },
//    "actions": {
//        "fade": {
//            "title": "Fade",
//            "description": "Fade the lamp to a given level",
//            "input": {
//                "type": "object",
//                "required": [
//                    "brightness",
//                    "duration"
//                ],
//                "properties": {
//                    "brightness": {
//                        "type": "integer",
//                        "minimum": 0,
//                        "maximum": 100,
//                        "unit": "percent"
//                    },
//                    "duration": {
//                        "type": "integer",
//                        "minimum": 1,
//                        "unit": "milliseconds"
//                    }
//                }
//            },
//            "links": [
//                {
//                    "rel": "action",
//                    "href": "/0/actions/fade"
//                }
//            ]
//        }
//    },
//    "events": {
//        "overheated": {
//            "description": "The lamp has exceeded its safe operating temperature",
//            "type": "number",
//            "unit": "degree celsius",
//            "links": [
//                {
//                    "rel": "event",
//                    "href": "/0/events/overheated"
//                }
//            ]
//        }
//    },
//    "links": [
//        {
//            "rel": "properties",
//            "href": "/0/properties"
//        },
//        {
//            "rel": "actions",
//            "href": "/0/actions"
//        },
//        {
//            "rel": "events",
//            "href": "/0/events"
//        },
//        {
//            "rel": "alternate",
//            "href": "ws://127.0.0.1:8888/0"
//        }
//    ],
//    "description": "A web connected lamp",
//    "base": "http://127.0.0.1:8888/0",
//    "securityDefinitions": {
//        "nosec_sc": {
//            "scheme": "nosec"
//        }
//    },
//    "security": "nosec_sc"
//}

type SingleThingStruct struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Context    string   `json:"@context"`
	Type       []string `json:"@type"`
	Properties struct {
		On struct {
			AtType      string `json:"@type"`
			Title       string `json:"title"`
			Type        string `json:"type"`
			Description string `json:"description"`
			Links       []struct {
				Rel  string `json:"rel"`
				Href string `json:"href"`
			} `json:"links"`
		} `json:"on"`
		Brightness struct {
			AtType      string `json:"@type"`
			Title       string `json:"title"`
			Type        string `json:"type"`
			Description string `json:"description"`
			Minimum     int    `json:"minimum"`
			Maximum     int    `json:"maximum"`
			Unit        string `json:"unit"`
			Links       []struct {
				Rel  string `json:"rel"`
				Href string `json:"href"`
			} `json:"links"`
		} `json:"brightness"`
	} `json:"properties"`
	Actions struct {
		Fade struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Input       struct {
				Type       string   `json:"type"`
				Required   []string `json:"required"`
				Properties struct {
					Brightness struct {
						Type    string `json:"type"`
						Minimum int    `json:"minimum"`
						Maximum int    `json:"maximum"`
						Unit    string `json:"unit"`
					} `json:"brightness"`
					Duration struct {
						Type    string `json:"type"`
						Minimum int    `json:"minimum"`
						Unit    string `json:"unit"`
					} `json:"duration"`
				} `json:"properties"`
			} `json:"input"`
			Links []struct {
				Rel  string `json:"rel"`
				Href string `json:"href"`
			} `json:"links"`
		} `json:"fade"`
	} `json:"actions"`
	Events struct {
		Overheated struct {
			Description string `json:"description"`
			Type        string `json:"type"`
			Unit        string `json:"unit"`
			Links       []struct {
				Rel  string `json:"rel"`
				Href string `json:"href"`
			} `json:"links"`
		} `json:"overheated"`
	} `json:"events"`
	Links []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
	Description         string `json:"description"`
	Base                string `json:"base"`
	SecurityDefinitions struct {
		NosecSc struct {
			Scheme string `json:"scheme"`
		} `json:"nosec_sc"`
	} `json:"securityDefinitions"`
	Security string `json:"security"`
}
