package sync

import (
	"context"
	"fmt"
	syncserver "github.com/upmio/unitctl/apps/sync/impl"
	"github.com/upmio/unitctl/apps/unit/impl"
	"os"

	"github.com/spf13/cobra"
	mysql "github.com/upmio/unitctl/apps/mysql/impl"
)

var (
	syncUser, syncPass               string
	defaultHostGroup, maxConnections int
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "sync mysql user",
	Long:  "sync mysql user",
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

		if len(mysqlSet) == 0 {
			return fmt.Errorf("can't get mysql master by service group name")
		}

		mysqlClient, err := mysql.NewMysqlClient(syncUser, syncPass, mysqlSet[0].IpAddr, "mysql", mysqlSet[0].Port)
		if err != nil {
			return err
		}

		hostIp := os.Getenv("INTERNAL_IP")
		userSet, err := mysqlClient.GetMysqlUser(hostIp)
		if err != nil {
			return err
		}

		return syncer.SyncUsers(defaultHostGroup, maxConnections, userSet)
	},
}

func init() {
}
