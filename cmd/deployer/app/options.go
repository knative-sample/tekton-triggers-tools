package app

import (
	"github.com/spf13/cobra"
)

type Options struct {
	Version    bool
	MasterURL  string
	Kubeconfig string
	Namespace  string
	KsvName    string
	Image      string
}

func (s *Options) SetOps(ac *cobra.Command) {
	ac.Flags().BoolVar(&s.Version, "version", s.Version, "Print version information")
	ac.Flags().StringVar(&s.MasterURL, "master", s.MasterURL, "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	ac.Flags().StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig, "Path to a kubeconfig. Only required if out-of-cluster.")
	ac.Flags().StringVar(&s.Namespace, "namespace", "default", "deployer knative service namespace")
	ac.Flags().StringVar(&s.KsvName, "ksvcname", "", "deployer knative service name")
	ac.Flags().StringVar(&s.Image, "image", "", "deployer knative service image")
}
