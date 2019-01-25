package pools


const (
  entryDelimiter string = ":"
)


type Pool struct {
  Entries [][]string
  FilePath string
}

func (p *Pool) Load() {
  buf, err := ioutil.ReadFile(filePath)
  if err != nil {
    log.Fatal(err) 
  }

  entries := strings.Split(string(buf), "\n")
  if entries[len(entries) - 1] == "" {
    entries = entries[:len(entries) - 1] 
  }

  for _, entry := range entries {
    p.Entries = append(p.Entries, strings.Split(entry, entryDelimiter))
  }
}

func (p *Pool) Save() {
  lines := make([]string, len(p.Entries))

  for i, entry := range p.Entries {
    lines[i] = strings.Join(entry, entryDelimiter)
  } 

  data := []byte(strings.Join(lines, "\n"))

  return ioutil.WriteFile(p.FilePath, data, 0600)
}

func (p *Pool) Refresh() {

}

func NewPool(filePath string) *Pool {
  pool := &Pool{ FilePath: filePath }
  pool.Load()

  return pool
}
