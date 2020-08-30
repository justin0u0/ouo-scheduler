package main

import (
  "fmt"
  "os"
	"math/rand"
	"time"

  "k8s.io/component-base/logs"
	"github.com/justin0u0/ouo-scheduler/pkg/plugin"
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
