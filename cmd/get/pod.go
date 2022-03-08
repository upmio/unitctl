package get

import (
	"context"
	"fmt"
	"github.com/upmio/unitctl/apps/unit/impl"
	"go.uber.org/zap"
	"time"

	"github.com/spf13/cobra"
)

var (
	showLabels bool
)

var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "get pod dbscale label from kubernetes cluster",
	Long:  "get pod dbscale label from kubernetes cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("need to specify at least one pod name as arg")
		} else if len(args) > 1 {
			return fmt.Errorf("can only accept one pod name as arg")
		}

		var (
			logger, _   = zap.NewDevelopment()
			slogger     = logger.Sugar()
			ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		)
		defer cancel()

		unitClient, err := impl.NewUnit(slogger)
		if err != nil {
			return err
		}

		if showLabels {
			// 获取指定 podName 内容
			labelSet, err := unitClient.GetPodLabelSet(ctx, namespace, args[0])
			if err != nil {
				return err
			}
			output, err := labelSet.Marshal()
			if err != nil {
				return err
			}

			fmt.Println(string(output))
		} else {
			fmt.Println(args[0])
		}
		return nil
	},
}

func init() {
	podCmd.PersistentFlags().BoolVar(&showLabels, "show-labels", false, "When printing, show all labels as the last column (default hide labels column)")
}
