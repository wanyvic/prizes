package cli

import (
	"fmt"

	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/rootless"
	"github.com/moby/buildkit/util/apicaps"
	"github.com/spf13/cobra"
)

var (
	honorXDG bool
)

func newDaemonCommand() (*cobra.Command, error) {
	// opts := newDaemonOptions(config.New())

	cmd := &cobra.Command{
		Use:           "dockerd [OPTIONS]",
		Short:         "A self-sufficient runtime for containers.",
		SilenceUsage:  true,
		SilenceErrors: true,
		// Args:          cli.NoArgs,
		// RunE: func(cmd *cobra.Command, args []string) error {
		// 	// opts.flags = cmd.Flags()
		// 	// return runDaemon(opts)
		// },
		DisableFlagsInUseLine: true,
		Version:               fmt.Sprintf("%s, build %s", dockerversion.Version, dockerversion.GitCommit),
	}
	SetupRootCommand(cmd)

	flags := cmd.Flags()
	flags.BoolP("version", "v", false, "Print version information and quit")
	// defaultDaemonConfigFile := "/home/wany/.prizes.yaml" //getDefaultDaemonConfigFile()
	// if err != nil {
	// 	return nil, err
	// }
	// flags.StringVar(&opts.configFile, "config-file", defaultDaemonConfigFile, "Daemon configuration file")
	// opts.InstallFlags(flags)
	// if err := installConfigFlags(opts.daemonConfig, flags); err != nil {
	// 	return nil, err
	// }
	// installServiceFlags(flags)

	return cmd, nil
}

func init() {
	if dockerversion.ProductName != "" {
		apicaps.ExportedProduct = dockerversion.ProductName
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
