package ecslog

//---------------------------------------------------------

import (
  "fmt"
  "sync/atomic"
  "time"
)

//---------------------------------------------------------

type TraceT struct {
  Id string `json:"id,omitempty"` // traceId
}

//------------------

func (self TraceT) Clone() TraceT {
  return self
}

//---------------------------------------------------------

var logCount uint64

func NewTraceId(prefix string) string {
  TraceTs := time.Now().Format("20060102030405")
  counter := atomic.AddUint64(&logCount, 1)
  return fmt.Sprintf("%s-%s-%07d", prefix, TraceTs, counter%1000000)
}

//---------------------------------------------------------

func NewTrace(id string) TraceT {
  return TraceT{Id: id}
}

//---------------------------------------------------------

func NewTraceWithPrefix(prefix string) TraceT {
  traceId := NewTraceId(prefix)
  return NewTrace(traceId)
}

//---------------------------------------------------------
// End-of-file
//---------------------------------------------------------

