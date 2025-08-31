package packages_test

import (
	"testing"

	"github.com/0hJonny/python-deps-crawler/internal/fetcher/domain/packages"
)

func TestNewPackageName_NameValidation(t *testing.T) {
	t.Run("should return error for empty name", func(t *testing.T) {
		t.Parallel()
		packageNameStr := ""

		_, err := packages.NewPackageName(packageNameStr)

		if err == nil {
			t.Fatalf("package name must be not empty")
		}
	})

	t.Run("should replace spaces in the package name", func(t *testing.T) {
		t.Parallel()
		packageNameStr := "package name"

		packageName, err := packages.NewPackageName(packageNameStr)

		assertPackageNameError(t, err, packageNameStr)

		correctPackageName := "package_name"
		assertPackageNameEqualsError(t, packageName, correctPackageName)
	})

	t.Run("should create packageName with a valid name", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			PackageName        string
			CorrectPackageName string
		}{
			{
				PackageName:        "requests",
				CorrectPackageName: "requests",
			},
			{
				PackageName:        "  requests   ",
				CorrectPackageName: "requests",
			},
		}

		for _, tc := range testCases {
			packageName, err := packages.NewPackageName(tc.PackageName)

			assertPackageNameError(t, err, tc.PackageName)

			assertPackageNameEqualsError(t, packageName, tc.CorrectPackageName)
		}
	})
}

func assertPackageNameError(
	t *testing.T,
	err error,
	packageName string,
) {
	if err != nil {
		t.Fatalf("got error %v on a package name args: +%q", err, packageName)
	}
}

func assertPackageNameEqualsError(
	t *testing.T,
	packageName packages.PackageName,
	correctPackageName string,
) {
	if packageName.String() != correctPackageName {
		t.Fatalf("names are different: %q is not %q", packageName, correctPackageName)
	}
}
