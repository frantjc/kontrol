package command

import (
	"github.com/frantjc/kontrol"
	apiv1alpha1 "github.com/frantjc/kontrol/api/v1alpha1"
	"github.com/frantjc/kontrol/controllers"

	// cloud auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/spf13/cobra"
)

func NewKontroller() *cobra.Command {
	var (
		schm                   = runtime.NewScheme()
		leaderElect            bool
		metricsAddr, probeAddr string
		cmd                    = &cobra.Command{
			Use:           "kontroller",
			Version:       kontrol.GetSemver(),
			SilenceErrors: true,
			SilenceUsage:  true,
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx = cmd.Context()
					log = kontrol.LoggerFrom(ctx)
				)

				if err := scheme.AddToScheme(schm); err != nil {
					return err
				}

				if err := apiv1alpha1.AddToScheme(schm); err != nil {
					return err
				}

				cfg, err := ctrl.GetConfig()
				if err != nil {
					return err
				}

				mgr, err := ctrl.NewManager(cfg, ctrl.Options{
					Scheme:                 schm,
					MetricsBindAddress:     metricsAddr,
					Port:                   8080,
					HealthProbeBindAddress: probeAddr,
					LeaderElection:         leaderElect,
					LeaderElectionID:       "15aa2f8a.frantj.cc",
					Logger:                 log,
				})
				if err != nil {
					return err
				}

				if err = (&controllers.KontrollerReconciler{
					Client: mgr.GetClient(),
					Scheme: mgr.GetScheme(),
					Logger: mgr.GetLogger(),
				}).SetupWithManager(mgr); err != nil {
					return err
				}

				if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
					return err
				}

				if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
					return err
				}

				return mgr.Start(ctx)
			},
		}
	)

	cmd.Flags().String("kubeconfig", "", "kubeconfig")
	cmd.Flags().StringVar(&metricsAddr, "metrics-bind-address", ":8081", "address the metric endpoint binds to")
	cmd.Flags().StringVar(&probeAddr, "health-probe-bind-address", ":8082", "address the probe endpoint binds to")
	cmd.Flags().BoolVar(&leaderElect, "leader-elect", false, "enable leader election for controller manager")

	return cmd
}
