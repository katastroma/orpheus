package helm

import (
	"testing"
)

func TestManifest_InvalidReleaseType(t *testing.T) {
	if _, err := manifest("not a release"); err == nil {
		t.Fatal("expected error for invalid release type")
	}
}
