package main

import (
  "os"

  "github.com/prql/prql/lib"
  "github.com/sirupsen/logrus"
)


type Logger struct {
  Console *logrus.Logger
  system  *logrus.Logger
  config  *lib.Config
}

var log *Logger

func setupLogger(config *lib.Config) {
  if log == nil {
    log = &Logger{
      system: logrus.New(),
      Console: logrus.New(),
      config: config,
    } 
  }

}

func (logger *Logger) openSystemLogger() (*os.File, error) {
  fd, err := os.OpenFile(logger.config.LogFile, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
  if err != nil {
    return nil, err 
  }

  log.system.Out = fd 

  return fd, nil
}

func (logger *Logger) Info(args ...interface{}) {
  logFile, err := logger.openSystemLogger() 
  if err == nil {
    logger.system.Info(args)
    defer logFile.Close()
  }

  logger.Console.Info(args)
}

func (logger *Logger) Panic(args ...interface{}) {
  logFile, err := logger.openSystemLogger() 
  if err == nil {
    logger.system.Panic(args)
    defer logFile.Close()
  }

  logger.Console.Panic(args)
}

func (logger *Logger) Error(args ...interface{}) {
  logFile, err := logger.openSystemLogger() 
  if err == nil {
    logger.system.Error(args)
    defer logFile.Close()
  }

  logger.Console.Error(args)
}

func (logger *Logger) Fatal(args ...interface{}) {
  logFile, err := logger.openSystemLogger() 
  if err == nil {
    logger.system.Fatal(args)
    defer logFile.Close()
  }

  logger.Console.Fatal(args)
}
