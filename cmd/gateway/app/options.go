package app

import (
	"github.com/spf13/cobra"
)

type Options struct {
	TriggerService string
	Port           int
	Version        bool
}

func (s *Options) SetOps(ac *cobra.Command) {
	ac.Flags().StringVar(&s.TriggerService, "trigger-service", s.TriggerService, "tekton trigger serivice name")
	ac.Flags().IntVar(&s.Port, "port", 8080, "port")
	ac.Flags().BoolVar(&s.Version, "version", s.Version, "Print version information")
}
