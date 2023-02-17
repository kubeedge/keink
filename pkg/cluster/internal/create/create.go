package create

import (
	"sigs.k8s.io/kind/pkg/cluster/shared/create"
	"sigs.k8s.io/kind/pkg/cluster/shared/create/actions"
	"sigs.k8s.io/kind/pkg/cluster/shared/delete"
	"sigs.k8s.io/kind/pkg/cluster/shared/providers"
	"sigs.k8s.io/kind/pkg/log"
	"sigs.k8s.io/kind/pkg/shared/cli"

	"github.com/kubeedge/keink/pkg/cluster/internal/create/actions/kubeedge"
)

// ClusterOptions wraps kind cluster creation options with KubeEdge customized options
type ClusterOptions struct {
	create.ClusterOptions

	AdvertiseAddress string
	ContainerMode    bool
}

// Cluster creates a cluster
func Cluster(logger log.Logger, p providers.Provider, opts *ClusterOptions) error {
	// setup a status object to show progress to the user
	status := cli.StatusForLogger(logger)

	actionsToRun := []actions.Action{
		kubeedge.NewAction(opts.AdvertiseAddress, opts.ContainerMode), // run kubeedge install
	}

	actionsContext := actions.NewActionContext(logger, status, p, opts.Config)

	for _, action := range actionsToRun {
		if err := action.Execute(actionsContext); err != nil {
			if !opts.Retain {
				_ = delete.Cluster(logger, p, opts.Config.Name, opts.KubeconfigPath)
			}
			return err
		}
	}

	return nil
}
