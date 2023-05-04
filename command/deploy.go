package command

import (
	"github.com/frantjc/kontrol/cr"
	"github.com/frantjc/kontrol/pkg"
	"github.com/frantjc/kontrol/pkg/lbl"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/printers"
)

func NewDeploy() *cobra.Command {
	var (
		name      string
		namespace string
		cmd       = &cobra.Command{
			Use:           "deploy",
			SilenceErrors: true,
			SilenceUsage:  true,
			Args:          cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx      = cmd.Context()
					packager = new(lbl.Packager)
					printer  = new(printers.YAMLPrinter)
					ref      = args[0]
				)

				image, err := cr.Pull(ctx, ref)
				if err != nil {
					return err
				}

				packaging, err := packager.Unpackage(ctx, image)
				if err != nil {
					return err
				}

				for _, o := range pkg.PackagingToObjects(ref, name, namespace, packaging) {
					if err = printer.PrintObj(o, cmd.OutOrStdout()); err != nil {
						return err
					}
				}

				return nil
			},
		}
	)

	cmd.Flags().StringVar(&name, "name", "kontroller", "Name for the controller")
	cmd.Flags().StringVar(&namespace, "namespace", "kontroller", "Namespace for the controller")

	return cmd
}
