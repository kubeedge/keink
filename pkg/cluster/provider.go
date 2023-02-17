package cluster

import (
	"fmt"

	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"sigs.k8s.io/kind/pkg/cluster"
	sharedcreate "sigs.k8s.io/kind/pkg/cluster/shared/create"
	shareddocker "sigs.k8s.io/kind/pkg/cluster/shared/providers/docker"
	"sigs.k8s.io/kind/pkg/shared/apis/config"

	"github.com/kubeedge/keink/pkg/apis/config/defaults"
	internalcreate "github.com/kubeedge/keink/pkg/cluster/internal/create"
)

type Provider struct {
	*cluster.ProviderWrapper
}

func NewProvider(options ...cluster.ProviderOption) *Provider {
	p := cluster.NewProviderWrapper(options...)
	return &Provider{
		ProviderWrapper: p,
	}

}

func (p *Provider) CreateKubeEdge(name string, options ...CreateOption) error {
	// apply options
	opts := &internalcreate.ClusterOptions{
		ClusterOptions: sharedcreate.ClusterOptions{
			NameOverride: name,
		},
	}

	for _, o := range options {
		if err := o.apply(opts); err != nil {
			return err
		}
	}

	PreProcessClusterOptions(opts)

	// create k8s cluster using kind library directly.
	err := sharedcreate.Cluster(p.Logger, p.Provider, &opts.ClusterOptions)
	if err != nil {
		return fmt.Errorf("failed to create k8s cluster: %v", err)
	}

	// create kubeedge cluster
	return internalcreate.Cluster(p.Logger, p.Provider, opts)
}

// PreProcessClusterOptions do some pre-processing on ClusterOptions so that kind api can recognize it
// will overwrite the input argument directly
func PreProcessClusterOptions(opts *internalcreate.ClusterOptions) {
	// we must ensure there's at least one Edge Node
	// so create an Edge Node if it not exist
	edgeNodeExist := false
	for _, node := range opts.ClusterOptions.Config.Nodes {
		if string(node.Role) == string(defaults.EdgeNodeRole) {
			edgeNodeExist = true
			break
		}
	}
	if !edgeNodeExist {
		opts.ClusterOptions.Config.Nodes = append(opts.ClusterOptions.Config.Nodes, config.Node{
			Role: config.WorkerRole,
			Labels: map[string]string{
				// use this label to indicate it is a KubeEdge edge node
				// this label will be applied to kubelet, and we can use `kubectl get node` to get this label info
				// so that we can run this command to predict whether a node belong to k8s worker node or KubeEdge edge node
				shareddocker.EdgeNodeLabelKey: shareddocker.EdgeNodeLabelValue,
			},
		})
	}

	// Due to kind only can recognize node role: control-plane and worker, but not edge-node
	// so we need to convert edge-node to worker before create k8s cluster
	for index := range opts.ClusterOptions.Config.Nodes {
		if string(opts.ClusterOptions.Config.Nodes[index].Role) == string(defaults.EdgeNodeRole) {
			// convert node role from edge-node to worker
			opts.ClusterOptions.Config.Nodes[index].Role = config.NodeRole(v1alpha4.WorkerRole)

			// apply label to indicate it is a KubeEdge edge node
			if len(opts.ClusterOptions.Config.Nodes[index].Labels) == 0 {
				opts.ClusterOptions.Config.Nodes[index].Labels = make(map[string]string)
			}
			opts.ClusterOptions.Config.Nodes[index].Labels[shareddocker.EdgeNodeLabelKey] = shareddocker.EdgeNodeLabelValue
			break
		}
	}

}
