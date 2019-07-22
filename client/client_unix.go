// +build linux freebsd openbsd darwin

package client // import "github.com/wanyvic/prizes/client"

// DefaultDockerHost defines os specific default if DOCKER_HOST is unset
const DefaultDockerHost = "unix:///var/run/prizes.sock"

const defaultProto = "unix"
const defaultAddr = "/var/run/prizes.sock"
