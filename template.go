package ecslog

//---------------------------------------------------------

import (
  "context"
  "encoding/json"
  "net/http"
)

//---------------------------------------------------------

type BodyT struct {
  Content string `json:"content,omitempty"`
}

//------------------

func (self BodyT) Clone() BodyT {
  return self
}

//---------------------------------------------------------

type RequestT struct {
  Referrer string `json:"referrer,omitempty"`
  Method   string `json:"method,omitempty"`
  Header   string `json:"header,omitempty"`
  Body     BodyT  `json:"body,omitempty"`
}

//------------------

func (self *RequestT) Clone() *RequestT {
  if self == nil { return nil }

  return &RequestT{
    Referrer:   self.Referrer,
    Method:     self.Method,
    Header:     self.Header,
    Body:       self.Body.Clone(),
  }
}

//---------------------------------------------------------

type ResponseT struct {
  Header     string `json:"header,omitempty"`
  Content    string `json:"content,omitempty"`
  StatusCode string `json:"status_code,omitempty"`
}

//------------------

func (self *ResponseT) Clone() *ResponseT {
  if self == nil { return nil }

  return &ResponseT{
    Header:     self.Header,
    Content:    self.Content,
    StatusCode: self.StatusCode,
  }
}

//---------------------------------------------------------

type HttpT struct {
  Version  string     `json:"version,omitempty"`
  Request  *RequestT  `json:"request,omitempty"`
  Response *ResponseT `json:"response,omitempty"`
}

//------------------

func (self *HttpT) Clone() *HttpT {
  if self == nil { return nil }

  return &HttpT{
    Version:  self.Version,
    Request:  self.Request.Clone(),
    Response: self.Response.Clone(),
  }
}

//---------------------------------------------------------

type UrlT struct {
  Full string `json:"full,omitempty"` // full-URL
}

//------------------

func (self *UrlT) Clone() *UrlT {
  if self == nil { return nil }

  return &UrlT{
    Full: self.Full,
  }
}

//---------------------------------------------------------

type ServiceT struct {
  Id string `json:"id"` // process
}

//------------------

func (self ServiceT) Clone() ServiceT {
   return self
}

//---------------------------------------------------------

type Fields map[string]interface{}

//------------------

func (self Fields) String() string {
  buff, _ := json.Marshal(self)
  return string(buff)
}

//------------------

func (self Fields) Clone(newFieldss ...Fields) Fields {
  tmp := Fields{}
  for k, v := range self {
    tmp[k] = v
  }

  if newFieldss != nil {
    for _, newFields := range newFieldss {
      for k, v := range newFields {
        tmp[k] = v
      }
    }
  }

  return tmp
}

//---------------------------------------------------------

type LogT struct {
  Level  string   `json:"level"`
  Origin *OriginT `json:"origin,omitempty"`
  Fields Fields   `json:"tags,omitempty"`
}

func (self LogT) Clone() LogT {
  return self
}

//---------------------------------------------------------

type ErrorT struct {
  // Code   *string   `json:"code,omitempty"` // code
  Fields     Fields   `json:"fields,omitempty"`
  Message    string   `json:"message,omitempty"`
  StackTrace []string `json:"stack_trace,omitempty"`
}

func (self ErrorT) Clone() ErrorT {
  return self
}

//---------------------------------------------------------

type Template struct {
  Event   *EventT  `json:"event,omitempty"`
  Service ServiceT `json:"service,omitempty"`
  Log     LogT     `json:"log,omitempty"`
  Trace   TraceT   `json:"trace,omitempty"`
  Error   ErrorT   `json:"error,omitempty"`
  Http    *HttpT   `json:"http,omitempty"`
  Url     *UrlT    `json:"url,omitempty"`
}

//------------------

func (self *Template) Clone() *Template {
  if self == nil { return nil }

  return &Template{
    Event:   self.Event.Clone(),
    Service: self.Service.Clone(),
    Log:     self.Log.Clone(),
    Error:   self.Error.Clone(),
    Http:    self.Http.Clone(),
    Url:     self.Url.Clone(),
  }
}

//------------------

func (self Template) String() string {
  buff, _ := json.Marshal(self)
  return string(buff)
}

//------------------

func (self *Template) SetHttp(url, version string) *Template {
  if self == nil { return nil }

  if self.Http == nil {
    self.Http = &HttpT{}
  }
  self.Http.Version = version

  if self.Url == nil {
    self.Url = &UrlT{}
  }
  self.Url.Full = url

  return self
}

//------------------

func (self *Template) SetRequest(referrer, method, content string) *Template {
  if self == nil { return nil }

  if self.Http == nil {
    self.Http = &HttpT{}
  }

  if self.Http.Request == nil {
    self.Http.Request = &RequestT{}
  }

  self.Http.Request.Body.Content = content
  self.Http.Request.Referrer = referrer
  self.Http.Request.Method = method
  return self
}

func (self *Template) SetRequestHeader(header string) *Template {
  if self == nil {
    return nil
  }

  if self.Http == nil {
    self.Http = &HttpT{}
  }

  if self.Http.Request == nil {
    self.Http.Request = &RequestT{}
  }

  self.Http.Request.Header = header
  return self
}

func (self *Template) ClearRequest() *Template {
  self.Http.Request = nil
  return self
}

//------------------

// SetResponse - Set http status and response body
func (self *Template) SetResponse(statusCode, content string) *Template {
  if self == nil { return nil }

  if self.Http == nil {
    self.Http = &HttpT{}
  }

  if self.Http.Response == nil {
    self.Http.Response = &ResponseT{}
  }

  self.Http.Response.StatusCode = statusCode
  self.Http.Response.Content = content
  return self
}

