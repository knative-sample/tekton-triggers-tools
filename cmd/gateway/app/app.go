package app

import (
	"strings"

	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/knative-sample/tekton-triggers-tools/pkg/gateway"
	"github.com/knative-sample/tekton-triggers-tools/pkg/version"
	"github.com/spf13/cobra"
)

// start edas api
func NewCommandStartServer(stopCh <-chan struct{}) *cobra.Command {
	ops := &Options{}
	mainCmd := &cobra.Command{
		Short: "hello world runner",
		Long:  "hello world runner",
		RunE: func(c *cobra.Command, args []string) error {
			glog.V(2).Infof("NewCommandStartServer main:%s", strings.Join(args, " "))
			run(stopCh, ops)
			return nil
		},
	}

	ops.SetOps(mainCmd)
	return mainCmd
}

func run(stopCh <-chan struct{}, ops *Options) {
	vs := version.Version().Info("trigger-proxy")
	if ops.Version {
		fmt.Println(vs)
		os.Exit(0)
	}

	if ops.TriggerService == "" {
		glog.Fatalf("--trigger-service is empty")
	}

	startApiArgs := &gateway.Proxy{
		Port:           ops.Port,
		TriggerService: ops.TriggerService,
	}
	gateway.StartProxy(startApiArgs)

	<-stopCh
}
