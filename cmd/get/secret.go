package get

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	secret string
)

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "get secret from kubernetes cluster",
	Long:  "get secret from kubernetes cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("need to specify at least one secret name as arg")
		} else if len(args) > 1 {
			return fmt.Errorf("can only accept one secret name as arg")
		}

		ctx := context.Background()

		data, err := unitClient.GetSecret(ctx, namespace, args[0])
		if err != nil {
			return err
		}

		output, err := data.Marshal()
		if err != nil {
			return err
		}

		fmt.Println(string(output))

		return nil
	},
}
