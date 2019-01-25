package pools

import (
  "os"
  "strings"
  "strconv"
  "io/ioutil"

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

type TokenEntry struct{
  Living bool

  User     string
  Password string
  HostName string
  DBName string

  Origins []string
}


/**
* Database File Entry Schema
*
* name:driver:host:port:ssl
*
* name - A string used to identify the database server.
*
* driver -  The type of database server.
* 
* host - The address of the database server.
* 
* port - The host's port number where the database server is listening.
*
* ssl - A boolean that indicates whether we should verify ssl or not.
*/
func GetDatabaseEntries() map[string]DatabaseEntry {
  databases := make(map[string]DatabaseEntry) 
  entries := ParseEntryFile(lib.Sys.DatabaseFile)

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

    databases[parts[0]] = DatabaseEntry{
      HostName: parts[0],
      Driver: parts[1],
      Host: parts[2],
      Port: port,
      SSL: ssl,
    }
  }

  return databases
}

/**
* Token File Entry Schema
*
* token:username:password:hostName:dbname:origins:living
*
* token - 32-character string generated by the cli. Used to identify credentials 
*         and passed to the program from the client in the X-PrQL-Token header.
*
* username - The database username.
* 
* password - The database user's password.
* 
* hostName - The tag that was defined by the user while using the cli to store the credentials
*              of the database server.
*
* dbname - The name of the database to create the connection to.
*
* origins - A comma separated list of authorized origins that this entry's token can
*           be used with. If left empty, then all origins are authorized.
*
* living - A boolean which indicates whether this entry will spawn a lifelong connection
*          to the specified database whenever the system starts up.
*/

func GetTokenEntries() map[string]TokenEntry {
  tokens := make(map[string]TokenEntry) 
  entries := ParseEntryFile(lib.Sys.TokenFile)

  for i, parts := range entries {
    if len(parts) != 7 {
      log.Error("Invalid token entry at line " + strconv.Itoa(i + 1)) 
      continue
    }

    originEntries := strings.Split(parts[5], ",")
    origins := make([]string, len(originEntries))
    originIndex := 0
    for _, origin := range originEntries {
      if origin != "" {
        origins[originIndex] = origin
        originIndex += 1
      }
    }

    if len(origins) == 1 && origins[0] == "" {
      origins = nil 
    }

    living, err := strconv.ParseBool(parts[6])
    if err != nil {
      living = false 
    }

    password, err := lib.InsecureDecryptString(parts[2])
    if err != nil {
      log.Error("Couldn't decrypt password at line " + strconv.Itoa(i + 1)) 
      continue
    }

    tokens[parts[0]] = TokenEntry{
      User: parts[1], 
      Password: password, 
      HostName: parts[3], 
      DBName: parts[4], 
      Origins: origins, 
      Living: living,
    }
  }

  return tokens
}


func WriteEntryFile(filePath string, entries [][]string) error {
}


func AppendEntry(filePath string, entry []string) error { 
  fd, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
  if err != nil {
    return err
  }
  
  defer fd.Close()

  _, err = fd.WriteString(strings.Join(entry, entryDelimiter) + "\n")

  return err
}
