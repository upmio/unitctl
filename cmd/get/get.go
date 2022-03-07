package get

import (
	"errors"
	"github.com/spf13/cobra"
)

var (
	namespace string
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "get resource from kubernetes cluster",
	Long:  "Display one or many resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("no command find")
	},
}

func init() {
	GetCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "specify namespace")
	GetCmd.AddCommand(secretCmd)
	GetCmd.AddCommand(configMapCmd)
}
