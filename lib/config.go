package lib

import (
  "path"

  "github.com/BurntSushi/toml"
)


type Config struct {
  Port int16 
  Secret string
}


var config Config


func loadConfig() (Config, error) {
  var loadedConfig Config

  tomlPath := path.Join(Sys.FilesPath, Sys.ConfigFile)

  if _, err := toml.DecodeFile(tomlPath, &loadedConfig); err != nil {
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
