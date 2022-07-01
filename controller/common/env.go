package common

import (
	"fmt"
	"os"
	"strconv"
)

const (
	ENV_GRPC_PORT = "COFFEE_GRPC_PORT"
	ENV_LOG_LEVEL = "COFFEE_LOG_LEVEL"
)

func GetPortFromEnv(v string) (uint, error) {
	x, ok := os.LookupEnv(v)
	if !ok {
		return 0, fmt.Errorf("Environment variable not set: %s", v)
	}

	p, err := strconv.ParseUint(x, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s set to an invalid value: %s", v, x)
	}
	return uint(p), nil
}
