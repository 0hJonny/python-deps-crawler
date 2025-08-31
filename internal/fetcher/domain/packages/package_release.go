package packages

import (
	"errors"
	"time"
)

type PackageRelease struct {
	id          PackageReleaseID
	packageName PackageName
	version     Version
	dependency  []Dependency
	uploadTime  time.Time
	url         string
}

func NewPackageRelease(
	id PackageReleaseID,
	nameStr string,
	versionStr string,
	uploadTime time.Time,
	url string,
) (*PackageRelease, error) {
	if id.IsEmpty() {
		return nil, errors.New("package release ID is required")
	}

	name, err := NewPackageName(nameStr)

	if err != nil {
		return nil, err
	}

	version, err := NewVersion(versionStr)

	if err != nil {
		return nil, err
	}

	return &PackageRelease{
		id:          id,
		packageName: name,
		version:     version,
		dependency:  make([]Dependency, 0),
		uploadTime:  uploadTime,
		url:         url,
	}, nil
}

func (p *PackageRelease) GetPackageName() PackageName {
	return p.packageName
}

func (p *PackageRelease) GetVersion() Version {
	return p.version
}

func (p *PackageRelease) GetID() PackageReleaseID {
	return p.id
}
