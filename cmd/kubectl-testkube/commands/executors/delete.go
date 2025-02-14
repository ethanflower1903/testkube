package executors

import (
	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common"
	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common/validator"
	"github.com/kubeshop/testkube/pkg/ui"
	"github.com/spf13/cobra"
)

func NewDeleteExecutorCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "executor [executorName]",
		Short: "Delete Executor",
		Long:  `Delete Executor Resource, pass name to delete by name`,
		Args:  validator.ExecutorName,
		Run: func(cmd *cobra.Command, args []string) {
			name = args[0]

			client, namespace := common.GetClient(cmd)

			err := client.DeleteExecutor(name, namespace)
			ui.ExitOnError("deleting executor: "+name, err)

			ui.Success("Executor deleted")
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "unique executor name, you can also pass it as first argument")

	return cmd
}
