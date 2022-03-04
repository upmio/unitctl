package cmd

import (
	"errors"
	"fmt"
	"github.com/upmio/unitctl/cmd/get"
	"os"

	"github.com/spf13/cobra"

	"github.com/upmio/unitctl/version"
)

var (
	vers, helpFlag bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "unitctl",
	Short: "unitctl get service group configmap or secret from Kubernetes cluster.",
	Long: `  This tool used for unit to get configmap or secret immediately in kubernetes which
can solve the problem that configmap or secret will not update immediately in pod.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(version.FullVersion())
			return nil
		}
		return errors.New("no flags find")
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&helpFlag, "help", "", false, "Help default flag")
	RootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, "the proxysql-initializer version")
	RootCmd.AddCommand(get.GetCmd)
}
