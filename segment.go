package newrelic

type Segment struct {
	tid, id  int
	funcChan chan func() error
	segChan  chan SegPromise
}

func (s *Segment) StartGenericSegment(name string) Segmentable {
	sp := SegPromise{
		func() (Segmentable, error) {
			subsegmentId, err := StartGenericSegment(s.tid, s.id, name)
			if err != nil {
				return nil, err
			}
			return &Segment{s.tid, subsegmentId, s.funcChan, s.segChan}, nil
		},
		make(chan Segmentable),
	}
	s.segChan <- sp
	return <-sp.Seg

}

func (s *Segment) StartExternalSegment(host, name string) Segmentable {
	sp := SegPromise{
		func() (Segmentable, error) {
			subsegmentId, err := StartExternalSegment(s.tid, s.id, host, name)
			if err != nil {
				return nil, err
			}
			return &Segment{s.tid, subsegmentId, s.funcChan, s.segChan}, nil
		},
		make(chan Segmentable),
	}
	s.segChan <- sp
	return <-sp.Seg
}

func (s *Segment) StartDatastoreSegment(table, operation, sql, rollup_name string) Segmentable {
	sp := SegPromise{
		func() (Segmentable, error) {
			subsegmentId, err := StartDatastoreSegment(s.tid, s.id, table, operation, sql, rollup_name)
			if err != nil {
				return nil, err
			}
			return &Segment{s.tid, subsegmentId, s.funcChan, s.segChan}, nil
		},
		make(chan Segmentable),
	}
	s.segChan <- sp
	return <-sp.Seg
}

func (s *Segment) End() {
	s.funcChan <- func() error {
		_, err := EndSegment(s.tid, s.id)
		return err
	}
}
