package main

import (
  "database/sql"

  _ "github.com/lib/pq"
  log "github.com/sirupsen/logrus"
)


func ConnectDatabase() {
  _, err := sql.Open("postgres", "user=eyoungberg dbname=myscprod sslmod=verify-full")
  if err != nil {
    log.Fatal(err) 
  }
}