// SetResponseHeader - Set header response
func (self *Template) SetResponseHeader(header string) *Template {
  if self == nil {
    return nil
  }

  if self.Http == nil {
    self.Http = &HttpT{}
  }

  if self.Http.Response == nil {
    self.Http.Response = &ResponseT{}
  }

  self.Http.Response.Header = header
  return self
}

func (self *Template) ClearResponse() *Template {
  self.Http.Response = nil
  return self
}

//------------------

func (self *Template) SetError(e error) *Template {
  if self == nil { return nil }

  if e == nil {
    self.Error.StackTrace = nil
    self.Error.Fields = nil
    self.Error.Message = ""
  } else {
    stackTrace := GetStackTrace(e)
    self.Error.StackTrace = stackTrace
    self.Error.Fields = GetErrorFields(e)

    if self.Error.Message == "" {
      if len(stackTrace) > 0 {
        self.Error.Message = stackTrace[0]
      }
    }
  }

  return self
}

//------------------

func (self *Template) SetStatus(level, text string, origin *OriginT) *Template {
  if self == nil { return nil }

  self.Error.Message = text
  self.Log.Level = level
  if origin == nil {
    self.Log.Origin = NewOrigin(1)
  } else {
    self.Log.Origin = origin
  }

  return self
}

//------------------

func (self *Template) ResetField(key string) *Template {
  if self == nil { return nil }

  if self.Log.Fields != nil {
    delete(self.Log.Fields, key)
  }

  if key == API_SID {
    self.Service.Id = ""
  } else if key == API_TID {
    self.Trace.Id = ""
  }

  return self
}

//------------------

func (self *Template) ResetFields() *Template {
  if self == nil { return nil }

  self.Log.Fields = nil
  self.Service.Id = ""
  self.Trace.Id = ""
  return self
}

//------------------

func (self *Template) SetField(key string, value interface{}) *Template {
  if self == nil { return nil }

  if self.Log.Fields == nil {
    self.Log.Fields = Fields{}
  }

  self.Log.Fields[key] = value
  /*
    if key == API_SID {
      self.Service.Id = value
    } else if key == API_TID {
      self.Trace.Id = value
    }
    //*/

  return self
}

//------------------

func (self *Template) ServiceId() string {
  if self == nil { return "" }

  return self.Service.Id
}

func (self *Template) SetServiceId(value string) *Template {
  if self == nil { return nil }

  // self.SetField(API_SID, value)
  self.Service.Id = value
  return self
}

//------------------

func (self *Template) TraceId() string {
  if self == nil { return "" }

  return self.Trace.Id
}

func (self *Template) SetTraceId(value string) *Template {
  if self == nil { return nil }

  // self.SetField(API_TID, value)
  self.Trace.Id = value
  return self
}

//------------------

func (self *Template) Continue() *Template {
  if self == nil { return nil }

  if self.Event == nil {
    self.Event = &EventT{}
  }

  self.Event.Continue()
  return self
}

//------------------

func (self *Template) Start() *Template {
  if self == nil { return nil }

  if self.Event == nil {
    self.Event = &EventT{}
  }

  self.Event.Start()
  return self
}

//------------------

func (self *Template) End() *Template {
  if self.Event == nil {
    self.Event = &EventT{}
  }

  self.Event.End()
  return self
}

//------------------

func (self *Template) Reset() *Template {
  if self == nil { return nil }

  if self.Event == nil {
    self.Event = &EventT{}
  }

  self.Event.Reset()
  return self
}

//---------------------------------------------------------

func NewTemplate(serviceId string) *Template {
  return &Template{
    Service: ServiceT{
      Id: serviceId,
    },
  }
}

//---------------------------------------------------------

func NewTemplateWithTraceId(serviceId, traceId string) *Template {
  return &Template{
    Service: ServiceT{
      Id: serviceId,
    },

    Trace: TraceT{
      Id: traceId,
    },
  }
}

//---------------------------------------------------------

func NewTemplateFromRequest(r *http.Request) *Template {
  if r == nil {
    return &Template{
      Service: ServiceT{Id: ServiceId()},
    }
  }

  serviceId, ok := ServiceId(), false
  if serviceId == "" {
    serviceId = r.Header.Get(API_SID)
    if serviceId == "" {
      if serviceId, ok = r.Context().Value(API_SID).(string); ok && serviceId != "" {
        c := context.WithValue(r.Context(), API_SID, serviceId)
        *r = *r.WithContext(c)
      }
    }
  }

  traceId := r.Header.Get(API_TID)
  if traceId == "" {
    traceId = NewTraceId(serviceId)
    c := context.WithValue(r.Context(), API_TID, traceId)
    *r = *r.WithContext(c)
    r.Header.Set(API_TID, traceId)
  }

  return &Template{
    Http:    &HttpT{Version: r.Proto},
    Service: ServiceT{Id: serviceId},
    Trace:   TraceT{Id: traceId},
    Url:     &UrlT{Full: r.RequestURI},
  }
}

//-----------------------------------------------

func NewTemplateFromContext(c context.Context) *Template {
  if c == nil {
    return &Template{
      Service: ServiceT{Id: ServiceId()},
    }
  }

  var (
    serviceId = ServiceId()
    ok        = false
    traceId   = ""
  )

  if serviceId == "" {
    if serviceId, ok = c.Value(API_SID).(string); !ok {
      serviceId = ""
    }
  }

  if traceId, ok = c.Value(API_TID).(string); !ok {
    traceId = ""
  }

  return &Template{
    Service: ServiceT{Id: serviceId},
    Trace:   TraceT{Id: traceId},
  }
}

//---------------------------------------------------------
// End-of-file
//---------------------------------------------------------

