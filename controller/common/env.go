package common

import (
	"fmt"
	"os"
	"strconv"
)

const (
	// EnvGRPCPort is the environment variable for the controller's GRPC port
	EnvGRPCPort = "COFFEE_GRPC_PORT"
	// EnvLogLevel is the environment variable for the controller's log level
	EnvLogLevel = "COFFEE_LOG_LEVEL"
)

// Port is a network port
type Port uint

// GetPortFromEnv returns a network port from an environment variable
func GetPortFromEnv(v string) (Port, error) {
	x, ok := os.LookupEnv(v)
	if !ok {
		return 0, fmt.Errorf("environment variable not set: %s", v)
	}

	p, err := strconv.ParseUint(x, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s set to an invalid value: %s", v, x)
	}
	return Port(p), nil
}
