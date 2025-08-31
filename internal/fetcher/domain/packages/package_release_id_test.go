package packages

import (
	"testing"
)

func TestNewPackageReleaseID_IDValidation(t *testing.T) {
	t.Run("empty id validation ", func(t *testing.T) {

		id := PackageReleaseID("")

		if id.IsEmpty() == false {
			t.Fatalf("package id is empty")
		}
	})
}
