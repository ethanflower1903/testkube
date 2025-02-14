package executors

import (
	"os"

	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common"
	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common/render"
	"github.com/kubeshop/testkube/pkg/ui"
	"github.com/spf13/cobra"
)

func NewGetExecutorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "executor [executorName]",
		Aliases: []string{"executors", "er"},
		Short:   "Gets executor details",
		Long:    `Gets executor, you can change output format`,
		Run: func(cmd *cobra.Command, args []string) {
			client, namespace := common.GetClient(cmd)

			if len(args) > 0 {
				name := args[0]

				executor, err := client.GetExecutor(name, namespace)
				ui.ExitOnError("getting executor: "+name, err)
				err = render.Obj(cmd, executor, os.Stdout)
				ui.ExitOnError("rendering executor", err)

			} else {
				executors, err := client.ListExecutors(namespace)
				ui.ExitOnError("listing executors: ", err)
				err = render.List(cmd, executors, os.Stdout)
				ui.ExitOnError("rendering executors", err)
			}
		},
	}

	return cmd
}
