package ecslog

//---------------------------------------------------------

import (
  "runtime"
  "strings"
)

//---------------------------------------------------------

type FileT struct {
  Name string `json:"name"` // fileName
  Line int    `json:"line"` // fileLine
}

//------------------

func (self FileT) Clone() FileT {
  return FileT{
    Name: self.Name,
    Line: self.Line,
  }
}

//---------------------------------------------------------

type OriginT struct {
  Function string `json:"function"` // function
  File     FileT  `json:"file,omitempty"`
}

//------------------

func (self OriginT) Clone() OriginT {
  return OriginT{
    Function: self.Function,
    File:     self.File.Clone(),
  }
}

//---------------------------------------------------------

func NewOrigin(skip int) *OriginT {
  origin := OriginT{}

  if pc, file, line, ok := runtime.Caller(skip); ok {
    if i := strings.Index(file, "/src/"); i != -1 {
      file = file[i+5:]
    }

    origin.File.Name = file
    origin.File.Line = line

    if fp := runtime.FuncForPC(pc); fp != nil {
      origin.Function = fp.Name()
    }
  }

  return &origin
}

//---------------------------------------------------------
// End-of-file
//---------------------------------------------------------

