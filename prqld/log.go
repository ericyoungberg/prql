package main

import (
  "os"

  "github.com/prql/prql/lib"
  "github.com/sirupsen/logrus"
)


type Logger struct {
  logrus.Logger

  Console *logrus.Logger
  system  *logrus.Logger

  config  *lib.Config
}


var (
  log *Logger
)


func setupLogger(config *lib.Config) {
  if log == nil {
    log = &Logger{
      system: logrus.New(),
      Console: logrus.New(),
    } 
  }

  log.config = config
}

func (logger *Logger) onSystem(fn func(...interface{}), args ...interface{}) {
  logFile, err := os.OpenFile(logger.config.LogFile(), os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
  if err != nil {
    return
  }
  defer logFile.Close()

  logger.system.Out = logFile
  fn(args...)
}

func (logger *Logger) Info(args ...interface{}) {
  logger.onSystem(logger.system.Info, args...)
  logger.Console.Info(args...)
}

func (logger *Logger) Panic(args ...interface{}) {
  logger.onSystem(logger.system.Panic, args...)
  logger.Console.Panic(args...)
}

func (logger *Logger) Error(args ...interface{}) {
  logger.onSystem(logger.system.Error, args...)
  logger.Console.Error(args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
  logger.onSystem(logger.system.Fatal, args...)
  logger.Console.Fatal(args...)
}
