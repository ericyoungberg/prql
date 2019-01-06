package databases

import (
  "fmt"
  "errors"
  "database/sql"

  "github.com/prql/prql/lib"
  "github.com/go-sql-driver/mysql"
)

type Connection struct {
  IsClosed bool

  connection *sql.DB
}


func (conn *Connection) Close() {
  conn.connection.Close()  
  conn.IsClosed = true
}


func NewConnection(tokenEntry *lib.TokenEntry, databaseEntry *lib.DatabaseEntry) (*Connection, error) {
  var db *sql.DB

  dbConnString, err := generateDSN(tokenEntry, databaseEntry)
  if err != nil {
    return &Connection{}, err
  } else {
    db, err = sql.Open(databaseEntry.Driver, dbConnString)
    if err != nil {
      return &Connection{}, err
    }

    conn := &Connection{
      IsClosed: false,
      connection: db,
    }

    return conn, nil
  }
}


func generateDSN(token *lib.TokenEntry, database *lib.DatabaseEntry) (string, error) {
  var dsn string = ""
  var err error = nil

  switch database.Driver {
    case "mysql":
      dsnConfig := mysql.NewConfig()
      dsnConfig.User = token.User
      dsnConfig.Passwd = token.Password
      dsnConfig.Net = "tcp"
      dsnConfig.Addr = fmt.Sprintf("%s:%d", database.Host, database.Port)
      dsnConfig.DBName = token.DBName
      dsn = dsnConfig.FormatDSN()

    case "postgres":
      dbConnStringFmt := "user=%s password=%s dbname=%s host=%s port=%d sslmode=disable"
      dbConnStringVars := []interface{}{token.User, token.Password, token.DBName, database.Host, database.Port}
      dsn = fmt.Sprintf(dbConnStringFmt, dbConnStringVars...)

    default:
      err = errors.New(fmt.Sprintf("database driver %s does not exist", database.Driver))
  }

  return dsn, err
}


