package cmd

import (
	"io"
	"io/ioutil"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/cmd/kind/completion"
	"sigs.k8s.io/kind/pkg/cmd/kind/delete"
	"sigs.k8s.io/kind/pkg/cmd/kind/export"
	"sigs.k8s.io/kind/pkg/cmd/kind/get"
	"sigs.k8s.io/kind/pkg/cmd/kind/load"
	"sigs.k8s.io/kind/pkg/cmd/kind/version"
	"sigs.k8s.io/kind/pkg/log"

	"github.com/kubeedge/keink/pkg/cmd/build"
	"github.com/kubeedge/keink/pkg/cmd/create"
)

type flagpole struct {
	LogLevel  string
	Verbosity int32
	Quiet     bool
}

// NewCommand returns a new cobra.Command implementing the root command for kind
func NewCommand(logger log.Logger, streams cmd.IOStreams) *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "keink",
		Short: "keink is a tool for managing local Kubernetes and KubeEdge clusters",
		Long:  "keink creates and manages local Kubernetes and KubeEdge clusters using Docker container 'nodes'",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return runE(logger, flags, cmd)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.Version(),
	}

	cmd.SetOut(streams.Out)
	cmd.SetErr(streams.ErrOut)
	cmd.PersistentFlags().StringVar(
		&flags.LogLevel,
		"loglevel",
		"",
		"DEPRECATED: see -v instead",
	)
	cmd.PersistentFlags().Int32VarP(
		&flags.Verbosity,
		"verbosity",
		"v",
		0,
		"info log verbosity",
	)
	cmd.PersistentFlags().BoolVarP(
		&flags.Quiet,
		"quiet",
		"q",
		false,
		"silence all stderr output",
	)

	// add all top level subcommands
	// kind sub cmds

	// For sub-commands that kubeedge not modify
	// for example: kind get clusters/nodes can be repleaced by keink get clusters/nodes directly
	// modification：just import kind commands directly
	cmd.AddCommand(completion.NewCommand(logger, streams))
	// TODO: maybe we can add a delete kubeedge subcommand, just for stopping edgenodes and stop cloudcore/edgecore components
	cmd.AddCommand(delete.NewCommand(logger, streams))
	cmd.AddCommand(export.NewCommand(logger, streams))
	cmd.AddCommand(get.NewCommand(logger, streams))
	cmd.AddCommand(version.NewCommand(logger, streams))
	cmd.AddCommand(load.NewCommand(logger, streams))

	// for build and create commands, we should do some customization

	// keink create kubeedge command
	// for create command, we will not import kind related create command
	// just overwrite it with keink own create command
	// and we should use kubeedge image, or we cannot start cloudcore on master nodes
	createCmd := create.NewCommand(logger, streams)
	cmd.AddCommand(createCmd)

	// keink build edge-image command
	// for build command, we will not import kind related build command
	// just overwrite it with keink own build command
	buildCmd := build.NewCommand(logger, streams)
	cmd.AddCommand(buildCmd)

	return cmd
}

func runE(logger log.Logger, flags *flagpole, command *cobra.Command) error {
	// handle limited migration for --loglevel
	setLogLevel := command.Flag("loglevel").Changed
	setVerbosity := command.Flag("verbosity").Changed
	if setLogLevel && !setVerbosity {
		switch flags.LogLevel {
		case "debug":
			flags.Verbosity = 3
		case "trace":
			flags.Verbosity = 2147483647
		}
	}
	// normal logger setup
	if flags.Quiet {
		// NOTE: if we are coming from app.Run handling this flag is
		// redundant, however it doesn't hurt, and this may be called directly.
		maybeSetWriter(logger, ioutil.Discard)
	}
	maybeSetVerbosity(logger, log.Level(flags.Verbosity))
	// warn about deprecated flag if used
	if setLogLevel {
		if cmd.ColorEnabled(logger) {
			logger.Warn("\x1b[93mWARNING\x1b[0m: --loglevel is deprecated, please switch to -v and -q!")
		} else {
			logger.Warn("WARNING: --loglevel is deprecated, please switch to -v and -q!")
		}
	}
	return nil
}

// maybeSetWriter will call logger.SetWriter(w) if logger has a SetWriter method
func maybeSetWriter(logger log.Logger, w io.Writer) {
	type writerSetter interface {
		SetWriter(io.Writer)
	}
	v, ok := logger.(writerSetter)
	if ok {
		v.SetWriter(w)
	}
}

// maybeSetVerbosity will call logger.SetVerbosity(verbosity) if logger
// has a SetVerbosity method
func maybeSetVerbosity(logger log.Logger, verbosity log.Level) {
	type verboser interface {
		SetVerbosity(log.Level)
	}
	v, ok := logger.(verboser)
	if ok {
		v.SetVerbosity(verbosity)
	}
}
