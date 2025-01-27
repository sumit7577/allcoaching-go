package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// HandleUniqueConstraintError generates a user-friendly error message for UNIQUE constraint violations.
func HandleUniqueConstraintError(err error) error {
	if err == nil {
		return nil
	}
	re := regexp.MustCompile(`UNIQUE constraint failed: ([a-zA-Z_]+\.[a-zA-Z_]+)`)
	matches := re.FindStringSubmatch(err.Error())

	if len(matches) == 2 {
		// Extract the full column name (e.g., "user.phone" or "order.email")
		fullFieldName := matches[1]

		parts := strings.Split(fullFieldName, ".")
		if len(parts) == 2 {
			field := strings.Title(strings.ReplaceAll(parts[1], "_", " "))
			return errors.New(fmt.Sprintf("%s is already used", field))
		}
		return errors.New(fmt.Sprintf("Please Input valid %s data", parts[0]))
	}
	return err
}
