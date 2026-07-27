package flags

import (
	"fmt"
	"strconv"
	"strings"
)

type HostPort string

func (h *HostPort) Type() string {
	return "host:port"
}

func (h *HostPort) String() string {
	return string(*h)
}

func (h *HostPort) Set(value string) error {
	// Format of argument is "host:port"
	if err := ParseHostPort(value); err != nil {
		return err
	}
	*h = HostPort(value)
	return nil
}

func ParseHostPort(value string) error {
	sep := strings.Split(value, ":")
	if len(sep) != 2 {
		return fmt.Errorf("host:port format expected")
	}
	_, err := strconv.Atoi(sep[1])
	if err != nil {
		return fmt.Errorf("port must be a number")
	}
	return nil
}
