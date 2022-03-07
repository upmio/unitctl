package sync

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	syncserver "github.com/upmio/unitctl/apps/sync/impl"
	"github.com/upmio/unitctl/apps/unit/impl"
)

var (
	serviceType                  string
	rwHostGroupId, roHostGroupId int
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "sync mysql server",
	Long:  "sync mysql server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("need to specify service group name")
		} else if len(args) > 1 {
			return fmt.Errorf("only accept one service group name")
		}

		unitClient, err := impl.NewUnit()
		if err != nil {
			return err
		}

		syncer, err := syncserver.NewSyncer(adminUser, adminPass, adminHost, "main", adminPort)
		if err != nil {
			return err
		}

		mysqlSet, err := unitClient.GetMysqlSet(context.TODO(), namespace, args[0])
		if err != nil {
			return err
		}

		switch serviceType {
		case "mysql-replication":
			return syncer.SyncServer(rwHostGroupId, roHostGroupId, args[0], mysqlSet)
		default:
			return fmt.Errorf("not support service type")
		}
	},
}

func init() {
	serverCmd.PersistentFlags().StringVarP(&serviceType, "service-type", "t", "mysql-replication", "the backend service type")
	serverCmd.PersistentFlags().IntVarP(&rwHostGroupId, "rw-hostgroup", "", 10, "the proxysql read hostgroup id")
	serverCmd.PersistentFlags().IntVarP(&roHostGroupId, "ro-hostgroup", "", 20, "the proxysql write hostgroup id")
}
