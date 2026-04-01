//revive:disable:package-comments
package helm

import (
	"fmt"

	"helm.sh/helm/v4/pkg/release"
)

func manifest(rel release.Releaser) ([]byte, error) {
	accessor, err := release.NewAccessor(rel)
	if err != nil {
		return nil, fmt.Errorf("accessing release: %w", err)
	}

	return []byte(accessor.Manifest()), nil
}
