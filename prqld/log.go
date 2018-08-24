package main

import (
  "github.com/sirupsen/logrus"
)


type Logger struct {
  System *logrus.Logger
  Console *logrus.Logger
}

var log *Logger

func setupLogger() {
  log = &Logger{
    System: logrus.New(),
    Console: logrus.New(),
  } 
}

func (logger *Logger) Info(args ...interface{}) {
  logger.System.Info(args)
  logger.Console.Info(args)
}

func (logger *Logger) Panic(args ...interface{}) {
  logger.System.Panic(args)
  logger.Console.Panic(args)
}

func (logger *Logger) Error(args ...interface{}) {
  logger.System.Error(args)
  logger.Console.Error(args)
}

func (logger *Logger) Fatal(args ...interface{}) {
  logger.System.Fatal(args)
  logger.Console.Fatal(args)
}
