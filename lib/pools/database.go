package pools

import (
  "strconv"

  "github.com/prql/prql/lib"
  log "github.com/sirupsen/logrus"
)


type DatabaseEntry struct{
  SSL bool

  Port int

  HostName string
  Driver   string
  Host     string
}

type DatabasePool struct {
  pool

  Entries map[string]DatabaseEntry
}

func (p *DatabasePool) build() {
  databases := make(map[string]DatabaseEntry) 

  for i, parts := range p.records {
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

    databases[parts[0]] = DatabaseEntry{
      HostName: parts[0],
      Driver: parts[1],
      Host: parts[2],
      Port: port,
      SSL: ssl,
    }
  }

  p.Entries = databases
}

func GetDatabasePool() *DatabasePool {
  if __DB_INIT {
    __DB_PROVIDER = loadDatabasePool() 
  }

  return __DB_PROVIDER
}

func loadDatabasePool() *DatabasePool {
  __DB_INIT = false

  databasePool := &DatabasePool{
    pool: pool{FilePath: lib.Sys.DatabaseFile},
  }
  databasePool.self = databasePool
  databasePool.Build()
  
  return databasePool
}

var (
  __DB_PROVIDER *DatabasePool
  __DB_INIT bool = true
)
