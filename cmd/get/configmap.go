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
	fileDir string
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

		// 获取指定configmap 内容
		data, err := unitClient.GetConfigmap(ctx, namespace, args[0])
		if err != nil {
			return err
		}

		// 将configmap内容落到文件中
		return data.CreateConfig(fileDir)
	},
}

func init() {
	configMapCmd.PersistentFlags().StringVarP(&fileDir, "directory", "d", "/tmp", "specify directory")
}
