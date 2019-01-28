package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"

  "github.com/prql/prql/lib"
  "github.com/prql/prql/lib/pools"
  "github.com/prql/prql/prqld/version"
)

type postRequestBody struct {
  Query string
}

type Server struct {
  port int
  host string
}

func (server *Server) StartFromConfig(config *lib.Config) {
  server.host = config.Host()
  server.port = config.Port()

  server.start()
}

func (server *Server) start() {
  log.Info("Starting server")

  refreshTokens := lib.SecretExec(pools.GetTokenPool().Build)
  refreshDatabases := lib.SecretExec(pools.GetDatabasePool().Build)

  mux := http.NewServeMux()
  mux.HandleFunc("/", handleDataRequest)
  mux.HandleFunc("/refresh-tokens", refreshTokens)
  mux.HandleFunc("/refresh-databases", refreshDatabases)
  mux.HandleFunc("/check", respondOK)
  mux.HandleFunc("/version", respondVersion)

  go lib.CheckServerStatus()

  http.ListenAndServe(fmt.Sprintf(":%d", server.port), mux)
}

func authorizedOrigin(origin string, entry pools.TokenEntry) bool {
  unauthorized := true

  if entry.Origins != nil && len(entry.Origins) > 0 {
    for _, authorized := range entry.Origins {
      if authorized == origin {
        unauthorized = false
        break 
      }
    }
  } else {
    unauthorized = false 
  }

  return !unauthorized
}

func fail(w http.ResponseWriter, message string) {
  log.Error(message)
  http.Error(w, message, http.StatusBadRequest)
}

func handleDataRequest(w http.ResponseWriter, r *http.Request) {
  token := readToken(r)
  if token == "" {
    fail(w, "No token")
    return
  }

  requestOrigin := r.Header.Get("Origin")
  if requestOrigin == "" {
    requestOrigin = "[Unknown Origin]"
  }
  log.Info(fmt.Sprintf("Received request from %s using token %s", requestOrigin, token))

  tokenEntry := pools.GetTokenPool().Entries[token]
  if !authorizedOrigin(requestOrigin, tokenEntry) {
    fail(w, fmt.Sprintf("Origin %s is not authorized to use token", requestOrigin)) 
    return
  }

  query := readQuery(r)
  if query == "" {
    fail(w, "No query")
    return
  }
  log.Info(fmt.Sprintf("Token %s was used for the following query: \"%s\"", token, query))

  data, err := performQuery(query, token)
  if err != nil {
    fail(w, err.Error())
    return
  }

  json.NewEncoder(w).Encode(data)
}

func readQuery(r *http.Request) string {
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
      if err == nil && body.Query != "" {
        query = body.Query
      }
    }
  }

  return query
} 

func readToken(r *http.Request) string {
  config, err := lib.GetConfig()
  if err != nil {
    log.Error(err) 
    return ""
  }

  token := r.Header.Get(config.Headers().Token)
  if token == "" {
    tokenParam := r.URL.Query()["token"]

    if len(tokenParam) != 0 {
      token = tokenParam[0] 
    }
  }

  return token
}

func respondOK(w http.ResponseWriter, req *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func respondVersion(w http.ResponseWriter, req *http.Request) {
  w.Write([]byte(version.VERSION))
}
