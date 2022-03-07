package sync

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	syncserver "github.com/upmio/unitctl/apps/sync/impl"
	"github.com/upmio/unitctl/apps/unit/impl"
	"go.uber.org/zap"
	"time"
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
		// 判断输入参数是否合法
		if len(args) == 0 {
			return fmt.Errorf("need to specify service group name")
		} else if len(args) > 1 {
			return fmt.Errorf("only accept one service group name")
		}

		var (
			logger, _   = zap.NewDevelopment()
			slogger     = logger.Sugar()
			ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		)
		defer cancel()

		// 构造unit对象以获取pod信息
		unitClient, err := impl.NewUnit(slogger)
		if err != nil {
			return err
		}

		// 查询mysql主节点ip地址
		mysqlSet, err := unitClient.GetMysqlSet(ctx, namespace, args[0])
		if err != nil {
			return err
		}

		// 此处不做是否为空判断 如果获取不到则清空proxysql mysql_server 表

		// 构造syncer对象用来同步server信息
		syncer, err := syncserver.NewSyncer(adminUser, adminPass, adminHost, "main", adminPort, slogger)
		if err != nil {
			return err
		}

		switch serviceType {
		case "mysql-replication":
			// 同步server信息
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
