package migration

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-gormigrate/gormigrate/v2"
)

func AllMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		// ... other migrations
	}
}

var (
	// ErrInvalidIdFormat is invalid id format.
	ErrInvalidIdFormat = fmt.Errorf("invalid id format")
	// ErrConvertIdToInt is convert id to int error.
	ErrConvertIdToInt = fmt.Errorf("convert id to int error")
)

func IDToVersion(id string) (int, error) {
	// convert migration ID to version
	// e.g. "0001" -> 1
	re := regexp.MustCompile(`^(\d+)$`)
	matches := re.FindStringSubmatch(id)
	if len(matches) != 2 {
		return 0, ErrInvalidIdFormat
	}
	version_string := matches[1]
	version, err := strconv.Atoi(strings.TrimLeft(version_string, "0"))
	if err != nil {
		return 0, ErrConvertIdToInt
	}

	return version, nil
}
