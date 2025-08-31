package packages

import "errors"

type VersionSpecifier string

func NewVersionSpecifier(versionSpecifierStr string) (VersionSpecifier, error) {

	if versionSpecifierStr == "" {
		return "", errors.New("version specifier must be not empty")
	}

	return VersionSpecifier(versionSpecifierStr), nil
}

func (v VersionSpecifier) String() string {
	return string(v)
}
