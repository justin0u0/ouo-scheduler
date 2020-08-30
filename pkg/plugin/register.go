package plugin

import (
	"github.com/justin0u0/ouo-scheduler/pkg/plugin/ouo"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func Register() *cobra.Command {
	return app.NewSchedulerCommand(app.WithPlugin(ouo.Name, ouo.New))
}
