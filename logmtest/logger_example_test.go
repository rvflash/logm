package logmtest_test

import (
	"fmt"

	"github.com/rvflash/logm"
	"github.com/rvflash/logm/logmtest"
)

func ExampleNewRecorder() {
	var (
		rec = logmtest.NewRecorder()
		log = logm.DefaultLogger("testing", rec)
	)
	log.Info("hello")
	log.Warn("beautiful")
	log.Error("world")

	err := rec.Expect(
		logmtest.Record{Contains: "hello"},
		logmtest.Record{Contains: "beau"},
		logmtest.Record{Contains: "world"},
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Output:
}
