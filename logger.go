package ecslog

//---------------------------------------------------------

import (
  "context"
  "fmt"
  "math/rand"
  "net/http"
  "os"
  "time"
)

//---------------------------------------------------------

type Logger struct {
  Template *Template
}

//------------------

func (self *Logger) ServiceId() string {
  return self.Template.ServiceId()
}

func (self *Logger) SetServiceId(value string) *Logger {
  self.Template.SetServiceId(value)
  return self
}

//------------------

func (self *Logger) TraceId() string {
  return self.Template.TraceId()
}

func (self *Logger) SetTraceId(value string) *Logger {
  self.Template.SetTraceId(value)
  return self
}

//------------------

func (self *Logger) format(level string, args ...interface{}) string {
  orgn := NewOrigin(5)
  buff := CensorValues(args...)
  text := fmt.Sprint(buff...)
  return self.Template.SetStatus(level, text, orgn).String()
}

func (self *Logger) formatf(level string, format string, args ...interface{}) string {
  orgn := NewOrigin(5)
  buff := CensorValues(args...)
  text := fmt.Sprintf(format, buff...)
  return self.Template.SetStatus(level, text, orgn).String()
}

//------------------

func (self *Logger) Panic(args ...interface{}) {
  text := self.format("PANIC", args...)
  panic(text)
}

func (self *Logger) Panicln(args ...interface{}) {
  text := self.format("PANIC", args...)
  panic(text)
}

func (self *Logger) Panicf(format string, args ...interface{}) {
  text := self.formatf("PANIC", format, args...)
  panic(text)
}

//------------------

func (self *Logger) print(level string, args ...interface{}) *Logger {
  if self.Template == nil {
    self.Template = NewTemplate("")
  }
  self.Template.Start()
  text := self.format(level, args...)
  fmt.Printf("%s\n\n", text)
  return self
}

func (self *Logger) printf(level string, format string, args ...interface{}) *Logger {
  if self.Template == nil {
    self.Template = NewTemplate("")
  }
  self.Template.Start()
  text := self.formatf(level, format, args...)
  fmt.Printf("%s\n\n", text)
  return self
}

//------------------

func (self *Logger) print2(level string, args ...interface{}) *Logger {
  self.Template.Start().End()
  text := self.format(level, args...)
  fmt.Printf("%s\n\n", text)
  self.Template.Continue()
  return self
}

func (self *Logger) printf2(level string, format string, args ...interface{}) *Logger {
  self.Template.Start().End()
  text := self.formatf(level, format, args...)
  fmt.Printf("%s\n\n", text)
  self.Template.Continue()
  return self
}

//------------------

func (self *Logger) Fatal(args ...interface{}) {
  self.print("FATAL", args...)
  os.Exit(1)
}

func (self *Logger) Fatalln(args ...interface{}) {
  self.print("FATAL", args...)
  os.Exit(1)
}

func (self *Logger) Fatalf(format string, args ...interface{}) {
  self.printf("FATAL", format, args...)
  os.Exit(1)
}

//------------------

func (self *Logger) Trace(args ...interface{}) *Logger {
  return self.print2("TRACE", args...)
}

func (self *Logger) Traceln(args ...interface{}) *Logger {
  return self.print2("TRACE", args...)
}

func (self *Logger) Tracef(format string, args ...interface{}) *Logger {
  return self.printf2("TRACE", format, args...)
}

//------------------

func (self *Logger) Debug(args ...interface{}) *Logger {
  return self.print("DEBUG", args...)
}

func (self *Logger) Debugln(args ...interface{}) *Logger {
  return self.print("DEBUG", args...)
}

func (self *Logger) Debugf(format string, args ...interface{}) *Logger {
  return self.printf("DEBUG", format, args...)
}

//------------------

func (self *Logger) Info(args ...interface{}) *Logger {
  return self.print("INFO", args...)
}

func (self *Logger) Infoln(args ...interface{}) *Logger {
  return self.print("INFO", args...)
}

func (self *Logger) Infof(format string, args ...interface{}) *Logger {
  return self.printf("INFO", format, args...)
}

//------------------

func (self *Logger) Warn(args ...interface{}) *Logger {
  return self.print("WARN", args...)
}

func (self *Logger) Warnln(args ...interface{}) *Logger {
  return self.print("WARN", args...)
}

func (self *Logger) Warnf(format string, args ...interface{}) *Logger {
  return self.printf("WARN", format, args...)
}

//------------------

func (self *Logger) Error(args ...interface{}) *Logger {
  return self.print("ERROR", args...)
}

func (self *Logger) Errorln(args ...interface{}) *Logger {
  return self.print("ERROR", args...)
}

func (self *Logger) Errorf(format string, args ...interface{}) *Logger {
  return self.printf("ERROR", format, args...)
}

//------------------

func (self *Logger) Printf(format string, args ...interface{}) {
  fmt.Printf(fmt.Sprintf("%s\n\n", format), args...)
}

//------------------

func (self *Logger) PrintStackTrace(e error, args ...interface{}) *Logger {
  self.Template.SetError(e)
  self.print("ERROR", args...)
  self.Template.SetError(nil)
  return self
}

