package pools

import (
  "fmt"
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
    splitEntries = append(splitEntries, strings.Split(entry, entryDelimiter))
  }

  return splitEntries
}

func removeByColumn(values []string, dataset [][]string, col int) [][]string {
  for _, value := range values {
    for i, row := range dataset {
      if row[col] == value {
        dataset = append(dataset[:i], dataset[i+1:]...)

        fmt.Printf("Deleting %s\n", value) 
        break
      }
    }
  }

  return dataset
}
