package packages_test

import (
	"testing"

	"github.com/0hJonny/python-deps-crawler/internal/fetcher/domain/packages"
)

func TestNewVersionSpecifier_NameValidation(t *testing.T) {
	t.Run("should return error for empty name", func(t *testing.T) {
		versionSpecifierStr := ""

		_, err := packages.NewVersionSpecifier(versionSpecifierStr)

		if err == nil {
			t.Fatal("version specifier must be not empty")
		}
	})

}
