package create

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
	kindcluster "sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/log"
	"sigs.k8s.io/kind/pkg/shared/cli"
	"sigs.k8s.io/kind/pkg/shared/runtime"

	"github.com/kubeedge/keink/pkg/apis/config/defaults"
	"github.com/kubeedge/keink/pkg/cluster"
)

type flagpole struct {
	Name             string
	Config           string
	ImageName        string
	Retain           bool
	Wait             time.Duration
	Kubeconfig       string
	AdvertiseAddress string
	ContainerMode    bool
}

// NewCommand returns a new cobra.Command for cluster creation
func NewCommand(logger log.Logger, streams cmd.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "create",
		Short: "Creates one of [cluster]",
		Long:  "Creates one of local Kubernetes cluster (cluster)",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Help()
			if err != nil {
				return err
			}
			return errors.New("Subcommand is required")
		},
	}
	cmd.AddCommand(newKubeEdgeCommand(logger, streams))
	return cmd
}

// newKubeEdgeCommand returns a new cobra.Command for kubeedge creation
func newKubeEdgeCommand(logger log.Logger, streams cmd.IOStreams) *cobra.Command {
	flags := &flagpole{}

	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "kubeedge",
		Short: "Creates a local KubeEdge cluster",
		Long:  "Creates a local KubeEdge cluster using Docker container 'nodes'",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.OverrideDefaultName(cmd.Flags())
			return runE(logger, streams, flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "cluster name, overrides KIND_CLUSTER_NAME, config (default kind)")
	cmd.Flags().StringVar(&flags.Config, "config", "", "path to a kind config file")
	cmd.Flags().StringVar(&flags.ImageName, "image", defaults.Image, "node docker image to use for booting the cluster")
	cmd.Flags().BoolVar(&flags.Retain, "retain", false, "retain nodes for debugging when cluster creation fails")
	cmd.Flags().DurationVar(&flags.Wait, "wait", time.Duration(120), "wait for control plane node to be ready (default 120s)")
	cmd.Flags().StringVar(&flags.Kubeconfig, "kubeconfig", "", "sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config")
	cmd.Flags().StringVar(&flags.AdvertiseAddress, "advertise-address", "", "sets cloudcore advertise-address")
	cmd.Flags().BoolVar(&flags.ContainerMode, "container-mode", false, "sets cloudcore in container mode")

	return cmd
}

func runE(logger log.Logger, streams cmd.IOStreams, flags *flagpole) error {
	kubeedgeProvider := cluster.NewProvider(
		kindcluster.ProviderWithLogger(logger),
		runtime.GetDefault(logger),
	)

	// handle config flag, we might need to read from stdin
	withConfig, err := configOption(flags.Config, streams.In)
	if err != nil {
		return err
	}

	if err := kubeedgeProvider.CreateKubeEdge(
		flags.Name,
		withConfig,
		cluster.CreateWithKubeconfigPath(flags.Kubeconfig),
		cluster.CreateWithNodeImage(flags.ImageName),
		cluster.CreateWithWaitForReady(flags.Wait),

		// the below options are KubeEdge customized configurations
		cluster.CreateWithAdvertiseAddress(flags.AdvertiseAddress),
		cluster.CreateWithContainerMode(flags.ContainerMode),
	); err != nil {
		return fmt.Errorf("failed to create kubeedge cluster: %v", err)
	}

	return nil
}

// configOption converts the raw --config flag value to a cluster creation
// option matching it. it will read from stdin if the flag value is `-`
func configOption(rawConfigFlag string, stdin io.Reader) (cluster.CreateOption, error) {
	// if not - then we are using a real file
	if rawConfigFlag != "-" {
		return cluster.CreateWithConfigFile(rawConfigFlag), nil
	}
	// otherwise read from stdin
	raw, err := ioutil.ReadAll(stdin)
	if err != nil {
		return nil, errors.Wrap(err, "error reading config from stdin")
	}
	return cluster.CreateWithRawConfig(raw), nil
}
