package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/pidfile"
	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/docker/rootless"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wanyvic/prizes/cli"
	"github.com/wanyvic/prizes/cmd/prizesd/config"
	"github.com/wanyvic/prizes/prizesversion"
)

var (
	honorXDG bool
	PIDFile  = "prizesd.pid"
)

func newDaemonCommand() (*cobra.Command, error) {
	opts := newDaemonOptions(config.New())

	cmd := &cobra.Command{
		Use:           "prizesd [OPTIONS]",
		Short:         "A monitor for docker swarm",
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
			opts.flags = cmd.Flags()
			return runDaemon(opts)
		},
		DisableFlagsInUseLine: true,
		Version:               fmt.Sprintf("%s, build %s", prizesversion.Version, prizesversion.GitCommit),
	}
	cli.SetupRootCommand(cmd)

	flags := cmd.Flags()
	flags.BoolP("version", "v", false, "Print version information and quit")
	defaultDaemonConfigFile, err := getDefaultDaemonConfigFile()
	if err != nil {
		return nil, err
	}
	flags.StringVar(&opts.configFile, "config-file", defaultDaemonConfigFile, "Daemon configuration file")
	opts.InstallFlags(flags)
	if err := installConfigFlags(opts.daemonConfig, flags); err != nil {
		return nil, err
	}
	return cmd, nil
}

func init() {
	// if dockerversion.ProductName != "" {
	// 	apicaps.ExportedProduct = dockerversion.ProductName
	// }
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
		TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00",
		FullTimestamp:   true,
	})

	// Set terminal emulation based on platform as required.
	_, stdout, stderr := term.StdStreams()

	onError := func(err error) {
		fmt.Fprintf(stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := PIDFileCheck(); err != nil {
		onError(err)
	}
	cmd, err := newDaemonCommand()
	if err != nil {
		onError(err)
	}
	cmd.SetOutput(stdout)
	if err := cmd.Execute(); err != nil {
		onError(err)
	}
	logrus.Debug("prizesd exit")
}
func PIDFileCheck() error {
	path := filepath.Join(os.TempDir(), PIDFile)
	file, err := pidfile.New(path)
	if err != nil {
		return errors.New("prizesd has been started, please stop first")
	}
	logrus.Debug(file)
	return nil
}
