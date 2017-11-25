package main

import (
  "github.com/prql/prqld/server"
)

func main() {
  server.Start(&server.Config{Port: 1999})
}
