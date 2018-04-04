package lib

type sys struct {
  ConfigFile string
  FilesPath  string
}

var Sys sys = sys{
  FilesPath: "/var/lib/prql",
  ConfigFile: "prql.toml",
}
