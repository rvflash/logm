package logm_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/logm"
)

const filename = "file.log"

func TestNewFile(t *testing.T) {
	t.Parallel()
	var (
		f   = logm.NewFile(filename)
		are = is.New(t)
	)
	are.Equal(filename, f.Filename) // mismatch filename
	are.Equal(f.MaxSize, 100)       // max size mismatch
	are.Equal(f.MaxAge, 0)          // max age mismatch
	are.Equal(f.MaxBackups, 0)      // max backups mismatch
	are.Equal(f.LocalTime, true)    // local time mismatch
	are.Equal(f.Compress, true)     // compress mismatch
}
