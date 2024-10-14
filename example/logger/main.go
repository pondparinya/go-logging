package main

//---------------------------------------------------------

import (
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "os/signal"
  // "strings"
  "syscall"
  "time"

  "go.kbtg.tech/733/go-ecslog"
)

//---------------------------------------------------------

func postRaw(url string, header http.Header, request []byte) ([]byte, error) {
  req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(request))
  if err != nil {
    return nil, err
  }

  if header != nil {
    for k, vs := range header {
      for _, v := range vs {
        req.Header.Add(k, v)
      }
    }
  }

  client := http.Client{}

  res, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  defer res.Body.Close()

  return ioutil.ReadAll(res.Body)
}

func postJson(url string, header http.Header, request, response interface{}) error {
  data, err := json.Marshal(request)
  if err != nil {
    return err
  }

  result, err := postRaw(url, header, data)
  if err != nil {
    return err
  }

  err = json.Unmarshal(result, response)
  return err
}

//---------------------------------------------------------

func handlerCall(w http.ResponseWriter, r *http.Request) {

  logger := ecslog.GetLogger(r)

  var req map[string]interface{}
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    logger.Errorf("Unable to parse request.")
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  fmt.Printf("map=%v\n", req)

  list, ok := req["to"].([]interface{})
  if !ok || len(list) == 0 {
    logger.Errorf("Invalid field 'to=%v'.", req["to"])
    http.Error(w, "Field 'to' is missing", http.StatusBadRequest)
    return
  }

  to, ok := list[0].(string)
  if !ok {
    logger.Errorf("Invalid field 'to=%v'.", list[0])
    http.Error(w, "Field 'to' is invalid", http.StatusBadRequest)
    return
  } else {
    req["to"] = list[1:]
  }

  logger = logger.Clone().SetFields(ecslog.Fields{"custom": "xyz"})

  r.Header.Set(ecslog.API_TID, logger.TraceId())
  var res map[string]interface{}
  if err := postJson(to, r.Header, &req, &res); err != nil {
    logger.Errorf("Unable to call since '%s'.", err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  s, err := json.Marshal(&res)
  if err != nil {
    logger.Errorf("Unable to marshal 'response=%v'.", res)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusOK)
  w.Write(s)
}

func handlerEcho(w http.ResponseWriter, r *http.Request) {
  buff, err := ioutil.ReadAll(r.Body)
  if err != nil {
    ecslog.Errorf("Unable to parse request.")
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  w.WriteHeader(http.StatusOK)
  w.Write(buff)
}

func handlerWait(w http.ResponseWriter, r *http.Request) {
  buff, err := ioutil.ReadAll(r.Body)
  if err != nil {
    ecslog.Errorf("Unable to parse request.")
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  time.Sleep(5 * time.Second)

  w.WriteHeader(http.StatusOK)
  w.Write(buff)
}

//---------------------------------------------------------

func startService(address string) {

  // Create listener for the 'SIGTERM'
  // from kernel
  trigger := make(chan os.Signal, 1)
  signal.Notify(trigger, os.Interrupt, syscall.SIGTERM)

  ecslog.Infof("Start serving at '%s'.", address)
  server := http.Server{Addr: address}

  // Wait for 'SIGTERM' from kernel
  var errRunning error
  go func() {
    http.HandleFunc("/call/", ecslog.LogHandler(handlerCall))
    http.HandleFunc("/echo/", ecslog.LogHandler(handlerEcho))
    http.HandleFunc("/wait/", ecslog.LogHandler(handlerWait))

    /*
       ecslog.Info(" |- /call/ - {\"to\": \"<URL>\", ...}")
       ecslog.Info(" |- /echo/ - {...}")
       ecslog.Info(" |- /wait/\n")
       //*/

    errRunning = server.ListenAndServe()
  }()
  <-trigger

  // Create the cancelable context for
  // help cancel the halted shutdown
  // process
  timeout := time.Duration(5) * time.Second
  srvCtx, srvCancel := context.WithTimeout(context.Background(), timeout)
  defer srvCancel()

  // Perform shutdown then wait until
  // the server finished the shutdown
  // process or the timeout had been
  // reached
  errShutdown := server.Shutdown(srvCtx)

  if errShutdown != nil {
    ecslog.Infof("Stop serving at '%s' since '%v'.", address, errShutdown)
  } else if errRunning != nil {
    ecslog.Infof("Stop serving at '%s' since '%v'.", address, errRunning)
  } else {
    ecslog.Infof("Stop serving at '%s'.", address)
  }
}

//---------------------------------------------------------

func main() {
  args := os.Args
  if len(args) != 3 {
    fmt.Printf("\nUsage: example <serviceId> <address=ip:port>\n")
    os.Exit(1)
  }

  ecslog.Setup(args[1])
  startService(args[2])
}

//---------------------------------------------------------
// End-of-file
//---------------------------------------------------------

