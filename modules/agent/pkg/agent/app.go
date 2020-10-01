package agent

import (
	"context"

	"github.com/rancher/fleet/modules/agent/pkg/controllers"
	"github.com/rancher/fleet/modules/agent/pkg/register"
	"github.com/rancher/wrangler/pkg/kubeconfig"
	"github.com/rancher/wrangler/pkg/ratelimit"
)

type Options struct {
	DefaultNamespace string
	ClusterID        string
	NoLeaderElect    bool
	CheckinInterval  time.Duration
}

func Register(ctx context.Context, kubeConfig, namespace, clusterID string) error {
	clientConfig := kubeconfig.GetNonInteractiveClientConfig(kubeConfig)
	kc, err := clientConfig.ClientConfig()
	if err != nil {
		return err
	}
	kc.RateLimiter = ratelimit.None

	_, err = register.Register(ctx, namespace, clusterID, kc)
	return err
}

func Start(ctx context.Context, kubeConfig, namespace string, opts *Options) error {
	if opts == nil {
		opts = &Options{}
	}
	if opts.DefaultNamespace == "" {
		opts.DefaultNamespace = "default"
	}

	clientConfig := kubeconfig.GetNonInteractiveClientConfig(kubeConfig)
	kc, err := clientConfig.ClientConfig()
	if err != nil {
		return err
	}

	agentInfo, err := register.Register(ctx, namespace, opts.ClusterID, kc)
	if err != nil {
		return err
	}

	fleetNamespace, _, err := agentInfo.ClientConfig.Namespace()
	if err != nil {
		return err
	}

	fleetRestConfig, err := agentInfo.ClientConfig.ClientConfig()
	if err != nil {
		return err
	}

	return controllers.Register(ctx,
		!opts.NoLeaderElect,
		fleetNamespace,
		namespace,
		opts.DefaultNamespace,
		agentInfo.ClusterNamespace,
		agentInfo.ClusterName,
		opts.CheckinInterval,
		fleetRestConfig,
		clientConfig)
}
