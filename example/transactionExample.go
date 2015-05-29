package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/barnjamin/newrelic"
)

func main() {
	log.Printf("Starting")

	rand.Seed(42)
	if err := os.Setenv("NEWRELIC_LOG_PROPERTIES_FILE", "../newrelic_sdk/config/log4cplus.properties"); err != nil {
		log.Fatalf("Failzore!")
	}

	newrelic.Initialize("TestApp")

	for {

		log.Printf("new transaction")
		t := newrelic.NewTransaction("testExternal", &log.Logger{})

		seg1 := t.StartExternalSegment("outToSea", "seg1")
		time.Sleep(time.Duration(50+rand.Intn(50)) * time.Millisecond)
		seg1.End()

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
		seg3 := t.StartGenericSegment("seg3")
		sig := make(chan struct{})
		go subfunc(seg3, sig)
		<-sig
		seg3.End()

		t.End()
	}

}

func subfunc(s newrelic.Segmentable, sig chan struct{}) {
	defer s.StartGenericSegment("seg3_1").End()
	time.Sleep(time.Duration(20+rand.Intn(10)) * time.Millisecond)
	sig <- struct{}{}
}
