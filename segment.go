package newrelic

type Segment struct {
	tid, id int
}

func (s *Segment) StartGenericSegment(name string) (*Segment, error) {
	subsegmentId, err := StartGenericSegment(s.tid, s.id, name)
	if err != nil {
		return nil, err
	}
	return &Segment{s.tid, subsegmentId}, nil
}

func (s *Segment) StartExternalSegment(host, name string) (*Segment, error) {
	subsegmentId, err := StartExternalSegment(s.tid, s.id, host, name)
	if err != nil {
		return nil, err
	}
	return &Segment{s.tid, subsegmentId}, nil
}

func (s *Segment) StartDatastoreSegment(table, operation, sql, rollup_name string) (*Segment, error) {
	subsegmentId, err := StartDatastoreSegment(s.tid, s.id, table, operation, sql, rollup_name)
	if err != nil {
		return nil, err
	}
	return &Segment{s.tid, subsegmentId}, nil
}

func (s *Segment) End() error {
	_, err := EndSegment(s.tid, s.id)
	return err
}
