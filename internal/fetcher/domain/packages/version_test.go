package packages_test

import (
	"testing"

	"github.com/0hJonny/python-deps-crawler/internal/fetcher/domain/packages"
)

func TestNewVersion_NameValidation(t *testing.T) {
	t.Run("should return error for empty name", func(t *testing.T) {
		t.Parallel()
		version := ""

		_, err := packages.NewVersion(version)

		if err == nil {
			t.Fatal("package version must be not empty")
		}
	})

	t.Run("should create a valid version", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			Version        string
			CorrectVersion string
		}{
			{
				Version:        "2.32.5",
				CorrectVersion: "2.32.5",
			},
			{
				Version:        "  2.32.5   ",
				CorrectVersion: "2.32.5",
			},
		}

		for _, tc := range testCases {
			version, err := packages.NewVersion(tc.Version)

			if err != nil {
				t.Fatalf("got error %v on a version args: +%q", err, tc)
			}

			if version != packages.Version(tc.CorrectVersion) {
				t.Fatalf("versions are different: %q is not %q", version, tc.CorrectVersion)
			}
		}
	})
}
