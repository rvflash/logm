package logmtest_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/logm"
	"github.com/rvflash/logm/logmtest"
)

const (
	msg01 = "hello world"
	msg02 = "another world"
	msg03 = "the world"
)

func TestRecorder_Expect(t *testing.T) {
	t.Parallel()

	are := is.New(t)
	for name, tc := range map[string]struct {
		// inputs
		rec *logmtest.Recorder
		in  []string
		// outputs
		out []logmtest.Record
		err error
	}{
		"Default": {rec: logmtest.NewRecorder()},
		"Default settings: one logged, none expected": {
			rec: logmtest.NewRecorder(),
			in:  []string{msg01},
			err: strconv.ErrRange,
		},
		"Default settings: one logged, one expected": {
			rec: logmtest.NewRecorder(),
			in:  []string{msg01},
			out: []logmtest.Record{{Contains: msg01}},
		},
		"Default settings: one logged, two expected": {
			rec: logmtest.NewRecorder(),
			in:  []string{msg01},
			out: []logmtest.Record{{Contains: msg01}, {Contains: msg02}},
			err: strconv.ErrRange,
		},
		"Default settings: two logged, two expected": {
			rec: logmtest.NewRecorder(),
			in:  []string{msg01, msg02},
			out: []logmtest.Record{{Contains: msg01}, {Contains: msg01}},
			err: strconv.ErrRange,
		},
		"Default settings: three logged, three expected": {
			rec: logmtest.NewRecorder(),
			in:  []string{msg01, msg02, msg03},
			out: []logmtest.Record{{Contains: msg01}, {Contains: msg03}, {Contains: msg02}},
			err: strconv.ErrRange,
		},
		"Any order: two logged, two expected ordered": {
			rec: func() *logmtest.Recorder {
				rec := logmtest.NewRecorder()
				rec.ExpectAnyOrder = true
				return rec
			}(),
			in:  []string{msg01, msg02},
			out: []logmtest.Record{{Contains: msg01}, {Contains: msg02}},
		},
		"Any order: two logged, two expected unordered": {
			rec: func() *logmtest.Recorder {
				rec := logmtest.NewRecorder()
				rec.ExpectAnyOrder = true
				return rec
			}(),
			in:  []string{msg02, msg01},
			out: []logmtest.Record{{Contains: msg01}, {Contains: msg02}},
		},
		"Any order: one logged, two expected": {
			rec: func() *logmtest.Recorder {
				rec := logmtest.NewRecorder()
				rec.ExpectAnyOrder = true
				return rec
			}(),
			in:  []string{msg01},
			out: []logmtest.Record{{Contains: msg01}, {Contains: msg02}},
			err: strconv.ErrRange,
		},
		"Any order: three logged, two expected": {
			rec: func() *logmtest.Recorder {
				rec := logmtest.NewRecorder()
				rec.ExpectAnyOrder = true
				return rec
			}(),
			in:  []string{msg01, msg02, msg03},
			out: []logmtest.Record{{Contains: msg03}, {Contains: msg02}},
			err: strconv.ErrRange,
		},
		"Unexpected: two logged, the one expected": {
			rec: func() *logmtest.Recorder {
				rec := logmtest.NewRecorder()
				rec.ExpectUnexpected = true
				return rec
			}(),
			in:  []string{msg01, msg02},
			out: []logmtest.Record{{Contains: msg01}},
		},
		"Unexpected: two logged, the second expected": {
			rec: func() *logmtest.Recorder {
				rec := logmtest.NewRecorder()
				rec.ExpectUnexpected = true
				return rec
			}(),
			in:  []string{msg01, msg02},
			out: []logmtest.Record{{Contains: msg02}},
		},
		"Unexpected: two logged, one expected": {
			rec: func() *logmtest.Recorder {
				rec := logmtest.NewRecorder()
				rec.ExpectUnexpected = true
				return rec
			}(),
			in:  []string{msg01, msg02},
			out: []logmtest.Record{{Contains: "oops"}},
			err: strconv.ErrRange,
		},
		"Unexpected in any order: two logged, one expected": {
			rec: func() *logmtest.Recorder {
				rec := logmtest.NewRecorder()
				rec.ExpectUnexpected = true
				rec.ExpectAnyOrder = true
				return rec
			}(),
			in:  []string{msg01, msg02},
			out: []logmtest.Record{{Contains: "oops"}},
			err: strconv.ErrRange,
		},
		"Unexpected in any order: three logged, two expected": {
			rec: func() *logmtest.Recorder {
				rec := logmtest.NewRecorder()
				rec.ExpectUnexpected = true
				rec.ExpectAnyOrder = true
				return rec
			}(),
			in:  []string{msg01, msg02, msg03},
			out: []logmtest.Record{{Contains: msg03}, {Contains: msg02}},
		},
	} {
		tt := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			log := logm.DefaultLogger("testing", tt.rec)
			for _, msg := range tt.in {
				log.Info(msg)
			}
			err := tt.rec.Expect(tt.out...)
			are.True(errors.Is(err, tt.err)) // mismatch error
		})
	}
}
