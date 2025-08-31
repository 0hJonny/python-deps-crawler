package packages

type PackageReleaseID string

func (id PackageReleaseID) String() string {
	return string(id)
}

func (id PackageReleaseID) IsEmpty() bool {
	return string(id) == ""
}
