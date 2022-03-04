package get

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/upmio/unitctl/apps/unit"
	"github.com/upmio/unitctl/apps/unit/impl"
)

var (
	namespace  string
	unitClient unit.UnitClient
	err        error
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
	unitClient, err = impl.NewUnitImpl()
	if err != nil {
		panic(err)
	}
	GetCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "specify namespace")
	GetCmd.AddCommand(secretCmd)
	GetCmd.AddCommand(configMapCmd)
}
