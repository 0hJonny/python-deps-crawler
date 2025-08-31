package packages

import (
	"errors"
	"strings"
)

type PackageName string

func NewPackageName(packageNameStr string) (PackageName, error) {
	packageNameStr = strings.TrimSpace(packageNameStr)

	if packageNameStr == "" {
		return "", errors.New("package name must be not empty")
	}

	packageNameStr = strings.ReplaceAll(packageNameStr, " ", "_")

	return PackageName(packageNameStr), nil
}

func (p PackageName) String() string {
	return string(p)
}
