package main

import (
  "strings"
  "io/ioutil"

  log "github.com/sirupsen/logrus"
)


func ParseEntryFile(filePath string) [][]string {
  var splitEntries [][]string

  buf, err := ioutil.ReadFile(filePath)
  if err != nil {
    log.Fatal(err) 
  }

  entries := strings.Split(string(buf), "\n")
  if entries[len(entries) - 1] == "" {
    entries = entries[:len(entries) - 1] 
  }

  for _, entry := range entries {
    splitEntries = append(splitEntries, strings.Split(entry, ":"))
  }

  return splitEntries
}
