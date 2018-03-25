package util

import (
  "os"
  "strings"
  "io/ioutil"

  log "github.com/sirupsen/logrus"
)


const (
  ENTRY_DELIMITER string = ":"
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
    splitEntries = append(splitEntries, strings.Split(entry, ENTRY_DELIMITER))
  }

  return splitEntries
}


func WriteEntryFile(filePath string, entries [][]string) error {
  lines := make([]string, len(entries))

  for i, entry := range entries {
    lines[i] = strings.Join(entry, ENTRY_DELIMITER)
  } 

  data := []byte(strings.Join(lines, "\n"))

  return ioutil.WriteFile(filePath, data, 0600)
}


func AppendEntry(filePath string, entry []string) error { 
  fd, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
  if err != nil {
    return err
  }
  
  defer fd.Close()

  _, err = fd.WriteString(strings.Join(entry, ENTRY_DELIMITER) + "\n")

  return err
}
