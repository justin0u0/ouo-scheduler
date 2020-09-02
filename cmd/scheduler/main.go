package main

import (
  "fmt"
  "math/rand"
  "os"
  "time"

  "github.com/justin0u0/ouo-scheduler/pkg/plugin"
  "k8s.io/component-base/logs"
)

func main() {
  rand.Seed(time.Now().UnixNano())

  command := plugin.Register()

  logs.InitLogs()
  defer logs.FlushLogs()

  if err := command.Execute(); err != nil {
    _, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
    os.Exit(1)
  }
}
