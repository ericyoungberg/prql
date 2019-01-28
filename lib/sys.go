package lib

import (
  "fmt"
  "path"
  "time"
  "net/url"
  "net/http"

  log "github.com/sirupsen/logrus"
)

const filesPath string = "/var/lib/prql"


type sys struct {
  ConfigFile   string
  TokenFile    string
  DatabaseFile string
  FilesPath    string
}

var (
  Sys sys = sys{
    FilesPath: filesPath,
    ConfigFile: path.Join(filesPath, "prql.toml"),
    TokenFile: path.Join(filesPath, "tokens"),
    DatabaseFile: path.Join(filesPath, "databases"),
  }
)

func CheckServerStatus() {
  config, err := GetConfig()
  if err != nil {
    log.Fatal(err)  
  }

  host := fmt.Sprintf("%s:%d", config.Host(), config.Port())
  running := false

  i := 0

  for ; i < 10; i += 1 {
    time.Sleep(time.Second) 

    endpoint := url.URL{Scheme: "http", Host: host, Path: "check"}
    res, err := http.Get(endpoint.String())
    if err != nil {
      fmt.Print(".")
      continue
    }

    res.Body.Close()

    if res.StatusCode != http.StatusOK {
      continue 
    }


    running = true
    break
  } 

  if i > 0 {
    fmt.Print("\n")
  }

  if running {
    log.Info(fmt.Sprintf("Server listening at %s", host))
  } else {
    log.Fatal(fmt.Sprintf("Cannot connect to server at %s.\nExiting...", host)) 
  }
}
