package lib

import (
  "path"
)

type sys struct {
  ConfigFile   string
  TokenFile    string
  DatabaseFile string
  FilesPath    string
}

const (
  filesPath string = "/var/lib/prql"
)

var (
  Sys sys = sys{
    FilesPath: filesPath,
    ConfigFile: path.Join(filesPath, "prql.toml"),
    TokenFile: path.Join(filesPath, "tokens"),
    DatabaseFile: path.Join(filesPath, "databases"),
  }
)
