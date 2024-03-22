package e

import "fmt"

func Wrap(logNumber int, msg string, err error) error {
	if err == nil {
		return fmt.Errorf("[WAR] [%d] %s", logNumber, msg)
	}
	return fmt.Errorf("[ERR] [%d] %s: %w", logNumber, msg, err)
}
