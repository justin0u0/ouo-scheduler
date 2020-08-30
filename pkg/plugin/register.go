package plugin

import (
  "github.com/spf13/cobra"
  "k8s.io/kubernetes/cmd/kube-scheduler/app"
  "github.com/justin0u0/ouo-scheduler/pkg/plugin/ouo"
)

func Register() *cobra.Command {
  return app.NewSchedulerCommand(
    app.WithPlugin(ouo.Name, ouo.New)
  )
}
