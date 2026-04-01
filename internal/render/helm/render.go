//revive:disable:package-comments
package helm

import (
	"fmt"
	"io"

	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart/loader"
	"helm.sh/helm/v4/pkg/cli"
)

const releaseName = "release"

// Render loads a Helm chart from the tar stream and returns the rendered
// manifest blob.
//
// The incoming tar contains bare paths (Chart.yaml, templates/foo.yaml).
// Helm's LoadArchive — the only stable in-memory chart loading API — requires
// a gzipped tar where every entry is nested under a directory prefix, because
// it was designed for pre-packaged chart archives (helm package output). It
// strips the first path segment from every entry during loading.
//
// To avoid leaking this helm packaging convention into upstream services, we
// repackage the tar here: entries are read one at a time, written into a new
// tar with a directory prefix, gzip-compressed, and piped to LoadArchive.
// This copies every byte twice — once through the repackage, once through
// LoadArchive's internal extraction — but keeps the helm SDK requirement
// contained in the only place that imports it.
func Render(r io.Reader) ([]byte, error) {
	archive := repackage(r)

	chart, err := loader.LoadArchive(archive)
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
