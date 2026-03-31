//revive:disable:package-comments
package helm

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"

	"github.com/katastroma/orpheus/internal/render"
)

const releaseName = "release"

// Render loads a Helm chart from the given files and returns the rendered manifest blob
func Render(files render.Files) ([]byte, error) {
	buffered := make([]*loader.BufferedFile, 0, len(files))
	for name, data := range files {
		buffered = append(buffered, &loader.BufferedFile{Name: name, Data: data})
	}

	chart, err := loader.LoadFiles(buffered)
	if err != nil {
		return nil, fmt.Errorf("loading chart: %w", err)
	}

	cfg := new(action.Configuration)
	install := action.NewInstall(cfg)
	install.DryRun = true
	install.ClientOnly = true
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

	return []byte(rel.Manifest), nil
}
