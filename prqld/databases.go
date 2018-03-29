package main

import (
  "fmt"
  "strconv"
  "database/sql"

  _ "github.com/lib/pq"
  "github.com/prql/prql/util"
  log "github.com/sirupsen/logrus"
)


type DatabaseEntry struct {
  ssl bool

  port int

  hostName string
  driver   string
  host     string
}

var (
  DatabasePool = make(map[string]DatabaseEntry)
  databaseConnections = make(map[string]*sql.DB)
)


/**
* Database File Entry Schema
*
* name:driver:host:port:ssl
*
* name - A string used to identify the database server.
*
* driver -  The type of database server. eg: postgres, mysql, ...
* 
* host - The address of the database server.
* 
* port - The host's port number where the database server is listening.
*
* ssl - A boolean that indicates whether we should verify ssl or not.
*/

func PopulateDatabasePool(refresh bool) {
  if refresh {
    DatabasePool = make(map[string]DatabaseEntry) 
  }

  entries := util.ParseEntryFile("/var/lib/prql/databases")

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


func CloseDatabaseConnections() {
  for k, v := range databaseConnections {
    v.Close()  
    delete(databaseConnections, k)
  }
}


func PerformQuery(query string, token string) (map[string]interface{}, error) {
  db := getDatabase(token)
  
  rows, err := db.Query(query)
  if err != nil {
    return nil, err
  }

  defer rows.Close()

  return structureData(rows)
}



/**
* Private
*/

func getDatabase(token string) *sql.DB {
  var db *sql.DB
  var ok bool
  var err error

  tokenEntry, ok := TokenPool[token]
  if ok != true {
    IpLogger.Panic("Invalid token") 
  }

  databaseEntry, ok := DatabasePool[tokenEntry.hostName]
  if ok != true {
    IpLogger.Panic("Invalid database server name")
  }

  db, ok = databaseConnections[tokenEntry.hostName] 
  if ok != true {
    dbConnStringFmt := "user=%s password=%s dbname=%s host=%s port=%d sslmode=disable"
    dbConnStringVars := []interface{}{tokenEntry.user, tokenEntry.password, tokenEntry.dbname, databaseEntry.host, databaseEntry.port}
    dbConnString := fmt.Sprintf(dbConnStringFmt, dbConnStringVars...)

    db, err = sql.Open(databaseEntry.driver, dbConnString)
    if err != nil {
      IpLogger.Error(err) 
    }

    databaseConnections[tokenEntry.hostName] = db
  }

  return db
}


func structureData(rows *sql.Rows) (map[string]interface{}, error) {
  var structured = make(map[string]interface{})

  colTypes, err := rows.ColumnTypes()
  if err != nil {
    return nil, err
  }

  var colNames = make([]string, len(colTypes))
  var fields = make(map[string]map[string]string)

  for i, colType := range colTypes {
    fields[colType.Name()] = map[string]string{ "type": colType.DatabaseTypeName() }
    colNames[i] = colType.Name()
  }
  structured["fields"] = fields

  rawData := make([][]byte, len(colTypes))
  buf := make([]interface{}, len(colTypes))

  for i, _ := range rawData {
    buf[i]  = &rawData[i]
  }

  var structuredRows []map[string]interface{}
  var rowNum int16 = 0

  for rows.Next() {
    err = rows.Scan(buf...)
    if err != nil {
      IpLogger.Error(err)
    }

    structuredRows = append(structuredRows, make(map[string]interface{}))

    for i, raw := range rawData {
      colName := colNames[i]

      if raw == nil {
        structuredRows[rowNum][colName] = nil
      } else {
        temp := string(raw)
        var err error

        switch fields[colName]["type"] {
        case "BOOL":
          structuredRows[rowNum][colName], err = strconv.ParseBool(temp)
        case "INT4", "INT8", "INT16", "INT32", "INT64":
          structuredRows[rowNum][colName], err = strconv.Atoi(temp)
        case "FLOAT4", "FLOAT8", "FLOAT16", "FLOAT32", "NUMERIC":
          structuredRows[rowNum][colName], err = strconv.ParseFloat(temp, 64)
        default:
          structuredRows[rowNum][colName] = temp
        }

        if err != nil {
          IpLogger.Error(err)
          structuredRows[rowNum][colName] = temp
        }
      }
    }

    rowNum = rowNum + 1
  }
  structured["rows"] = structuredRows
  structured["total_rows"] = rowNum

  return structured, nil
}
