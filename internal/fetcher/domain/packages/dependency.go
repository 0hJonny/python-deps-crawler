package packages

type Dependency struct {
	packageName      PackageName
	versionSpecifier VersionSpecifier
}

func NewDependency(packageNameStr, versionSpecifierStr string) (*Dependency, error) {
	packageName, err := NewPackageName(packageNameStr)

	if err != nil {
		return nil, err
	}

	versionSpecifier, err := NewVersionSpecifier(versionSpecifierStr)

	if err != nil {
		return nil, err
	}

	return &Dependency{
		packageName:      packageName,
		versionSpecifier: versionSpecifier,
	}, nil
}
