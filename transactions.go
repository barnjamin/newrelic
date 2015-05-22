package newrelic

import "log"

type Segmentable interface {
	StartGenericSegment(string) Segmentable
	StartDatastoreSegment(string, string, string, string) Segmentable
	StartExternalSegment(string, string) Segmentable
	End()
}

type SegPromise struct {
	Func func() (Segmentable, error)
	Seg  chan Segmentable
}

type Transaction struct {
	id       int
	funcChan chan func() error
	endChan  chan func() error
	segChan  chan SegPromise
}

func NewTransaction(name string, log *log.Logger) *Transaction {
	tChan := make(chan *Transaction)

	//Force all segments to be created in the same thread as the transaction
	go func() {
		id, err := StartTransaction()

		//We dont want to return nil so create a transaction and the errors will just get logged
		t := &Transaction{
			id:       id,
			funcChan: make(chan func() error),
			endChan:  make(chan func() error),
			segChan:  make(chan SegPromise),
		}

		tChan <- t

		if err != nil {
			log.Printf("Failed to start transaction: %+v", err)
			return
		}

		for {
			select {
			case efun := <-t.funcChan:
				log.Debug("Calling an Update function")
				err := efun()
				if err != nil {
					log.Errorf("Failed to call function: %+v", err)
				}
			case segp := <-t.segChan:
				log.Debug("Calling segment function")
				seg, err := segp.Func()
				if err != nil {
					log.Errorf("Failed to call segmenting function: %+v", err)
				}
				segp.Seg <- seg
			case quit := <-t.endChan:
				log.Debug("Quitting")
				quit()
				return
			}
		}
	}()

	t := <-tChan
	t.SetName(name)
	return t
}

func (t *Transaction) End() {
	t.endChan <- func() error {
		_, err := EndTransaction(t.id)
		return err
	}
}

func (t *Transaction) SetName(name string) {
	t.funcChan <- func() error {
		_, err := SetTransactionName(t.id, name)
		return err
	}
}

func (t *Transaction) SetType(web bool) {
	t.funcChan <- func() error {
		var err error
		switch web {
		case true:
			_, err = SetTransactionTypeWeb(t.id)
		case false:
			_, err = SetTransactionTypeOther(t.id)
		}
		return err
	}
}

func (t *Transaction) SetCategory(cat string) {
	t.funcChan <- func() error {
		_, err := SetTransactionCategory(t.id, cat)
		return err
	}
}

func (t *Transaction) AddAttribute(k, v string) {
	t.funcChan <- func() error {
		_, err := AddTransactionAttribute(t.id, k, v)
		return err
	}
}

func (t *Transaction) SetMaxSegments(max int) {
	t.funcChan <- func() error {
		_, err := SetTransactionMaxSegments(t.id, max)
		return err
	}
}

func (t *Transaction) SetURL(URL string) {
	t.funcChan <- func() error {
		_, err := SetTransactionRequestURL(t.id, URL)
		return err
	}
}

func (t *Transaction) SetError(errType, message, trace, traceDelim string) {
	t.funcChan <- func() error {
		_, err := SetTransactionError(t.id, errType, message, trace, traceDelim)
		return err
	}
}

func (t *Transaction) StartGenericSegment(name string) Segmentable {
	s := SegPromise{
		func() (Segmentable, error) {
			segID, err := StartGenericSegment(t.id, ROOT_SEGMENT, name)
			if err != nil {
				return nil, err
			}
			return &Segment{t.id, segID, t.funcChan, t.segChan}, nil
		},
		make(chan Segmentable),
	}
	t.segChan <- s
	return <-s.Seg
}

func (t *Transaction) StartExternalSegment(host, name string) Segmentable {
	s := SegPromise{
		func() (Segmentable, error) {
			segID, err := StartExternalSegment(t.id, ROOT_SEGMENT, host, name)
			if err != nil {
				return nil, err
			}
			return &Segment{t.id, segID, t.funcChan, t.segChan}, nil
		},
		make(chan Segmentable),
	}
	t.segChan <- s
	return <-s.Seg
}

func (t *Transaction) StartDatastoreSegment(table, operation, sql, rollup_name string) Segmentable {
	s := SegPromise{
		func() (Segmentable, error) {
			subsegmentId, err := StartDatastoreSegment(t.id, ROOT_SEGMENT, table,
				operation, sql, rollup_name)
			if err != nil {
				return nil, err
			}
			return &Segment{t.id, subsegmentId, t.funcChan, t.segChan}, nil
		},
		make(chan Segmentable),
	}
	t.segChan <- s
	return <-s.Seg
}
