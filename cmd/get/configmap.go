package get

import (
	"context"
	"fmt"
	"github.com/upmio/unitctl/apps/unit/impl"

	"github.com/spf13/cobra"
)

var (
	fileDir, configMapName string
)

var configMapCmd = &cobra.Command{
	Use:   "configmap",
	Short: "get configmap from kubernetes cluster",
	Long:  "get configmap from kubernetes cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("need to specify at least one configmap name as arg")
		} else if len(args) > 1 {
			return fmt.Errorf("can only accept one configmap name as arg")
		}

		ctx := context.Background()

		unitClient, err := impl.NewUnit()
		if err != nil {
			return err
		}

		data, err := unitClient.GetConfigmap(ctx, namespace, args[0])
		if err != nil {
			return err
		}

		return data.CreateConfig(fileDir)
	},
}

func init() {
	configMapCmd.PersistentFlags().StringVarP(&fileDir, "directory", "d", "/tmp", "specify directory")
}
