package logm

import "github.com/natefinch/lumberjack"

// maxMoFileSize is the maximum size in megabytes of the log file before it gets rotated.
const maxMoFileSize = 100

// NewFile returns a file with the given name and time localization setup.
// This file will be automatically rotated if it its size exceeds the 100 Mo.
// The newly created file will have this format `name-timestamp.ext`.
// It is the filename without the extension,
// timestamp is the time at which the log was rotated formatted with the time.Time
// format of `2006-01-02T15-04-05.000` and the extension is the original extension.
// For example, with `/data/log/server.log` as name, a backup created at 6:30pm on Nov 11 2016
// would use the filename `/data/log/server-2016-11-04T18-30-00.000.log`
func NewFile(name string, localTime bool) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   name,
		MaxSize:    maxMoFileSize,
		MaxAge:     0,
		MaxBackups: 0,
		LocalTime:  localTime,
		Compress:   true,
	}
}
