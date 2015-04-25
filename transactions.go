package newrelic

type Segmentable interface {
	StartGenericSegment(string) (Segmentable, error)
	StartDatastoreSegment(string, string, string, string) (Segmentable, error)
	StartExternalSegment(string, string) (Segmentable, error)
	End() (Segmentable, error)
}

type Transaction struct {
	id int
}

func NewTransaction(name string) (*Transaction, error) {
	id, err := StartTransaction()
	if err != nil {
		return nil, err
	}

	t := &Transaction{id}
	t.SetName(name)
	return t, nil
}

func (t *Transaction) SetName(name string) error {
	_, err := SetTransactionName(t.id, name)
	return err
}

func (t *Transaction) SetType(web bool) error {
	var err error
	switch web {
	case true:
		_, err = SetTransactionTypeWeb(t.id)
	case false:
		_, err = SetTransactionTypeOther(t.id)
	}
	return err
}

func (t *Transaction) SetCategory(cat string) error {
	_, err := SetTransactionCategory(t.id, cat)
	return err
}

func (t *Transaction) AddAttribute(k, v string) error {
	_, err := AddTransactionAttribute(t.id, k, v)
	return err
}

func (t *Transaction) SetMaxSegments(max int) error {
	_, err := SetTransactionMaxSegments(t.id, max)
	return err
}

func (t *Transaction) SetURL(URL string) error {
	_, err := SetTransactionRequestURL(t.id, URL)
	return err
}

func (t *Transaction) SetError(errType, message, trace, traceDelim string) error {
	_, err := SetTransactionError(t.id, errType, message, trace, traceDelim)
	return err
}

func (t *Transaction) StartGenericSegment(name string) (*Segment, error) {
	segID, err := StartGenericSegment(t.id, ROOT_SEGMENT, name)
	if err != nil {
		return nil, err
	}
	return &Segment{t.id, segID}, nil
}

func (t *Transaction) StartExternalSegment(host, name string) (*Segment, error) {
	segID, err := StartExternalSegment(t.id, ROOT_SEGMENT, host, name)
	if err != nil {
		return nil, err
	}
	return &Segment{t.id, segID}, nil
}

func (t *Transaction) StartDatastoreSegment(table, operation, sql, rollup_name string) (*Segment, error) {
	subsegmentId, err := StartDatastoreSegment(t.id, ROOT_SEGMENT, table, operation, sql, rollup_name)
	if err != nil {
		return nil, err
	}
	return &Segment{t.id, subsegmentId}, nil
}

func (t *Transaction) End() error {
	_, err := EndTransaction(t.id)
	return err
}
