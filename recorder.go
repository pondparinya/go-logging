package ecslog

//---------------------------------------------------------

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//---------------------------------------------------------

const API_TAG = "logger#"
const API_SID = "logger#service-id"
const API_TID = "logger#trace-id"

//---------------------------------------------------------

type Recorder struct {
	http.ResponseWriter

	// Logger *Logger

	RequestBody    []byte
	RequestHeader  []byte
	ResponseStatus int
	ResponseBody   []byte
	ResponseHeader []byte
}

//------------------

func (w *Recorder) WriteHeader(code int) {
	head, _ := json.Marshal(w.Header())
	w.ResponseHeader = head
	w.ResponseStatus = code
	w.ResponseWriter.WriteHeader(code)
}

//------------------

func (w *Recorder) Write(body []byte) (int, error) {
	w.ResponseBody = body
	return w.ResponseWriter.Write(body)
}

//------------------

func NewRecorder(w http.ResponseWriter, r *http.Request) *Recorder {

	w.Header().Set(API_SID, r.Header.Get(API_TID))
	w.Header().Set(API_TID, r.Header.Get(API_SID))

	h, _ := json.Marshal(r.Header)
	b, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return &Recorder{
		ResponseWriter: w,
		RequestHeader:  h,
		RequestBody:    b,
	}
}

//---------------------------------------------------------
// End-of-file
//---------------------------------------------------------
