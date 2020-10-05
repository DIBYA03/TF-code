package log

import (
	"io"
	"os"
	"runtime"
	"runtime/debug"

	lr "github.com/sirupsen/logrus"
)

const (
	stackKey          = "stack"
	reportLocationKey = "reportLocation"
)

//Logger interface, a subset of logrus
type Logger interface {
	Debug(string)
	DebugD(string, Fields)

	Info(string)
	InfoD(string, Fields)

	Warn(string)
	WarnD(string, Fields)

	Error(string)
	ErrorD(string, Fields)

	Panic(string)
	PanicD(string, Fields)
}

type logger struct {
	l *lr.Logger
}

//Fields and type lr.Fields are a map[string]interface
type Fields lr.Fields

//NewLogger returns a logger interface
func NewLogger(opts ...Option) Logger {
	//Default args
	la := &loggerOpts{
		output: os.Stdout,
		format: &lr.JSONFormatter{},
	}

	for _, opt := range opts {
		opt(la)
	}

	l := lr.New()

	l.SetOutput(la.output)
	l.SetFormatter(la.format)

	if os.Getenv("API_ENV") == "prod" {
		l.SetLevel(lr.InfoLevel)
	} else {
		l.SetLevel(lr.DebugLevel)
	}

	return logger{
		l: l,
	}
}

type loggerOpts struct {
	output io.Writer
	format lr.Formatter
}

//Option optional funcs passed into NewLogger
type Option func(*loggerOpts)

//SetOutput is optionally passed into NewLogger it's used to set the log output
func SetOutput(i io.Writer) Option {
	return func(opts *loggerOpts) {
		opts.output = i
	}
}

//SetFormat is optionally passed into NewLogger it's used to set the log format
func SetFormat(f lr.Formatter) Option {
	return func(opts *loggerOpts) {
		opts.format = f
	}
}

//Debug will log debug level logs, will not log in production
func (l logger) Debug(s string) {
	f := l.appendReportLocation(Fields{})
	l.l.WithFields(f.format()).Debug(s)
}

//DebugD will log debug level logs with Fields, will not log in production
func (l logger) DebugD(s string, f Fields) {
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Debug(s)
}

//Info will log info level logs, will appear in production logs
func (l logger) Info(s string) {
	f := l.appendReportLocation(Fields{})
	l.l.WithFields(f.format()).Info(s)
}

//Info will log info level logs with Fields, will appear in production logs
func (l logger) InfoD(s string, f Fields) {
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Info(s)
}

//Warn will log warn level logs and a stacktrace, will appear in production logs
func (l logger) Warn(s string) {
	f := l.appendStack(Fields{})
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Warn(s)
}

//WarnD will log want level logs with Fields and a stacktrace, will appear in production logs
func (l logger) WarnD(s string, f Fields) {
	f = l.appendStack(f)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Warn(s)
}

//Error will log error level logs and a stacktrace, will appear in production logs
func (l logger) Error(s string) {
	f := l.appendStack(Fields{})
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Error(s)
}

//ErrorD will log error level logs with Fields and a stacktrace, will appear in production logs
func (l logger) ErrorD(s string, f Fields) {
	f = l.appendStack(f)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Error(s)
}

//Panic will log panic level logs and a stacktrace, will appear in production logs
func (l logger) Panic(s string) {
	f := l.appendStack(Fields{})
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Panic(s)
}

//PanicD will log panic level logs with Fields and a stacktrace, will appear in production logs
func (l logger) PanicD(s string, f Fields) {
	f = l.appendStack(f)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Panic(s)
}

func (l logger) appendStack(f Fields) Fields {
	f[stackKey] = string(debug.Stack())

	return f
}

func (l logger) appendReportLocation(f Fields) Fields {
	pc, fn, line, _ := runtime.Caller(2)

	f[reportLocationKey] = map[string]interface{}{
		"filePath":     fn,
		"line":         line,
		"functionName": runtime.FuncForPC(pc).Name(),
	}

	return f
}

//Format to logrus formatted fields
//We could also put any log data sanitisazation in here
func (f Fields) format() lr.Fields {
	return lr.Fields(f)
}
