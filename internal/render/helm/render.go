//revive:disable:package-comments
package helm

import (
	"bytes"
	"fmt"
	"io"

	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart/loader"
	"helm.sh/helm/v4/pkg/cli"
)

const releaseName = "release"

// Backend implements render.Backend for Helm charts.
//
// The incoming tar contains bare paths (Chart.yaml, templates/foo.yaml).
// Helm's LoadArchive — the only stable in-memory chart loading API — requires
// a gzipped tar where every entry is nested under a directory prefix, because
// it was designed for pre-packaged chart archives (helm package output). It
// strips the first path segment from every entry during loading.
//
// To avoid leaking this helm packaging convention into upstream services, we
// repackage the tar here: entries are read one at a time, written into a new
// tar with a directory prefix, gzip-compressed, and buffered. This copies
// every byte twice — once through the repackage, once through LoadArchive's
// internal extraction — but keeps the helm SDK requirement contained in the
// only place that imports it.
type Backend struct {
	archive bytes.Buffer
}

// New creates a Helm rendering backend.
func New() *Backend {
	return &Backend{}
}

// Receive repackages the tar stream into the gzipped format LoadArchive expects.
func (b *Backend) Receive(r io.Reader) error {
	if _, err := io.Copy(&b.archive, repackage(r)); err != nil {
		return fmt.Errorf("repackaging chart archive: %w", err)
	}

	return nil
}

// Render loads the buffered chart archive and produces a rendered manifest blob.
func (b *Backend) Render() ([]byte, error) {
	chart, err := loader.LoadArchive(&b.archive)
	if err != nil {
		return nil, fmt.Errorf("loading chart: %w", err)
	}

	cfg := new(action.Configuration)
	install := action.NewInstall(cfg)
	install.DryRunStrategy = action.DryRunClient
	install.ReleaseName = releaseName
	// TODO: releaseName should come from the source target ConfigMap data
	// The tenant will configure this via the front-end. For now, hardcoded
	install.Namespace = cli.New().Namespace()

	// TODO: values should come from the source target ConfigMap data, passed
	// through from the source handler as gRPC metadata or similar.
	var vals map[string]any
	rel, err := install.Run(chart, vals)
	if err != nil {
		return nil, fmt.Errorf("rendering chart: %w", err)
	}

	return manifest(rel)
}
