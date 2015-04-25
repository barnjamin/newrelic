package main

import (
	"github.com/barnjamin/newrelic"
	"log"
	"math/rand"
	"os"
	"time"
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
		t, _ := newrelic.NewTransaction("testExternal")

		seg1, _ := t.StartExternalSegment("outToSea", "seg1")
		time.Sleep(time.Duration(50+rand.Intn(50)) * time.Millisecond)
		seg1.End()

		newrelic.RecordMetric("customMetric", 1)

		seg2, _ := t.StartGenericSegment("seg2")
		time.Sleep(10 * time.Millisecond)

		//Try out subsegments
		seg2_1, _ := seg2.StartGenericSegment("seg2_1")
		time.Sleep(time.Duration(50+rand.Intn(50)) * time.Millisecond)
		seg2_1.End()

		time.Sleep(10 * time.Millisecond)
		seg2.End()

		newrelic.RecordMetric("customMetric2", 2)
		seg3, _ := t.StartGenericSegment("seg3")
		sig := make(chan struct{})
		go subfunc(seg3, sig)
		<-sig
		seg3.End()

		t.End()
	}

}

func subfunc(s *newrelic.Segment, sig chan struct{}) {
	subseg, _ := s.StartGenericSegment("seg3_1")
	time.Sleep(time.Duration(20+rand.Intn(10)) * time.Millisecond)
	subseg.End()
	sig <- struct{}{}
}