func (self *Logger) PrintStackTraceln(e error, args ...interface{}) *Logger {
  self.Template.SetError(e)
  self.print("ERROR", args...)
  self.Template.SetError(nil)
  return self
}

func (self *Logger) PrintStackTracef(e error, format string, args ...interface{}) *Logger {
  self.Template.SetError(e)
  self.printf("ERROR", format, args...)
  self.Template.SetError(nil)
  return self
}

//---------------------------------------------------------

func (self Logger) ForwardWith(traceId string) *Logger {
  self.SetTraceId(traceId)
  return &self
}

//------------------

func (self Logger) Forward() *Logger {
  now := time.Now()
  prefix := now.Format("20060102150405")
  suffix := (now.Second() + rand.Int()) % 1023
  return self.ForwardWith(fmt.Sprintf("%s%03d", prefix, suffix))
}

//---------------------------------------------------------

func (self *Logger) ResetField(key string) *Logger {
  if self.Template != nil {
    self.Template.ResetField(key)
  }

  return self
}

//---------------------------------------------------------

func (self *Logger) ResetFields() *Logger {
  if self.Template != nil {
    self.Template.ResetFields()
  }

  return self
}

//---------------------------------------------------------

func (self *Logger) SetField(key string, value interface{}) *Logger {
  if self.Template != nil {
    self.Template.SetField(key, value)
  }

  return self
}

func (self *Logger) SetFields(fields Fields) *Logger {
  if fields != nil {
    for key, value := range fields {
      self.SetField(key, value)
    }
  }

  return self
}

func (self *Logger) WithFields(fields Fields) *Logger {
  return self.Clone().SetFields(fields)
}

//---------------------------------------------------------

func (self *Logger) NewResult(args ...interface{}) *Result {
  return &Result{
    Origin: *NewOrigin(2),
    Detail: ErrorT{
      Fields:  self.Template.Log.Fields,
      Message: fmt.Sprint(args...),
    },
  }
}

func (self *Logger) NewResultf(format string, args ...interface{}) *Result {
  return &Result{
    Origin: *NewOrigin(2),
    Detail: ErrorT{
      Fields:  self.Template.Log.Fields,
      Message: fmt.Sprintf(format, args...),
    },
  }
}

//---------------------------------------------------------

func (self *Logger) WrapResult(previous error, args ...interface{}) *Result {
  return &Result{
    Origin: *NewOrigin(2),
    Detail: ErrorT{
      Fields:  self.Template.Log.Fields,
      Message: fmt.Sprint(args...),
    },
    Previous: previous,
  }
}

func (self *Logger) WrapResultf(previous error, format string, args ...interface{}) *Result {
  return &Result{
    Origin: *NewOrigin(2),
    Detail: ErrorT{
      Fields:  self.Template.Log.Fields,
      Message: fmt.Sprintf(format, args...),
    },
    Previous: previous,
  }
}

//---------------------------------------------------------

func (self *Logger) Clone() *Logger {
  return &Logger{Template: self.Template.Clone()}
}

//---------------------------------------------------------

func (self *Logger) Fork() *Logger {
  return NewLogger(self.ServiceId())
}

//---------------------------------------------------------

func (self *Logger) SetRequest(referrer, method, content string) *Logger {
  if self.Template != nil {
    self.Template.SetRequest(referrer, method, content)
  }

  return self
}

func (self *Logger) SetRequestHeader(header string) *Logger {
  if self.Template != nil {
    self.Template.SetRequestHeader(header)
  }

  return self
}

func (self *Logger) ClearRequest() *Logger {
  if self.Template != nil {
    self.Template.ClearRequest()
  }
  return self
}

//---------------------------------------------------------

func (self *Logger) SetResponse(statusCode, content string) *Logger {
  if self.Template != nil {
    self.Template.SetResponse(statusCode, content)
  }

  return self
}

func (self *Logger) SetResponseHeader(header string) *Logger {
  if self.Template != nil {
    self.Template.SetResponseHeader(header)
  }

  return self
}

func (self *Logger) ClearResponse() *Logger {
  if self.Template != nil {
    self.Template.ClearResponse()
  }
  return self
}

//---------------------------------------------------------

func NewLogger(serviceId string) *Logger {
  template := NewTemplate(serviceId)
  return &Logger{Template: template}
}

//------------------

func GetLogger(v interface{}) *Logger {
  if r, ok := v.(*http.Request); ok {
    return &Logger{Template: NewTemplateFromRequest(r)}
  }
  if c, ok := v.(context.Context); ok {
    return &Logger{Template: NewTemplateFromContext(c)}
  }
  return &Logger{Template: NewTemplate(ServiceId())}
}

//-----------------------------

func GetLogEntry(r *http.Request) *Logger {
  newLogger := logger
  if r != nil {
    if rid := r.Context().Value("uuid").(string); rid != "" {
      return newLogger.WithFields(Fields{"request_id": rid})
    }
  }
  return &newLogger
}

//---------------------------------------------------------
// End-of-file
//---------------------------------------------------------
