package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/barnjamin/newrelic"
)

func main() {
	log.Println("Starting")

	rand.Seed(42)

	//This path would have to be updated with the location of the properties file
	if err := os.Setenv("NEWRELIC_LOG_PROPERTIES_FILE", "../newrelic_sdk/config/log4cplus.properties"); err != nil {
		log.Fatalf("Failzore!")
	}

	//Initialize newrelic with name of application being monitored
	newrelic.Initialize("TestApp")

	for {

		log.Println("New transaction")

		//Start a new transaction and pass a logger for errors
		t := newrelic.NewTransaction("testExternal", &log.Logger{})

		//Start an external segment, get back a Segmentable interface
		seg1 := t.StartExternalSegment("outToSea", "seg1")
		time.Sleep(time.Duration(50+rand.Intn(50)) * time.Millisecond)
		seg1.End()

		//Record a custom metric
		newrelic.RecordMetric("customMetric", 1)

		seg2 := t.StartGenericSegment("seg2")
		time.Sleep(10 * time.Millisecond)

		//Try out subsegments
		seg2_1 := seg2.StartGenericSegment("seg2_1")
		time.Sleep(time.Duration(50+rand.Intn(50)) * time.Millisecond)
		seg2_1.End()

		time.Sleep(10 * time.Millisecond)
		seg2.End()

		newrelic.RecordMetric("customMetric2", 2)

		sig := make(chan struct{})
		seg3 := t.StartGenericSegment("seg3")
		//Try out segmenting within a goroutine
		go subfunc(seg3, sig)

		//block until we return
		<-sig
		seg3.End()

		//end the transaction
		t.End()
	}

}

func subfunc(s newrelic.Segmentable, sig chan struct{}) {
	//Start tracking at the beginning and defer the End until the function returns
	defer s.StartGenericSegment("seg3_1").End()
	time.Sleep(time.Duration(20+rand.Intn(10)) * time.Millisecond)
	sig <- struct{}{}
}
