package packages_test

import (
	"testing"

	"github.com/0hJonny/python-deps-crawler/internal/fetcher/domain/packages"
)

const DP_ERROR = "dependency error: %v"

func TestNewDependency_NameValidation(t *testing.T) {
	t.Run("should return error for empty package name", func(t *testing.T) {
		t.Parallel()
		var (
			packageNameStr      = ""
			versionSpecifierStr = ">=1.2.0, <2.0.0"
		)

		_, err := packages.NewDependency(packageNameStr, versionSpecifierStr)

		assertDependencyNilError(t, err)

	})

	t.Run("should return error for empty version specifier", func(t *testing.T) {
		t.Parallel()
		var (
			packageNameStr      = "requests"
			versionSpecifierStr = ""
		)

		_, err := packages.NewDependency(packageNameStr, versionSpecifierStr)

		assertDependencyNilError(t, err)
	})

	t.Run("should create dependency with a valid fields", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			PackageName      string
			VersionSpecifier string
		}{
			{
				PackageName:      "requests",
				VersionSpecifier: ">=1.2.0, <2.0.0",
			},
		}
		for _, tc := range testCases {
			_, err := packages.NewDependency(tc.PackageName, tc.VersionSpecifier)

			if err != nil {
				t.Fatalf(DP_ERROR, err)
			}
		}
	})

}

func assertDependencyNilError(t *testing.T, err error) {
	if err == nil {
		t.Fatalf(DP_ERROR, err)
	}
}
