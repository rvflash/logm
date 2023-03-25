package logm

import "github.com/natefinch/lumberjack"

// maxMoFileSize is the maximum size in megabytes of the log file before it gets rotated.
const maxMoFileSize = 100

// File is a log file.
type File = lumberjack.Logger

// NewFile returns a file with this name.
// This file will be automatically rotated if it its size exceeds the 100 Mo.
// The newly created file will have this format `name-timestamp.ext` and will be compressed.
// timestamp is the time at which the log was rotated formatted with the time.Time
// format of `2006-01-02T15-04-05.000` and the extension is the original extension.
// For example, with `/data/log/server.log` as name, a backup created at 6:30pm on Nov 11 2016
// would use the filename `/data/log/server-2016-11-04T18-30-00.000.log`
// Each record will have a local time.
func NewFile(name string) *File {
	return &File{
		Filename:   name,
		MaxSize:    maxMoFileSize,
		MaxAge:     0,
		MaxBackups: 0,
		LocalTime:  true,
		Compress:   true,
	}
}
