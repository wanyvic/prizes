package main

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/wanyvic/prizes/cmd/prizesd/config"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
)

type daemonOptions struct {
	configFile   string
	daemonConfig *config.Config
	flags        *pflag.FlagSet
	Debug        bool
	Hosts        []string
	LogLevel     string
	TimeScale    int
}

// newDaemonOptions returns a new daemonFlags
func newDaemonOptions(config *config.Config) *daemonOptions {
	return &daemonOptions{
		daemonConfig: config,
	}
}

// InstallFlags adds flags for the common options on the FlagSet
func (o *daemonOptions) InstallFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&o.Debug, "debug", "D", false, "Enable debug mode")
	flags.StringVarP(&o.LogLevel, "log-level", "l", "info", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)
	flags.IntVarP(&o.TimeScale, "time-Scale", "t", refresh.DefaultTimeScale, "Set record Millisecond time scale to database")
}

// SetDefaultOptions sets default values for options after flag parsing is
// complete
func (o *daemonOptions) SetDefaultOptions(flags *pflag.FlagSet) {

	refresh.TimeScale = time.Duration(refresh.DefaultTimeScale) * time.Millisecond
	// Regardless of whether the user sets it to true or false, if they
	// specify --tlsverify at all then we need to turn on TLS
	// TLSVerify can be true even if not set due to DOCKER_TLS_VERIFY env var, so we need
	// to check that here as well

}
