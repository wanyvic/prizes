package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/docker/rootless"
	"github.com/moby/buildkit/util/apicaps"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	prizesversion "github.com/wanyvic/prizes/version"
)

var (
	honorXDG bool
)

func newDaemonCommand() (*cobra.Command, error) {

	cmd := &cobra.Command{
		Use:           "prizes [OPTIONS]",
		Short:         "A mointor  for docker swarm manager.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			if cmd.HasSubCommands() {
				return errors.Errorf("\n" + strings.TrimRight(cmd.UsageString(), "\n"))
			}
			return errors.Errorf(
				"\"%s\" accepts no argument(s).\nSee '%s --help'.\n\nUsage:  %s\n\n%s",
				cmd.CommandPath(),
				cmd.CommandPath(),
				cmd.UseLine(),
				cmd.Short,
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		DisableFlagsInUseLine: true,
		Version:               fmt.Sprintf("%s, build %s", prizesversion.Version, prizesversion.GitCommit),
	}
	// cli.SetupRootCommand(cmd)

	flags := cmd.Flags()
	flags.BoolP("version", "v", false, "Print version information and quit")

	return cmd, nil
}

func init() {
	if prizesversion.ProductName != "" {
		apicaps.ExportedProduct = prizesversion.ProductName
	}
	// When running with RootlessKit, $XDG_RUNTIME_DIR, $XDG_DATA_HOME, and $XDG_CONFIG_HOME needs to be
	// honored as the default dirs, because we are unlikely to have permissions to access the system-wide
	// directories.
	//
	// Note that even running with --rootless, when not running with RootlessKit, honorXDG needs to be kept false,
	// because the system-wide directories in the current mount namespace are expected to be accessible.
	// ("rootful" dockerd in rootless dockerd, #38702)
	honorXDG = rootless.RunningWithRootlessKit()
}

func main() {
	if reexec.Init() {
		return
	}

	// initial log formatting; this setting is updated after the daemon configuration is loaded.
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: jsonmessage.RFC3339NanoFixed,
		FullTimestamp:   true,
	})

	// Set terminal emulation based on platform as required.
	_, stdout, stderr := term.StdStreams()

	// initLogging(stdout, stderr)

	onError := func(err error) {
		fmt.Fprintf(stderr, "%s\n", err)
		os.Exit(1)
	}

	cmd, err := newDaemonCommand()
	if err != nil {
		onError(err)
	}
	cmd.SetOutput(stdout)
	if err := cmd.Execute(); err != nil {
		onError(err)
	}
}
