package command

import (
	"runtime"

	"github.com/frantjc/kontrol"
	"github.com/spf13/cobra"
)

func NewKontrol() *cobra.Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:           "kontrol",
			Version:       kontrol.GetSemver(),
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				cmd.SetContext(
					kontrol.WithLogger(cmd.Context(), kontrol.NewLogger().V(2-verbosity)),
				)
			},
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.PersistentFlags().CountVarP(&verbosity, "verbose", "V", "verbosity for kontrol")
	cmd.AddCommand(NewPackage(), NewDeploy())

	return cmd
}
