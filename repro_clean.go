package main

import (
  "fmt"
  "os"
  "path/filepath"
  "github.com/xef/xefcli/internal/app/file"
  "github.com/xef/xefcli/internal/core/logger"
  "github.com/xef/xefcli/internal/infrastructure/filesystem"
)

func main() {
  dir, err := os.MkdirTemp("", "repro-clean")
  if err != nil { panic(err) }
  defer os.RemoveAll(dir)
  if err := os.WriteFile(filepath.Join(dir, "a.tmp"), []byte("temp"), 0644); err != nil { panic(err) }
  if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("keep"), 0644); err != nil { panic(err) }
  fs := filesystem.NewOSFileSystem()
  cmd := file.NewCommand(fs, logger.Nop())
  cmd.SetArgs([]string{"clean", dir})
  fmt.Println("before execute")
  err = cmd.Execute()
  fmt.Println("after execute", err)
}
