package ecslog

//---------------------------------------------------------

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

//---------------------------------------------------------

func Handler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rec := NewRecorder(w, r)
	log := GetLogger(r)

	if r.URL.Path == "/status" {
		next(rec, r)
		return
	}
	defer func() {
		log.SetResponse(
			strconv.FormatInt(int64(rec.ResponseStatus), 10),
			string(rec.ResponseBody))

		if rcv := recover(); rcv != nil {
			log.Panicln(rcv)
			panic(rcv)
		}

		log.Traceln("DONE")
	}()

	{
		log.SetRequest(
			r.Header.Get("Referer"),
			r.Method,
			string(rec.RequestBody))

		log.Traceln("START")

		log.ClearRequest()
	}

	next(rec, r)
}

//------------------

func LogHandler(next http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Handler(w, r, next)
	}
}

//---------------------------------------------------------

func LogMiddleWare() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger := logger
			logFields := Fields{}
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			logFields["url"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
			if rid := r.Context().Value("uuid"); rid != "" {
				logFields["request_id"] = rid
			}
			logFields["http_scheme"] = scheme
			logFields["http_proto"] = r.Proto
			logFields["http_method"] = r.Method
			logFields["remote_addr"] = r.RemoteAddr
			logFields["user_agent"] = r.UserAgent()

			// try to get body from request
			var bodyReq interface{}
			content := r.Header.Get("Content-Type")
			if !strings.Contains(content, "image") && !strings.Contains(content, "multi") && r.Body != nil {
				buf, err := ioutil.ReadAll(r.Body)
				if err == nil {
					rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
					rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

					if err := json.NewDecoder(rdr1).Decode(bodyReq); err != nil {
						bodyReq = string(buf)
					}
					r.Body = rdr2
				}
			}
			logFields["request"] = bodyReq
			logger.WithFields(logFields).Info("request started")
			rec := NewRecorder(w, r)
			t1 := time.Now()
			defer func() {
				tLog := logger.WithFields(logFields).WithFields(Fields{
					"resp_status":       rec.ResponseStatus,
					"resp_bytes_length": len(rec.ResponseBody),
					"resp_elapsed_ms":   float64(time.Since(t1).Nanoseconds()) / 1000000.0,
				})

				tLog.Info("request completed")
			}()

			next.ServeHTTP(rec, r)
		}
		return http.HandlerFunc(fn)
	}
}

//---------------------------------------------------------

func LogMiddleWareGin(c *gin.Context) {
	r := c.Request
	w := c.Writer
	logger := GetLogger(r)
	logFields := Fields{}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["url"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
	if rid, ok := r.Context().Value("uuid").(string); ok && rid != "" {
		logFields["request_id"] = rid
	}

	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method
	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	// try to get body from request
	var bodyReq interface{}
	content := r.Header.Get("Content-Type")
	if !strings.Contains(content, "image") && !strings.Contains(content, "multi") && r.Body != nil {
		buf, err := ioutil.ReadAll(r.Body)
		if err == nil {
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

			if err := json.NewDecoder(rdr1).Decode(bodyReq); err != nil {
				bodyReq = string(buf)
			}
			r.Body = rdr2
		}
	}
	logFields["request"] = bodyReq
	logger.WithFields(logFields).Info("request started")
	rec := NewRecorder(w, r)
	t1 := time.Now()
	defer func() {
		tLog := logger.WithFields(logFields).WithFields(Fields{
			"resp_status":       rec.ResponseStatus,
			"resp_bytes_length": len(rec.ResponseBody),
			"resp_elapsed_ms":   float64(time.Since(t1).Nanoseconds()) / 1000000.0,
		})

		tLog.Info("request completed")
	}()
	c.Next()
}

//---------------------------------------------------------

func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value("uuid") == "" {
				ctx := context.WithValue(r.Context(), "uuid", uuid.NewV4().String())
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

//---------------------------------------------------------
// End-of-file
//---------------------------------------------------------
