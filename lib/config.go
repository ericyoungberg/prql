package lib

import (
  "errors"

  "github.com/BurntSushi/toml"
  "github.com/prql/prql/lib/defaults"
)

type headers struct {
  Token  string 
  Secret string
}

type configFile struct {
  Port    int

  Host    string
  Secret  string
  LogFile string

  Headers headers
}

type Config struct {
  file configFile
}

func (c *Config) Port() int {
  if c.file.Port != 0 {
    return c.file.Port 
  }

  return defaults.Port
}

func (c *Config) Host() string {
  if c.file.Host != "" {
    return c.file.Host 
  }

  return defaults.Host
}

func (c *Config) Secret() (string, error) {
  if c.file.Secret == "" {
    return "", NoSecretErr
  }

  return c.file.Secret, nil
}

func (c *Config) LogFile() string {
  if c.file.LogFile != "" {
    return c.file.LogFile
  }

  return defaults.LogFile
}

func (c *Config) Headers() headers {
  if c.file.Headers.Token == "" {
    c.file.Headers.Token = defaults.HeadersToken 
  }

  if c.file.Headers.Secret == "" {
    c.file.Headers.Secret = defaults.HeadersSecret
  }
  
  return c.file.Headers
}


func loadConfig() (Config, error) {
  var loadedConfig Config

  if _, err := toml.DecodeFile(Sys.ConfigFile, &loadedConfig.file); err != nil {
    return loadedConfig, err 
  }

  return loadedConfig, nil
}


func GetConfig() (Config, error) {
  var err error

  if __PROVIDER == (Config{}) {
    __PROVIDER, err = loadConfig()
  }

  return __PROVIDER, err
}

var (
  __PROVIDER Config

  NoSecretErr = errors.New("No `Secret` value defined in prql.toml")
)
