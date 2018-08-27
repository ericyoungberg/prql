package main

import (
  "fmt"
  "time"
  "net/url"
  "net/http"
  "io/ioutil"
  "encoding/json"

  "github.com/prql/prql/lib"
)

type postRequestBody struct {
  Query string
}

var (
  host string
)

func startServer(config *lib.Config) {
  mux := http.NewServeMux()
  port := fmt.Sprintf(":%d", config.Port)
  host = fmt.Sprintf("127.0.0.1%s", port)

  refreshTokens := lib.SecretExec(populateTokenPool)
  refreshDatabases := lib.SecretExec(populateDatabasePool)

  mux.HandleFunc("/", handler)
  mux.HandleFunc("/refresh-tokens", refreshTokens)
  mux.HandleFunc("/refresh-databases", refreshDatabases)
  mux.HandleFunc("/check", func(w http.ResponseWriter, req *http.Request) { 
    w.WriteHeader(http.StatusOK)
  })

  go checkServerStatus()

  log.Info("Starting server")
  http.ListenAndServe(port, mux)
}

func checkServerStatus() {
  running := false

  for i := 0; i < 10; i += 1 {
    time.Sleep(time.Second) 

    endpoint := url.URL{Scheme: "http", Host: host, Path: "check"}
    res, err := http.Get(endpoint.String())
    if err != nil {
      log.Error(err)
      continue
    }

    res.Body.Close()

    if res.StatusCode != http.StatusOK {
      continue 
    }

    running = true
    break
  } 

  if running {
    log.Info(fmt.Sprintf("Server listening at %s", host))
  } else {
    log.Panic(fmt.Sprintf("Cannot connect to server at %s.\nExiting...", host)) 
  }
}

func handler(w http.ResponseWriter, r *http.Request) {
  config, err := lib.GetConfig()
  if err != nil {
    log.Fatal(err) 
  }

  var token string

  token = r.Header.Get(config.Headers.Token)
  if token == "" {
    tokenParam := r.URL.Query()["token"]

    if len(tokenParam) != 0 {
      token = tokenParam[0] 
    } else {
      fail(w, "No token")
      return
    }
  }
  log.Info(fmt.Sprintf("Received request using token %s", token))

  tokenEntry := tokenPool[token]
  if tokenEntry.Origins != nil && len(tokenEntry.Origins) > 0 {
    unauthorized := true
    for _, authorized := range tokenEntry.Origins {
      log.Info(fmt.Sprintf("%s|", authorized))

      if authorized == r.RemoteAddr {
        unauthorized = false
        break 
      }
    }

    if unauthorized {
      fail(w, fmt.Sprintf("Origin %s is not authorized to use token", r.RemoteAddr)) 
      return
    }
  }

  query := getQuery(r)
  if query == "" {
    fail(w, "No query")
    return
  }
  log.Info(fmt.Sprintf("Token %s was used for the following query: \"%s\"", token, query))

  data, err := performQuery(query, token)
  if err != nil {
    fail(w, err.Error())
  }

  json.NewEncoder(w).Encode(data)
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
      log.Panic(err) 
    }

    if len(bodyBytes) != 0 {
      err = json.Unmarshal(bodyBytes, &body)
      if err != nil {
        log.Panic(err) 
      }

      if body.Query != "" {
        query = body.Query
      }
    }
  }

  return query
} 

func fail(w http.ResponseWriter, message string) {
  log.Error(message)
  http.Error(w, message, http.StatusBadRequest)
}
