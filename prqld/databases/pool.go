package databases

import (
  "strconv"
  "database/sql"

  _ "github.com/lib/pq"
  "github.com/prql/prql/lib"
  log "github.com/sirupsen/logrus"
  "github.com/prql/prql/prqld/pools"
)

var (
  databasePool map[string]lib.DatabaseEntry
  databaseConnections = make(map[string]*Connection)
)


func populateDatabasePool() {
  databasePool = make(map[string]lib.DatabaseEntry)
  databasePool = lib.GetDatabaseEntries()
}


func closeDatabaseConnections() {
  for k, v := range databaseConnections {
    v.Close()  
    delete(databaseConnections, k)
  }
}


func performQuery(query string, token string) (map[string]interface{}, error) {
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
  var ok bool
  pools := pool.GetPoolBroker()

  tokenEntry := pools.Get("token", token)
  if tokenEntry == nil {
    log.Panic("Invalid token") 
  }

  databaseEntry := pool.Get("database", tokenEntry.HostName)
  if datasebaseEntry == nil {
    log.Panic("Invalid database server name")
  }

  db, ok := databaseConnections[token] 
  if ok != true {
    conn, err := NewConnection(&lib.TokenEntry(tokenEntry), &lib.DatabaseEntry(databaseEntry))
    if err != nil {
      log.Panic("Couldn't connect to database %s", databaseEntry.HostName) 
    }

    databaseConnections[token] = conn
    db = databaseConnections[token]
  }

  return db.connection
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
    buf[i] = &rawData[i]
  }

  var structuredRows []map[string]interface{}
  var rowNum uint32 = 0

  for rows.Next() {
    err = rows.Scan(buf...)
    if err != nil {
      log.Error(err)
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
          log.Error(err)
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
