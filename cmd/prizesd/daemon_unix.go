// +build !windows

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/cli/debug"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/wanyvic/gobtclib/client"
	"github.com/wanyvic/prizes/cmd/db"
	"github.com/wanyvic/prizes/cmd/db/mongodb"
	"github.com/wanyvic/prizes/cmd/prizesd/config"
	dockerapi "github.com/wanyvic/prizes/cmd/prizesd/docker"
	"github.com/wanyvic/prizes/cmd/prizesd/massgrid"
	"github.com/wanyvic/prizes/cmd/prizesd/prizeservice"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh/calculagraph"
	httpserver "github.com/wanyvic/prizes/cmd/prizesd/server"
)

type DaemonCli struct {
	*config.Config
	configFile *string
	flags      *pflag.FlagSet
}

// NewDaemonCli returns a daemon CLI
func NewDaemonCli() *DaemonCli {
	return &DaemonCli{}
}

func getDefaultDaemonConfigDir() (string, error) {
	if !honorXDG {
		return "/etc/prizes", nil
	}
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home := os.Getenv("HOME")
		if home == "" {
			return "", errors.New("could not get either XDG_CONFIG_HOME or HOME")
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "prizes"), nil
}

func getDefaultDaemonConfigFile() (string, error) {
	dir, err := getDefaultDaemonConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "daemon.json"), nil
}

func getDaemonConfDir(_ string) (string, error) {
	return getDefaultDaemonConfigDir()
}

// installConfigFlags adds flags to the pflag.FlagSet to configure the daemon
func installConfigFlags(conf *config.Config, flags *pflag.FlagSet) error {
	return nil
}

func runDaemon(opts *daemonOptions) error {
	daemonCli := NewDaemonCli()
	return daemonCli.start(opts)
}

func initLogging(_, stderr io.Writer) {
	logrus.SetOutput(stderr)
}
func loadDaemonCliConfig(opts *daemonOptions) (*config.Config, error) {
	conf := opts.daemonConfig
	conf.LogLevel = opts.LogLevel
	refresh.TimeScale = time.Duration(opts.TimeScale) * time.Millisecond
	logrus.Info("data refresh time duration ", refresh.TimeScale)
	prizeservice.StatementDuration = time.Duration(opts.TimeStatement) * time.Minute
	logrus.Info("statement time cycle ", prizeservice.StatementDuration)
	if opts.TestNet {
		massgrid.DefaultNetParams = massgrid.DefaultTestNetParams
		logrus.Info("set massgrid testnet")
	}
	if len(opts.RPC.MassGridHost) > 0 {
		massgrid.DefaultClientConfig = &client.Config{
			Host: opts.RPC.MassGridHost[0],
			User: opts.RPC.Username,
			Pass: opts.RPC.Password,
		}
		logrus.Info(fmt.Sprintf("set massgrid rpc host %s username %s password %s", opts.RPC.MassGridHost, opts.RPC.Username, opts.RPC.Password))
	}
	if len(opts.DataBaseHost) > 0 {
		db.DBimplement = &mongodb.MongDBClient{
			URI:      opts.DataBaseHost[0],
			DataBase: db.DBDefaultDataBase,
		}
		logrus.Info(fmt.Sprintf("set database host %s", opts.DataBaseHost[0]))
	}
	return conf, nil
}
func (cli *DaemonCli) start(opts *daemonOptions) (err error) {
	stopc := make(chan bool)
	defer close(stopc)

	opts.SetDefaultOptions(opts.flags)

	if cli.Config, err = loadDaemonCliConfig(opts); err != nil {
		return err
	}

	if err := configureDaemonLogs(cli.Config); err != nil {
		return err
	}

	cli.configFile = &opts.configFile
	cli.flags = opts.flags

	if cli.Config.Debug {
		debug.Enable()
	}

	logrus.Info("Starting up")
	if err := dockerapi.NewProxy(); err != nil {
		logrus.Warning(err)
	}
	if err := configureServer(opts.Hosts); err != nil {
		logrus.Warning(err)
	}
	loop := refresh.NewRefreshMoudle()
	loop.Start()
	if err := calculagraph.InitCalculagraph(); err != nil {
		logrus.Warning(err)
	}
	if err := CheckCalculagraph(); err != nil {
		logrus.Warning(err)
	}
	return nil
}

// configureDaemonLogs sets the logrus logging level and formatting
func configureDaemonLogs(conf *config.Config) error {
	if conf.LogLevel != "" {
		lvl, err := logrus.ParseLevel(conf.LogLevel)
		if err != nil {
			return fmt.Errorf("unable to parse logging level: %s", conf.LogLevel)
		}
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: jsonmessage.RFC3339NanoFixed,
		DisableColors:   conf.RawLogs,
		FullTimestamp:   true,
	})
	return nil
}
func configureServer(hosts []string) error {
	if len(hosts) > 0 {
		for i := 0; i < len(hosts); i++ {
			protoAddr := hosts[i]
			logrus.Info("configure Server: ", protoAddr)
			protoAddrParts := strings.SplitN(protoAddr, "://", 2)
			if len(protoAddrParts) != 2 {
				return fmt.Errorf("bad format %s, expected PROTO://ADDR", protoAddr)
			}
			proto := protoAddrParts[0]
			addr := protoAddrParts[1]

			server, err := httpserver.NewServerWithOpts(httpserver.ServerOpts{Proto: proto, Addr: addr})
			if err != nil {
				return err
			}
			go server.Start()
		}
	} else {
		logrus.Info("starting default unix sock server")
		server, err := httpserver.NewServer("unix")
		if err != nil {
			return err
		}
		go server.Start()
	}
	return nil
}
