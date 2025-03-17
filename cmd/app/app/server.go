package app

import (
	"github.com/spf13/cobra"

	"dp-command/cmd/app/app/options"
)

func NewServerCommand() *cobra.Command {
	o := options.NewAppOptions()
	cmd := &cobra.Command{
		Use:  "app",
		Long: `Long describe.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(o)
		},
	}
	// 指定配置参数
	cmd.Flags().StringVarP(&o.ConfFile, "config", "c", "", "Config file path.")
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of dp-command.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("command")
		},
	}

	cmd.AddCommand(versionCmd)

	return cmd
}

func run(o *options.AppOptions) (err error) {
	server, err := o.NewServer()
	server.Run()
	return
}
