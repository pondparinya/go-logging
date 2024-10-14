package ecslog

//---------------------------------------------------------

import (
	"fmt"
	"net/http"
)

//---------------------------------------------------------

var logger Logger

//------------------

func ServiceId() string {
	return logger.ServiceId()
}

func SetServiceId(value string) {
	logger.SetServiceId(value)
}

//------------------

func TraceId() string {
	return logger.TraceId()
}

func SetTraceId(value string) {
	logger.SetTraceId(value)
}

//------------------

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicln(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

//------------------

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

//------------------

func Trace(args ...interface{}) *Logger {
	return logger.Trace(args...)
}

func Traceln(args ...interface{}) *Logger {
	return logger.Traceln(args...)
}

func Tracef(format string, args ...interface{}) *Logger {
	return logger.Tracef(format, args...)
}

//------------------

func Debug(args ...interface{}) *Logger {
	return logger.Debug(args...)
}

func Debugln(args ...interface{}) *Logger {
	return logger.Debugln(args...)
}

func Debugf(format string, args ...interface{}) *Logger {
	return logger.Debugf(format, args...)
}

//------------------

func Info(args ...interface{}) *Logger {
	return logger.Info(args...)
}

func Infoln(args ...interface{}) *Logger {
	return logger.Infoln(args...)
}

func Infof(format string, args ...interface{}) *Logger {
	return logger.Infof(format, args...)
}

//------------------

func Warn(args ...interface{}) *Logger {
	return logger.Warn(args...)
}

func Warnln(args ...interface{}) *Logger {
	return logger.Warnln(args...)
}

func Warnf(format string, args ...interface{}) *Logger {
	return logger.Warnf(format, args...)
}

//------------------

func Error(args ...interface{}) *Logger {
	return logger.Error(args...)
}

func Errorln(args ...interface{}) *Logger {
	return logger.Errorln(args...)
}

func Errorf(format string, args ...interface{}) *Logger {
	return logger.Errorf(format, args...)
}

//------------------

func PrintStackTrace(e error, args ...interface{}) *Logger {
	return logger.PrintStackTrace(e, args...)
}

func PrintStackTraceln(e error, args ...interface{}) *Logger {
	return logger.PrintStackTraceln(e, args...)
}

func PrintStackTracef(e error, format string, args ...interface{}) *Logger {
	return logger.PrintStackTracef(e, format, args...)
}

//---------------------------------------------------------

func ResetField(key string) *Logger {
	return logger.ResetField(key)
}

//---------------------------------------------------------

func ResetFields() *Logger {
	return logger.ResetFields()
}

//---------------------------------------------------------

func SetField(key, value string) *Logger {
	return logger.SetField(key, value)
}

func SetFields(fields Fields) *Logger {
	return logger.SetFields(fields)
}

func WithFields(fields Fields) *Logger {
	return logger.WithFields(fields)
}

//---------------------------------------------------------

func NewResult(args ...interface{}) *Result {
	return &Result{
		Origin: *NewOrigin(2),
		Detail: ErrorT{
			Fields:  logger.Template.Log.Fields,
			Message: fmt.Sprint(args...),
		},
	}
}

func NewResultf(format string, args ...interface{}) *Result {
	return &Result{
		Origin: *NewOrigin(2),
		Detail: ErrorT{
			Fields:  logger.Template.Log.Fields,
			Message: fmt.Sprintf(format, args...),
		},
	}
}

//---------------------------------------------------------

func WrapResult(previous error, args ...interface{}) *Result {
	return &Result{
		Origin: *NewOrigin(2),
		Detail: ErrorT{
			Fields:  logger.Template.Log.Fields,
			Message: fmt.Sprint(args...),
		},
		Previous: previous,
	}
}

func WrapResultf(previous error, format string, args ...interface{}) *Result {
	return &Result{
		Origin: *NewOrigin(2),
		Detail: ErrorT{
			Fields:  logger.Template.Log.Fields,
			Message: fmt.Sprintf(format, args...),
		},
		Previous: previous,
	}
}

//---------------------------------------------------------

func SetupWithTraceId(serviceId, traceId string) *Logger {
	logger.Template = NewTemplateWithTraceId(serviceId, traceId)
	return &logger
}

//------------------

func Setup(serviceId string) *Logger {
	logger.Template = NewTemplate(serviceId)
	return &logger
}

//---------------------------------------------------------

func (log *Logger) ForwardTraceId(req *http.Request) {
	req.Header.Set(API_TID, log.TraceId())
}

//---------------------------------------------------------

//---------------------------------------------------------
// End-of-file
//---------------------------------------------------------
