package packages

import (
	"errors"
	"strings"
)

type Version string

func NewVersion(versionStr string) (Version, error) {
	versionStr = strings.TrimSpace(versionStr)

	if versionStr == "" {
		return "", errors.New("package version must be not empty")
	}

	return Version(versionStr), nil
}

func (v Version) String() string {
	return string(v)
}
