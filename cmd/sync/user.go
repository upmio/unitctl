package sync

import (
	"context"
	"fmt"
	syncserver "github.com/upmio/unitctl/apps/sync/impl"
	"github.com/upmio/unitctl/apps/unit/impl"
	"go.uber.org/zap"
	"os"
	"time"

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
		// 判断输入参数是否合法
		if len(args) == 0 {
			return fmt.Errorf("need to specify service group name")
		} else if len(args) > 1 {
			return fmt.Errorf("only accept one service group name")
		}

		//初始化logger和context对象
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

		// 判断是否能获取到mysql master及诶点信息
		if len(mysqlSet) == 0 {
			return fmt.Errorf("can't get mysql master by service group name")
		}

		// 构造mysqlclient对象用来查询所有用户信息
		mysqlClient, err := mysql.NewMysqlClient(syncUser, syncPass, mysqlSet[0].IpAddr, "mysql", mysqlSet[0].Port, slogger)
		if err != nil {
			return err
		}

		// 获取mysql 用户信息
		hostIp := os.Getenv("INTERNAL_IP")
		userSet, err := mysqlClient.GetMysqlUser(hostIp)
		if err != nil {
			return err
		}

		// 构造syncer对象用来同步user信息
		syncer, err := syncserver.NewSyncer(adminUser, adminPass, adminHost, "main", adminPort, slogger)
		if err != nil {
			return err
		}

		// 同步用户
		return syncer.SyncUsers(defaultHostGroup, maxConnections, userSet)
	},
}

func init() {
	userCmd.PersistentFlags().StringVarP(&syncUser, "sync-username", "", "root", "sync username")
	userCmd.PersistentFlags().StringVarP(&syncPass, "sync-password", "", "", "sync password")
	userCmd.PersistentFlags().IntVarP(&maxConnections, "max-connection", "", 1024, "proxysql user max connection")
	userCmd.PersistentFlags().IntVarP(&defaultHostGroup, "default-hostgroup", "", 10, "proxysql user default hostgroup id")
}
