package cluster

import (
	"time"

	"sigs.k8s.io/kind/pkg/shared/apis/config/encoding"

	internalcreate "github.com/kubeedge/keink/pkg/cluster/internal/create"
)

// CreateOption is a Provider.Create option
type CreateOption interface {
	apply(*internalcreate.ClusterOptions) error
}

type createOptionAdapter func(*internalcreate.ClusterOptions) error

func (c createOptionAdapter) apply(o *internalcreate.ClusterOptions) error {
	return c(o)
}

// TODO: add more options

// CreateWithKubeconfigPath sets the explicit --kubeconfig path
func CreateWithKubeconfigPath(explicitPath string) CreateOption {
	return createOptionAdapter(func(o *internalcreate.ClusterOptions) error {
		o.KubeconfigPath = explicitPath
		return nil
	})
}
func CreateWithNodeImage(nodeImage string) CreateOption {
	return createOptionAdapter(func(o *internalcreate.ClusterOptions) error {
		o.NodeImage = nodeImage
		return nil
	})
}

func CreateWithWaitForReady(waitTime time.Duration) CreateOption {
	return createOptionAdapter(func(o *internalcreate.ClusterOptions) error {
		o.WaitForReady = waitTime
		return nil
	})
}

// CreateWithConfigFile configures the config file path to use
func CreateWithConfigFile(path string) CreateOption {
	return createOptionAdapter(func(o *internalcreate.ClusterOptions) error {
		var err error
		o.Config, err = encoding.Load(path)
		return err
	})
}

// CreateWithRawConfig configures the config to use from raw (yaml) bytes
func CreateWithRawConfig(raw []byte) CreateOption {
	return createOptionAdapter(func(o *internalcreate.ClusterOptions) error {
		var err error
		o.Config, err = encoding.Parse(raw)
		return err
	})
}

// CreateWithAdvertiseAddress sets the explicit --advertise-address ip
func CreateWithAdvertiseAddress(address string) CreateOption {
	return createOptionAdapter(func(o *internalcreate.ClusterOptions) error {
		o.AdvertiseAddress = address
		return nil
	})
}

// CreateWithContainerMode sets the explicit --container-mode
func CreateWithContainerMode(containerMode bool) CreateOption {
	return createOptionAdapter(func(o *internalcreate.ClusterOptions) error {
		o.ContainerMode = containerMode
		return nil
	})
}
