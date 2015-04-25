package newrelic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	NEWRELIC_INSIGHTS_KEY = "YOUR_KEY_HERE"
	URL = "https://insights-collector.newrelic.com/v1/accounts/" //Your URL here
)

var (
	hostname string
)

func init() {
	//if it errors we can leave hostname blank
	hostname, _ = os.Hostname()
}

type EventTracker struct {
	sync.Mutex
	EventType string
	Interval  time.Duration
	Events    []map[string]interface{}
}

func (e *EventTracker) RecordEvent(event map[string]interface{}) error {
	e.Lock()
	defer e.Unlock()

	event["hostname"] = hostname

	if _, ok := event["eventType"]; !ok {
		event["eventType"] = e.EventType
	}

	if _, ok := event["timestamp"]; !ok {
		event["timestamp"] = time.Now().Unix()
	}

	e.Events = append(e.Events, event)
	return nil
}

func (e *EventTracker) Run() error {
	for range time.Tick(e.Interval) {
		err := e.send()
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *EventTracker) send() error {
	e.Lock()
	defer e.Unlock()

	//Dont bother sending if there is nothing to send
	if len(e.Events) == 0 {
		return nil
	}

	b, err := json.MarshalIndent(e.Events, "", " ")
	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST", URL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	r.Header.Add("X-Insert-Key", NEWRELIC_INSIGHTS_KEY)

	//TODO::retry?
	c := &http.Client{}
	result, err := c.Do(r)
	defer result.Body.Close()

	if err != nil {
		return err
	}

	if result.StatusCode != 200 {
		return fmt.Errorf("Invalid response code: %+v", result)
	}

	//reset it
	e.Events = []map[string]interface{}{}

	//TODO::read the body and check for {success:true}?
	return nil
}
