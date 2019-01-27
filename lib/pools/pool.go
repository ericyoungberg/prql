package pools

import (
  "strings"
  "io/ioutil"

  log "github.com/sirupsen/logrus"
)

const (
  entryDelimiter string = ":"
)


type pool struct {
  FilePath string

  records [][]string
}

func (p *pool) Save() error {
  lines := make([]string, len(p.records))

  for i, entry := range p.records {
    lines[i] = strings.Join(entry, entryDelimiter)
  } 

  data := []byte(strings.Join(lines, "\n"))

  return ioutil.WriteFile(p.FilePath, data, 0600)
}

func (p *pool) Build() {
  log.Fatal("pool.Build() must be overriden")
}

func (p *pool) load() {
  buf, err := ioutil.ReadFile(p.FilePath)
  if err != nil {
    log.Fatal(err) 
  }

  entries := strings.Split(string(buf), "\n")
  if entries[len(entries) - 1] == "" {
    entries = entries[:len(entries) - 1] 
  }

  p.records = make([][]string, len(entries))
  for _, entry := range entries {
    p.records = append(p.records, strings.Split(entry, entryDelimiter))
  }
}
