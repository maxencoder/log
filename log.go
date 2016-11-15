package log

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//log level, from low to high, more high means more serious
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	Ltime  = 1 << iota //time format "2006/01/02 15:04:05"
	Lfile              //file.go:123
	Llevel             //[Trace|Debug|Info...]
)

var LevelName [6]string = [6]string{"Trace", "Debug", "Info", "Warn", "Error", "Fatal"}

const TimeFormat = "2006/01/02 15:04:05"

type Logger struct {
	sync.Mutex

	level int
	flag  int

	handler Handler

	buf *bytes.Buffer
}

//new a logger with specified handler and flag
func New(handler Handler, flag int) *Logger {
	var l = new(Logger)

	l.level = LevelInfo
	l.handler = handler

	l.flag = flag

	l.buf = &bytes.Buffer{}

	return l
}

//new a default logger with specified handler and flag: Ltime|Lfile|Llevel
func NewDefault(handler Handler) *Logger {
	return New(handler, Ltime|Lfile|Llevel)
}

func newStdHandler() *StreamHandler {
	h, _ := NewStreamHandler(os.Stdout)
	return h
}

var std = NewDefault(newStdHandler())

//set log level, any log level less than it will not log
func (l *Logger) SetLevel(level int) {
	l.Lock()
	l.level = level
	l.Unlock()
}

//a low interface, maybe you can use it for your special log format
//but it may be not exported later......
func (l *Logger) Output(callDepth int, level int, format string, v ...interface{}) {
	l.Lock()
	defer l.Unlock()

	if l.level > level {
		return
	}

	buf := l.buf

	if l.flag&Ltime > 0 {
		now := time.Now().Format(TimeFormat)
		buf.WriteByte('[')
		buf.WriteString(now)
		buf.WriteString("] ")
	}

	if l.flag&Lfile > 0 {
		_, file, line, ok := runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
		} else {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}
		}

		buf.WriteString(file)
		buf.WriteByte(':')

		buf.WriteString(strconv.Itoa(line))
		buf.WriteByte(' ')
	}

	if l.flag&Llevel > 0 {
		buf.WriteByte('[')
		buf.WriteString(LevelName[level])
		buf.WriteString("] ")
	}

	s := fmt.Sprintf(format, v...)

	buf.WriteString(s)

	if s[len(s)-1] != '\n' {
		buf.WriteByte('\n')
	}

	log.Println(buf)
	buf.Reset()
}

//log with Trace level
func (l *Logger) Trace(format string, v ...interface{}) {
	l.Output(2, LevelTrace, format, v...)
}

//log with Debug level
func (l *Logger) Debug(format string, v ...interface{}) {
	l.Output(2, LevelDebug, format, v...)
}

//log with info level
func (l *Logger) Info(format string, v ...interface{}) {
	l.Output(2, LevelInfo, format, v...)
}

//log with warn level
func (l *Logger) Warn(format string, v ...interface{}) {
	l.Output(2, LevelWarn, format, v...)
}

//log with error level
func (l *Logger) Error(format string, v ...interface{}) {
	l.Output(2, LevelError, format, v...)
}

//log with fatal level
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Output(2, LevelFatal, format, v...)
}

func SetLevel(level int) {
	std.SetLevel(level)
}

func Trace(format string, v ...interface{}) {
	std.Output(2, LevelTrace, format, v...)
}

func Debug(format string, v ...interface{}) {
	std.Output(2, LevelDebug, format, v...)
}

func Info(format string, v ...interface{}) {
	std.Output(2, LevelInfo, format, v...)
}

func Warn(format string, v ...interface{}) {
	std.Output(2, LevelWarn, format, v...)
}

func Error(format string, v ...interface{}) {
	std.Output(2, LevelError, format, v...)
}

func Fatal(format string, v ...interface{}) {
	std.Output(2, LevelFatal, format, v...)
}
