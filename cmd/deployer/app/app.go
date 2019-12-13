package app

import (
	"encoding/json"

	"log"

	deployer "github.com/knative-sample/tekton-triggers-tools/pkg/deployer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
)

var defaultZLC = []byte(`{
  "level": "info",
  "development": false,
  "outputPaths": ["stdout"],
  "errorOutputPaths": ["stderr"],
  "encoding": "json",
  "encoderConfig": {
    "timeKey": "ts",
    "levelKey": "level",
    "nameKey": "logger",
    "callerKey": "caller",
    "messageKey": "msg",
    "stacktraceKey": "stacktrace",
    "lineEnding": "",
    "levelEncoder": "",
    "timeEncoder": "iso8601",
    "durationEncoder": "",
    "callerEncoder": ""
  }
}`)

// start edas api
func NewCommandStartServer() *cobra.Command {
	ops := &Options{}
	mainCmd := &cobra.Command{
		Short: "deployer",
		Long:  "deployer",
		RunE: func(c *cobra.Command, args []string) error {
			run(ops)
			return nil
		},
	}

	ops.SetOps(mainCmd)
	return mainCmd
}

func run(ops *Options) {
	var logger *zap.SugaredLogger
	// setup logger
	var zlc zap.Config
	if err := json.Unmarshal(defaultZLC, &zlc); err != nil {
		log.Fatalf("Unmarshal zap.Logger config error:%s ", err)
	}
	if l, err := zlc.Build(); err != nil {
		log.Fatalf("Build Logger error:%s", err)
	} else {
		logger = l.Sugar()
	}
	defer logger.Sync()

	if ops.KsvName == "" {
		logger.Fatalf("--ksvcname is empty")
	}
	if ops.Namespace == "" {
		logger.Fatalf("--namespace is empty")
	}
	if ops.Image == "" {
		logger.Fatalf("--image is empty")
	}

	ctx := signals.NewContext()
	ctx = logging.WithLogger(ctx, logger)
	logger.Info("logger construction succeeded")

	// setup informer
	cfg, err := sharedmain.GetConfig(ops.MasterURL, ops.Kubeconfig)
	if err != nil {
		logger.Fatal("Error building kubeconfig:", err)
	}

	logger.Infof("Registering %d clients", len(injection.Default.GetClients()))
	logger.Infof("Registering %d informer factories", len(injection.Default.GetInformerFactories()))
	logger.Infof("Registering %d informers", len(injection.Default.GetInformers()))

	ctx, _ = injection.Default.SetupInformers(ctx, cfg)

	dployer := deployer.NewDeployer(ctx)
	dployer.Deploy(ops.KsvName, ops.Namespace, ops.Image)
}
