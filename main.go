package main

import (
  "github.com/prql/prqld/server"

  //"github.com/coreos/go-systemd/daemon"
)

func main() {
  go server.Startup(&server.Config{Port: 1999})
}
