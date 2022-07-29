package build

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/log"

	"github.com/kubeedge/keink/pkg/build/edgeimage"
)

type flagpole struct {
	Source    string
	BuildType string
	Image     string
	BaseImage string
	KubeRoot  string
	Arch      string
}

// NewCommand returns a new cobra.Command for building
func NewCommand(logger log.Logger, streams cmd.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.NoArgs,
		// TODO(bentheelder): more detailed usage
		Use:   "build",
		Short: "Build one of [edge-image]",
		Long:  "Build one of [edge-image]",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Help()
			if err != nil {
				return err
			}
			return errors.New("Subcommand is required")
		},
	}
	// add subcommands
	cmd.AddCommand(newBuildCommand(logger, streams))
	return cmd
}

// NewCommand returns a new cobra.Command for building the node image
func newBuildCommand(logger log.Logger, streams cmd.IOStreams) *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args: cobra.MaximumNArgs(1),
		// TODO(bentheelder): more detailed usage
		Use:   "edge-image [kubeedge-source]",
		Short: "Build the edge image",
		Long:  "Build the edge image which contains KubeEdge build artifacts and other kind requirements",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Lookup("kube-root").Changed {
				if len(args) != 0 {
					return errors.New("passing an argument and deprecated --kube-root is not supported, please switch to just the argument")
				}
				logger.Warn("--kube-root is deprecated, please switch to passing this as an argument")
			}
			if cmd.Flags().Lookup("type").Changed {
				return errors.New("--type is no longer supported, please remove this flag")
			}
			return runE(logger, flags, args)
		},
	}
	cmd.Flags().StringVar(
		&flags.BuildType, "type",
		"docker", "build type, default is docker",
	)
	cmd.Flags().StringVar(
		&flags.Image, "image",
		edgeimage.DefaultImage,
		"name:tag of the resulting image to be built",
	)
	cmd.Flags().StringVar(
		&flags.KubeRoot, "kube-root",
		"",
		"path to the KubeEdge source directory (if empty, the path is autodetected)",
	)
	cmd.Flags().StringVar(
		&flags.BaseImage, "base-image",
		edgeimage.DefaultBaseImage,
		"name:tag of the base image to use for the build",
	)
	cmd.Flags().StringVar(
		&flags.Arch, "arch",
		"",
		"architecture to build for, defaults to the host architecture",
	)
	return cmd
}

func runE(logger log.Logger, flags *flagpole, args []string) error {
	kubeRoot := flags.KubeRoot
	if len(args) > 0 {
		kubeRoot = args[0]
	}
	if err := edgeimage.Build(
		edgeimage.WithImage(flags.Image),
		edgeimage.WithBaseImage(flags.BaseImage),
		edgeimage.WithKubeEdgeRoot(kubeRoot),
		edgeimage.WithLogger(logger),
		edgeimage.WithArch(flags.Arch),
	); err != nil {
		return fmt.Errorf("failed to build edge-image: %v", err)
	}
	return nil
}
