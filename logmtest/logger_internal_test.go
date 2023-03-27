package logmtest

import (
	"testing"

	"github.com/matryer/is"
)

func TestNewRecorder(t *testing.T) {
	t.Parallel()
	var (
		rec = NewRecorder()
		are = is.New(t)
	)
	are.True(rec != nil)        // missing recorder
	are.Equal(0, rec.buf.Len()) // unexpected buffer
}
