package main

import (
  "github.com/prql/prqld/server"
)

func main() {
  server.Startup(&server.Config{Port: 1999})
}
