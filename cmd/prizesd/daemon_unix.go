// +build !windows

package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/cli/debug"
	"github.com/docker/docker/daemon/listeners"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/wanyvic/prizes/cmd/prizesd/config"
	"github.com/wanyvic/prizes/cmd/prizesd/refresh"
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

	logrus.Info("Starting up")

	cli.configFile = &opts.configFile
	cli.flags = opts.flags

	if cli.Config.Debug {
		debug.Enable()
	}

	refresh.WhileLoop()
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
func newAPIServerConfig(cli *DaemonCli) (*apiserver.Config, error) {
	serverConfig := &apiserver.Config{
		Logging:     true,
		SocketGroup: cli.Config.SocketGroup,
		Version:     dockerversion.Version,
		CorsHeaders: cli.Config.CorsHeaders,
	}

	if cli.Config.TLS {
		tlsOptions := tlsconfig.Options{
			CAFile:             cli.Config.CommonTLSOptions.CAFile,
			CertFile:           cli.Config.CommonTLSOptions.CertFile,
			KeyFile:            cli.Config.CommonTLSOptions.KeyFile,
			ExclusiveRootPools: true,
		}

		if cli.Config.TLSVerify {
			// server requires and verifies client's certificate
			tlsOptions.ClientAuth = tls.RequireAndVerifyClientCert
		}
		tlsConfig, err := tlsconfig.Server(tlsOptions)
		if err != nil {
			return nil, err
		}
		serverConfig.TLSConfig = tlsConfig
	}

	if len(cli.Config.Hosts) == 0 {
		cli.Config.Hosts = make([]string, 1)
	}

	return serverConfig, nil
}

func loadListeners(cli *DaemonCli, serverConfig *apiserver.Config) ([]string, error) {
	var hosts []string
	seen := make(map[string]struct{}, len(cli.Config.Hosts))

	for i := 0; i < len(cli.Config.Hosts); i++ {
		var err error
		if cli.Config.Hosts[i], err = dopts.ParseHost(cli.Config.TLS, honorXDG, cli.Config.Hosts[i]); err != nil {
			return nil, errors.Wrapf(err, "error parsing -H %s", cli.Config.Hosts[i])
		}
		if _, ok := seen[cli.Config.Hosts[i]]; ok {
			continue
		}
		seen[cli.Config.Hosts[i]] = struct{}{}

		protoAddr := cli.Config.Hosts[i]
		protoAddrParts := strings.SplitN(protoAddr, "://", 2)
		if len(protoAddrParts) != 2 {
			return nil, fmt.Errorf("bad format %s, expected PROTO://ADDR", protoAddr)
		}

		proto := protoAddrParts[0]
		addr := protoAddrParts[1]

		// It's a bad idea to bind to TCP without tlsverify.
		if proto == "tcp" && (serverConfig.TLSConfig == nil || serverConfig.TLSConfig.ClientAuth != tls.RequireAndVerifyClientCert) {
			logrus.Warn("[!] DON'T BIND ON ANY IP ADDRESS WITHOUT setting --tlsverify IF YOU DON'T KNOW WHAT YOU'RE DOING [!]")
		}
		ls, err := listeners.Init(proto, addr, serverConfig.SocketGroup, serverConfig.TLSConfig)
		if err != nil {
			return nil, err
		}
		// If we're binding to a TCP port, make sure that a container doesn't try to use it.
		if proto == "tcp" {
			if err := allocateDaemonPort(addr); err != nil {
				return nil, err
			}
		}
		logrus.Debugf("Listener created for HTTP on %s (%s)", proto, addr)
		hosts = append(hosts, protoAddrParts[1])
		cli.api.Accept(addr, ls...)
	}

	return hosts, nil
}
