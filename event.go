package ecslog

//---------------------------------------------------------

import (
  "time"
)

//---------------------------------------------------------

type EventT struct {
  TimeStart *time.Time     `json:"start,omitempty"`
  TimeEnd   *time.Time     `json:"end,omitempty"`
  Duration  *time.Duration `json:"duration,omitempty"`
}

//------------------

func (self *EventT) Clone() *EventT {
  if self == nil { return nil }

  return &EventT{
    TimeStart: self.TimeStart,
    TimeEnd:   self.TimeEnd,
    Duration:  self.Duration,
  }
}

//------------------

func (self *EventT) Continue() *EventT {
  if self == nil { return nil }

  self.TimeStart = self.TimeEnd
  self.TimeEnd = nil
  self.Duration = nil
  return self
}

//------------------

func (self *EventT) Reset() *EventT {
  if self == nil { return nil }

  self.Duration = nil
  self.TimeEnd = nil
  self.TimeStart = nil
  return self
}

//------------------

func (self *EventT) Start() *EventT {
  if self == nil { return nil }

  if self.TimeStart == nil {
    dt := time.Now()
    self.TimeStart = &dt
  }

  return self
}

//------------------

func (self *EventT) End() *EventT {
  if self == nil { return nil }

  if self.TimeStart != nil &&
    self.TimeEnd == nil {
    dt := time.Now()
    self.TimeEnd = &dt

    diff := self.TimeEnd.Sub(*self.TimeStart)
    self.Duration = &diff
  }

  return self
}

//---------------------------------------------------------

func NewEvent() *EventT {
  return &EventT{}
}

//---------------------------------------------------------
// TimeEnd-of-file
//---------------------------------------------------------
