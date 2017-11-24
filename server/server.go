package server

import (
  "fmt"
  "log"
  "net/http"
)

type Config struct {
  Port int16
  ReadTimeout int16
  WriteTimeout int16
}

func handler(writer http.ResponseWriter, r *http.Request) {
  fmt.Println("Received request!")
}


/*
* Public
*/

func Start(config Config) {

  s := &http.Server {
    Addr: fmt.Sprintf(":%d", config.Port),
  }

  s.HandleFunc("/", handler)

  log.Fatal(s.ListenAndServe())
}
