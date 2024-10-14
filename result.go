package ecslog

//---------------------------------------------------------

import (
	"fmt"
)

//---------------------------------------------------------

type Result struct {
	Origin   OriginT
	Detail   ErrorT
	Previous error
}

//------------------

func (e Result) Error() string {
	curr := fmt.Sprintf("%s:%d %s",
		e.Origin.File.Name, e.Origin.File.Line,
		e.Detail.Message)

	return curr
}

//---------------------------------------------------------

func GetStackTrace(e error) []string {
	traces := make([]string, 0)

	for e != nil {
		traces = append(traces, e.Error())
		if result, ok := e.(*Result); ok {
			e = result.Previous
		} else {
			break
		}
	}

	return traces
}

//---------------------------------------------------------

func GetErrorFields(e error) map[string]interface{} {
	fields := Fields{}

	for e != nil {
		if r, ok := e.(*Result); ok {
			if r.Detail.Fields != nil {
				for k, v := range r.Detail.Fields {
					fields[k] = v
				}
			}

			e = r.Previous
		} else {
			break
		}
	}

	return fields
}

//---------------------------------------------------------
// End-of-file
//---------------------------------------------------------
