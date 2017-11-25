package server

import (
  "fmt"
  "log"
  "time"
  "net/http"
)

type Config struct {
  Port int16
  ReadTimeout time.Duration
  WriteTimeout time.Duration
}

func handler(writer http.ResponseWriter, r *http.Request) {
  fmt.Println("Received request!")
}


/*
* Public
*/

func Start(config *Config) {

  // Give defaults for timeouts
  if config.ReadTimeout == 0 {
    config.ReadTimeout = 10 * time.Second
  }

  if config.WriteTimeout == 0 {
    config.WriteTimeout = 10 * time.Second
  }

  // Setup server
  s := &http.Server {
    Addr: fmt.Sprintf(":%d", config.Port),
    ReadTimeout: config.ReadTimeout,
    WriteTimeout: config.WriteTimeout,
  }

  s.HandleFunc("/", handler)

  log.Fatal(s.ListenAndServe())
}
