// Package logmtest provides utilities for Log testing.
package logmtest

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

// NewRecorder returns a new instance of an in-memory recorder.
func NewRecorder() *Recorder {
	return &Recorder{
		buf: new(bytes.Buffer),
		mu:  sync.Mutex{},
	}
}

// Recorder is an in-memory recorder dedicated fo test purposes.
type Recorder struct {
	// ExpectAnyOrder gives an option whether to match all expectations in the order they were set or not.
	// By default, expectations are in the order they were set.
	// But when using goroutines, that option may be handy.
	ExpectAnyOrder bool

	// ExpectUnexpected gives an option whether to ignore not matching records.
	// By default, if a record is not expected, an error will be triggered.
	ExpectUnexpected bool

	buf *bytes.Buffer
	mu  sync.Mutex
}

// Write implements the io.Writer interface.
func (l *Recorder) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.Write(p)
}

// Expect returns if the recorded data matches these records.
func (l *Recorder) Expect(list ...Record) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	var (
		sc     = bufio.NewScanner(l.buf)
		rs     = newRecords(list)
		pp, cp int
		ok     bool
	)
	for sc.Scan() {
		cp, ok = rs.Contains(sc.Bytes())
		if !ok {
			if !l.ExpectUnexpected {
				return fmt.Errorf("any records are not within %q: %w", sc.Text(), strconv.ErrRange)
			}
			continue
		}
		if cp < pp && !l.ExpectAnyOrder {
			return fmt.Errorf("%q: %w", sc.Text(), strconv.ErrRange)
		}
		pp = cp
	}
	return errors.Join(rs.Err(), sc.Err())
}

// Record represents a record.
type Record struct {
	// Contains is the substr to find.
	Contains string

	ok bool
}

func newRecords(a []Record) records {
	r := make(records, len(a))
	copy(r, a)
	return r
}

type records []Record

// Contains reports whether b is contained by one of these records.
func (rs records) Contains(b []byte) (int, bool) {
	for k := range rs {
		if rs[k].ok {
			continue
		}
		rs[k].ok = bytes.Contains(b, []byte(rs[k].Contains))
		if rs[k].ok {
			return k, true
		}
	}
	return 0, false
}

// Err returns in error any records that has not been found.
func (rs records) Err() (err error) {
	for k, v := range rs {
		if !v.ok {
			err = errors.Join(err, fmt.Errorf("record #%d %q: %w", k, v.Contains, strconv.ErrRange))
		}
	}
	return err
}
