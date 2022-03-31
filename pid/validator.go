package pid

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/sys/unix"
)

func ValidatePIDFile(pidfile string) error {
	if _, err := os.Stat(pidfile); err == nil {
		err = unix.Access(pidfile, unix.W_OK|unix.R_OK)
		if err != nil {
			return fmt.Errorf("existing PID file is not readable/writable: %w", err)
		}
		return nil
	}

	var d []byte

	defer os.Remove(pidfile)

	if err := ioutil.WriteFile(pidfile, d, 0644); err != nil {
		return fmt.Errorf("PID file is not readable/writable at path: %w", err)
	}

	return nil
}
