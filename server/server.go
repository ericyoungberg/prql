package server

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"

  //"github.com/prql/prqld/auth"

  log "github.com/sirupsen/logrus"
)

const (
  HeaderName string = "X-PrQL-Token"
)

type Config struct {
  Port int16
}

func Startup(config *Config) {
  mux := http.NewServeMux()
  port := fmt.Sprintf(":%d", config.Port)

  defer func() {
    log.Println("Starting server")
    http.ListenAndServe(port, mux)
  }()

  mux.HandleFunc("/", handler)
}


/**
* Private
*/

var (
  ipLogger log.FieldLogger
)

type postRequestBody struct {
  Query string
}

func handler(w http.ResponseWriter, r *http.Request) {
  ipLogger = log.WithFields(log.Fields{"IP": r.RemoteAddr})

  token := r.Header.Get(HeaderName)
  if token == "" {
    fail(w, "No token")
    return
  }
  ipLogger.Info(fmt.Sprintf("Received request using token %s", token))

  query := getQuery(r)
  if query == "" {
    fail(w, "No query")
    return
  }
  ipLogger.Info(fmt.Sprintf("Token %s was used for the following query: \"%s\"", token, query))
}


func getQuery(r *http.Request) string {
  var query string

  if r.Method == "GET" {
    q := r.URL.Query()

    if len(q["query"]) != 0 {
      query = q["query"][0]
    }
  }

  if r.Method == "POST" {
    var body postRequestBody

    bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
      ipLogger.Panic(err) 
    }

    if len(bodyBytes) != 0 {
      err = json.Unmarshal(bodyBytes, &body)
      if err != nil {
        ipLogger.Panic(err) 
      }

      if body.Query != "" {
        query = body.Query
      }
    }
  }

  return query
} 

func fail(w http.ResponseWriter, message string) {
  ipLogger.Error(message)
  http.Error(w, message, http.StatusBadRequest)
}
