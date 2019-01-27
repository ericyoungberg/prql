package pools

import (
  "os"
  "strings"
)

func AppendEntry(filePath string, entry []string) error { 
  fd, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
  if err != nil {
    return err
  }
  
  defer fd.Close()

  _, err = fd.WriteString(strings.Join(entry, entryDelimiter) + "\n")

  return err
}
