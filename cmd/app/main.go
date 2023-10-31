package app

import (
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"
	kindcmd "sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/exec"
	"sigs.k8s.io/kind/pkg/log"

	"github.com/kubeedge/keink/pkg/cmd"
)

// Main is the kind main(), it will invoke Run(), if an error is returned
// it will then call os.Exit
func Main() {
	if err := Run(kindcmd.NewLogger(), kindcmd.StandardIOStreams(), os.Args[1:]); err != nil {
		os.Exit(1)
	}
}

// Run invokes the kind root command, returning the error.
// See: sigs.k8s.io/kind/pkg/cmd/kind
func Run(logger log.Logger, streams kindcmd.IOStreams, args []string) error {
	// NOTE: we handle the quiet flag here so we can fully silence cobra
	if checkQuiet(args) {
		// If we are in quiet mode, we want to suppress all status output
		// Only streams.Out should be written to (program output)
		logger = log.NoopLogger{}
		streams.ErrOut = ioutil.Discard
	}
	// Actually run the command
	c := cmd.NewCommand(logger, streams)
	c.SetArgs(args)
	if err := c.Execute(); err != nil {
		logError(logger, err)
		return err
	}
	return nil
}

// checkQuiet returns true if -q / --quiet was set in args
func checkQuiet(args []string) bool {
	flags := pflag.NewFlagSet("persistent-quiet", pflag.ContinueOnError)
	flags.ParseErrorsWhitelist.UnknownFlags = true
	quiet := false
	flags.BoolVarP(
		&quiet,
		"quiet",
		"q",
		false,
		"silence all stderr output",
	)
	// NOTE: pflag will error if -h / --help is specified
	// We don't care here. That will be handled downstream
	// It will also call flags.Usage so we're making that no-op
	flags.Usage = func() {}
	_ = flags.Parse(args)
	return quiet
}

// logError logs the error and the root stack trace if there is one
func logError(logger log.Logger, err error) {
	colorEnabled := kindcmd.ColorEnabled(logger)
	if colorEnabled {
		logger.Errorf("\x1b[31mERROR\x1b[0m: %v", err)
	} else {
		logger.Errorf("ERROR: %v", err)
	}
	// Display output if the error was from running a command ...
	if execErr := exec.RunErrorForError(err); execErr != nil {
		if colorEnabled {
			logger.Errorf("\x1b[31mCommand Output\x1b[0m: %s", execErr.Output)
		} else {
			logger.Errorf("\nCommand Output: %s", execErr.Output)
		}
	}
	// TODO: Stack trace should probably be guarded by a higher level ...?
	if logger.V(1).Enabled() {
		// Then display the stack trace if any (there should be one...)
		if trace := errors.StackTrace(err); trace != nil {
			if colorEnabled {
				logger.Errorf("\x1b[31mStack Trace\x1b[0m: %+v", trace)
			} else {
				logger.Errorf("\nStack Trace: %+v", trace)
			}
		}
	}
}
