package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	apiv1alpha1 "github.com/frantjc/kontrol/api/v1alpha1"
	"github.com/frantjc/kontrol/cr"
	"github.com/frantjc/kontrol/pkg"
	"github.com/frantjc/kontrol/pkg/lbl"
	"github.com/spf13/cobra"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func NewPackage() *cobra.Command {
	var (
		containerPorts []int32
		rules          []string
		crds           string
		roles          string
		argsStr        string
		commandStr     string
		cmd            = &cobra.Command{
			Use:           "package",
			SilenceErrors: true,
			SilenceUsage:  true,
			Args:          cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx                       = cmd.Context()
					packager                  = new(lbl.Packager)
					ref                       = args[0]
					policyRules               = []rbacv1.PolicyRule{}
					customResourceDefinitions = []apiextensionsv1.CustomResourceDefinition{}
				)

				if cmd.Flag("roles").Changed {
					var (
						r   = cmd.InOrStdin()
						err error
					)

					if roles != "-" {
						r, err = os.Open(roles)
						if err != nil {
							return err
						}
					}

					decoder := yaml.NewYAMLOrJSONDecoder(r, 0)
					for {
						role := struct {
							PolicyRules []rbacv1.PolicyRule `json:"rules,omitempty"`
						}{}
						if err := decoder.Decode(&role); errors.Is(err, io.EOF) {
							break
						} else if err != nil {
							return err
						}

						policyRules = append(policyRules, role.PolicyRules...)
					}
				}

				if cmd.Flag("crds").Changed {
					var (
						r   = cmd.InOrStdin()
						err error
					)

					if crds != "-" {
						r, err = os.Open(crds)
						if err != nil {
							return err
						}
					}

					decoder := yaml.NewYAMLOrJSONDecoder(r, 0)
					for {
						crd := apiextensionsv1.CustomResourceDefinition{}
						if err := decoder.Decode(&crd); errors.Is(err, io.EOF) {
							break
						} else if err != nil {
							return err
						}

						customResourceDefinitions = append(customResourceDefinitions, crd)
					}
				}

				for _, rule := range rules {
					parts := strings.Split(rule, ":")
					switch len(parts) {
					case 2:
						policyRules = append(
							policyRules,
							rbacv1.PolicyRule{
								APIGroups: []string{""},
								Resources: strings.Split(parts[0], ","),
								Verbs:     strings.Split(parts[1], ","),
							},
						)
					case 3:
						policyRules = append(
							policyRules,
							rbacv1.PolicyRule{
								APIGroups: strings.Split(parts[0], ","),
								Resources: strings.Split(parts[1], ","),
								Verbs:     strings.Split(parts[2], ","),
							},
						)
					default:
						return fmt.Errorf("unparsable rule %s", rule)
					}
				}

				image, err := cr.Pull(ctx, ref)
				if err != nil {
					return err
				}

				image, err = packager.Package(ctx, image, &pkg.Packaging{
					PolicyRules:               policyRules,
					CustomResourceDefinitions: customResourceDefinitions,
					Overrides: apiv1alpha1.Overrides{
						ContainerPorts: containerPorts,
						Command:        strings.Split(commandStr, " "),
						Args:           strings.Split(argsStr, " "),
					},
				})
				if err != nil {
					return err
				}

				return cr.Push(ctx, ref, image)
			},
		}
	)

	cmd.Flags().StringVar(&crds, "crds", "", "CRDs for the controller")
	_ = cmd.MarkFlagFilename("crds", "yaml", "yml")
	cmd.Flags().StringVar(&roles, "roles", "", "Roles for the controller")
	_ = cmd.MarkFlagFilename("role", "yaml", "yml")
	cmd.Flags().StringSliceVar(&rules, "rule", nil, "Extra rules for the controller")
	cmd.Flags().Int32SliceVarP(&containerPorts, "port", "p", nil, "Ports that the controller binds to")
	cmd.Flags().StringVar(&commandStr, "cmd", "", "Command for the controller")
	cmd.Flags().StringVar(&argsStr, "args", "", "Args for the controller")

	return cmd
}
