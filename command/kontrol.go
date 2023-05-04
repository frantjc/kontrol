package command

import (
	"github.com/spf13/cobra"
)

func NewKontrol() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kontrol",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(NewPackage())
	cmd.AddCommand(NewDeploy())

	return cmd
}
