package main

import (
	"time"

	"github.com/docker/docker/opts"
	"github.com/spf13/pflag"
	"github.com/wanyvic/prizes/cmd/prizesd/config"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
)

type MassgridRPCOptions struct {
	MassGridHost []string
	Username     string
	Password     string
}

type daemonOptions struct {
	configFile    string
	daemonConfig  *config.Config
	flags         *pflag.FlagSet
	Debug         bool
	Hosts         []string
	LogLevel      string
	TimeScale     int
	TimeStatement int
	TestNet       bool
	RPC           MassgridRPCOptions
	DataBaseHost  []string
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
	flags.IntVarP(&o.TimeStatement, "time-Scale-Statement", "", 5, "Set time cycle for statement Minute")
	flags.BoolVarP(&o.TestNet, "testnet", "", false, "Set massgrid testnet")

	hostOpt := opts.NewNamedListOptsRef("hosts", &o.Hosts, opts.ValidateHost)
	flags.VarP(hostOpt, "host", "H", "Daemon socket(s) to connect to")

	massgridHost := opts.NewNamedListOptsRef("massgridhosts", &o.RPC.MassGridHost, opts.ValidateHost)
	flags.VarP(massgridHost, "rpc-server", "", "MassGrid rpc host (default tcp://localhost:9442)")
	flags.StringVarP(&o.RPC.Username, "rpc-username", "", "user", "Set MassGrid rpc username")
	flags.StringVarP(&o.RPC.Password, "rpc-password", "", "password", "Set MassGrid rpc password")

	databaseHost := opts.NewNamedListOptsRef("databasehosts", &o.DataBaseHost, nil)
	flags.VarP(databaseHost, "db-server", "", "database host (default mongodb://localhost:27017)")
}

// SetDefaultOptions sets default values for options after flag parsing is
// complete
func (o *daemonOptions) SetDefaultOptions(flags *pflag.FlagSet) {

	refresh.TimeScale = time.Duration(refresh.DefaultTimeScale) * time.Millisecond
	o.DataBaseHost = append(o.DataBaseHost, "mongodb://localhost:27017")
	// Regardless of whether the user sets it to true or false, if they
	// specify --tlsverify at all then we need to turn on TLS
	// TLSVerify can be true even if not set due to DOCKER_TLS_VERIFY env var, so we need
	// to check that here as well

}
