package lib

import (
  "github.com/BurntSushi/toml"
)

type headers struct {
  Token  string 
  Secret string
}

type Config struct {
  Port    int
  Secret  string

  LogFile string

  Headers headers
}

var (
  config Config
)

func loadConfig() (Config, error) {
  var loadedConfig Config

  if _, err := toml.DecodeFile(Sys.ConfigFile, &loadedConfig); err != nil {
    return loadedConfig, err 
  }

  return loadedConfig, nil
}


func GetConfig() (Config, error) {
  var err error

  if config == (Config{}) {
    config, err = loadConfig()
  }

  return config, err
}
