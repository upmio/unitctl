package sync

import (
	"errors"
	"github.com/spf13/cobra"
)

var (
	helpFlag                                   bool
	namespace, adminUser, adminPass, adminHost string
	adminPort                                  int
	err                                        error
)

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync resource into proxysql admin interface",
	Long:  "Sync server or user to proxysql",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("no command find")
	},
}

func init() {
	SyncCmd.PersistentFlags().BoolVarP(&helpFlag, "help", "", false, "Help default flag")
	SyncCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "specify namespace")
	SyncCmd.PersistentFlags().StringVarP(&adminUser, "admin-username", "u", "admin", "proxysql admin user")
	SyncCmd.PersistentFlags().StringVarP(&adminPass, "admin-password", "p", "", "proxysql admin password")
	SyncCmd.PersistentFlags().StringVarP(&adminHost, "admin-host", "h", "127.0.0.1", "proxysql admin host")
	SyncCmd.PersistentFlags().IntVarP(&adminPort, "admin-port", "P", 6032, "proxysql admin port")
	SyncCmd.AddCommand(serverCmd)
	SyncCmd.AddCommand(userCmd)
}
