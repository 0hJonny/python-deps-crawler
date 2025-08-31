package packages_test

import (
	"testing"
	"time"

	"github.com/0hJonny/python-deps-crawler/internal/fetcher/domain/packages"
)

func TestNewPackage_NameValidation(t *testing.T) {
	t.Run("should return error for empty id", func(t *testing.T) {
		t.Parallel()
		var (
			id         = packages.PackageReleaseID("")
			name       = "requests"
			version    = "2.32.5"
			uploadTime = time.Now()
			url        = "example.com"
		)
		_, err := packages.NewPackageRelease(id, name, version, uploadTime, url)

		if err == nil {
			t.Fatal("package id must be not empty")
		}
	})
	t.Run("should return error for empty name", func(t *testing.T) {
		t.Parallel()
		var (
			id         = packages.PackageReleaseID("new")
			name       = ""
			version    = "2.32.5"
			uploadTime = time.Now()
			url        = "example.com"
		)
		_, err := packages.NewPackageRelease(id, name, version, uploadTime, url)

		if err == nil {
			t.Fatal("package name must be not empty")
		}
	})

	t.Run("should return error for empty version", func(t *testing.T) {
		t.Parallel()
		var (
			id         = packages.PackageReleaseID("new")
			name       = "requests"
			version    = ""
			uploadTime = time.Now()
			url        = "example.com"
		)

		_, err := packages.NewPackageRelease(id, name, version, uploadTime, url)

		if err == nil {
			t.Fatal("package version must be not empty")
		}
	})

	t.Run("should create package with a valid name", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			ID             packages.PackageReleaseID
			Name           string
			CorrectName    string
			Version        string
			CorrectVersion string
			UploadTime     time.Time
			Url            string
		}{
			{
				ID:             "New",
				Name:           "requests",
				CorrectName:    "requests",
				Version:        "2.32.5",
				CorrectVersion: "2.32.5",
				UploadTime:     time.Now(),
				Url:            "example.com",
			},
			{
				ID:             "New",
				Name:           "    requests",
				CorrectName:    "requests",
				Version:        "2.32.5   ",
				CorrectVersion: "2.32.5",
				UploadTime:     time.Now(),
				Url:            "example.com",
			},
			{

				ID:             "New",
				Name:           "    my lib",
				CorrectName:    "my_lib",
				Version:        "2.32.5   ",
				CorrectVersion: "2.32.5",
				UploadTime:     time.Now(),
				Url:            "example.com",
			},
		}

		for _, tc := range testCases {
			pkg, err := packages.NewPackageRelease(tc.ID, tc.Name, tc.Version, tc.UploadTime, tc.Url)

			if err != nil {
				t.Fatalf("got error %v on a package args: +%q", err, tc)
			}

			if packageName := pkg.GetPackageName(); packageName.String() != tc.CorrectName {
				t.Fatalf("names are different: %q is not %q", packageName, tc.CorrectName)
			}

			if version := pkg.GetVersion(); version.String() != tc.CorrectVersion {
				t.Fatalf("versions are different: %q is not %q", version, tc.CorrectVersion)
			}
		}
	})
}

func TestGetID_ValidateID(t *testing.T) {
	t.Run("should return valid id", func(t *testing.T) {
		t.Parallel()
		var (
			id         = packages.PackageReleaseID("requests-2.32.5")
			name       = "requests"
			version    = "2.32.5"
			uploadTime = time.Now()
			url        = "example.com"
		)

		pkg, _ := packages.NewPackageRelease(id, name, version, uploadTime, url)
		expected := "requests-2.32.5"
		if id := pkg.GetID().String(); id != expected {
			t.Fatalf("package id are different: %q is not %q", id, expected)
		}
	})
}
