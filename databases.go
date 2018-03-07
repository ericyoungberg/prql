package main

import (
  "strconv"
  "database/sql"

  _ "github.com/lib/pq"
  log "github.com/sirupsen/logrus"
)


type DatabaseEntry struct {
  tag string
  driver string
  host string
  port int
  ssl bool
}

var (
  DatabasePool = make(map[string]DatabaseEntry)
)


/**
* Database File Entry Schema
*
* tag:driver:host:port:ssl
*
* tag - A string used to identify the database server.
*
* driver -  The type of database server. eg: postgres, mysql, ...
* 
* host - The address of the database server.
* 
* port - The host's port number where the database server is listening.
*
* ssl - A boolean that indicates whether we should verify ssl or not.
*/

func PopulateDatabasePool() {
  entries := ParseEntryFile("/var/lib/prql/databases")

  for i, parts := range entries {
    if len(parts) != 5 {
      log.Error("Invalid database entry at line " + strconv.Itoa(i + 1)) 
      continue
    }

    ssl, err := strconv.ParseBool(parts[4])
    if err != nil {
      ssl = false 
    }

    port, err := strconv.Atoi(parts[3])
    if err != nil {
      port = 5432
    }

    DatabasePool[parts[0]] = DatabaseEntry {
      driver: parts[1],
      host: parts[2],
      port: port,
      ssl: ssl,
    }
  }
}


func ConnectDatabase() {
  _, err := sql.Open("postgres", "user=eyoungberg dbname=myscprod sslmod=verify-full")
  if err != nil {
    log.Fatal(err) 
  }
}
