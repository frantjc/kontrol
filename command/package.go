package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/frantjc/go-fn"
	"github.com/frantjc/kontrol"
	apiv1alpha1 "github.com/frantjc/kontrol/api/v1alpha1"
	"github.com/frantjc/kontrol/cr"
	"github.com/frantjc/kontrol/pkg"
	"github.com/frantjc/kontrol/pkg/lbl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func NewPackage() *cobra.Command {
	var (
		conf   string
		config = viper.New()
		rules  []string
		cmd    = &cobra.Command{
			Use:           "package",
			SilenceErrors: true,
			SilenceUsage:  true,
			Args:          cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx                       = cmd.Context()
					_                         = kontrol.LoggerFrom(ctx)
					packager                  = new(lbl.Packager)
					ref                       = args[0]
					policyRules               = []rbacv1.PolicyRule{}
					customResourceDefinitions = []apiextensionsv1.CustomResourceDefinition{}
					ext                       = filepath.Ext(conf)
					stdinUsed                 = false
				)

				config.SetConfigName(strings.TrimSuffix(conf, ext))
				config.SetConfigType(strings.TrimPrefix(ext, "."))

				if err := config.ReadInConfig(); err != nil {
					return err
				}

				for _, role := range config.GetStringSlice("package.roles") {
					var (
						r   = cmd.InOrStdin()
						err error
					)

					if role != "-" {
						r, err = os.Open(role)
						if err != nil {
							return err
						}
					} else {
						if stdinUsed {
							return fmt.Errorf("cannot use stdin twice")
						}

						stdinUsed = true
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

				for _, crd := range config.GetStringSlice("package.crds") {
					var (
						r   = cmd.InOrStdin()
						err error
					)

					if crd != "-" {
						r, err = os.Open(crd)
						if err != nil {
							return err
						}
					} else {
						if stdinUsed {
							return fmt.Errorf("cannot use stdin twice")
						}

						stdinUsed = true
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
						ContainerPorts: fn.Map(config.GetIntSlice("port"), func(port int, _ int) int32 {
							return int32(port)
						}),
						Command: strings.Split(config.GetString("package.cmd"), " "),
						Args:    strings.Split(config.GetString("package.args"), " "),
					},
				})
				if err != nil {
					return err
				}

				return cr.Push(ctx, ref, image)
			},
		}
	)

	config.AddConfigPath(".")

	cmd.Flags().StringVarP(&conf, "conf", "c", ".kontrol.yml", "Config file")
	_ = cmd.MarkFlagFilename("conf", "yaml", "yml")

	cmd.Flags().StringSlice("crds", nil, "CRDs for the controller")
	_ = cmd.MarkFlagFilename("crds", "yaml", "yml")
	_ = config.BindPFlag("package.crds", cmd.Flag("crds"))

	cmd.Flags().StringSlice("roles", nil, "Roles for the controller")
	_ = cmd.MarkFlagFilename("roles", "yaml", "yml")
	_ = config.BindPFlag("package.roles", cmd.Flag("crds"))

	cmd.Flags().Int32SliceP("port", "p", nil, "Ports that the controller binds to")
	_ = config.BindPFlag("package.ports", cmd.Flag("port"))

	cmd.Flags().String("cmd", "", "Command for the controller")
	_ = config.BindPFlag("package.cmd", cmd.Flag("cmd"))

	cmd.Flags().String("args", "", "Args for the controller")
	_ = config.BindPFlag("package.args", cmd.Flag("args"))

	cmd.Flags().StringSliceVar(&rules, "rule", nil, "Extra rules for the controller")

	return cmd
}
