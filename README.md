# Orpheus

Katastroma's renderer. Implements the
[keleustēs](https://github.com/katastroma/keleustes) interface.

Given source content, orpheus renders Kubernetes manifests using helm,
kustomize, or raw YAML.

## Rendering, Not Ordering

Orpheus renders source content into Kubernetes manifests. The output manifest
order depends on the render backend. Note that none of them guarantee a
consistent apply-safe order and will largely vary depending on the nature of the
source that is rendered:

- **Helm** — `action.Install` (the SDK equivalent of `helm template`) outputs
  manifests sorted by Helm's hardcoded `InstallOrder` (namespaces before RBAC
  before workloads, etc). This sort happens in `renderResources()` via
  `releaseutil.SortManifests()` before the dry-run return, so `Release.Manifest`
  comes back pre-sorted. CRDs go first when `IncludeCRDs` is set.
- **Kustomize** — `krusty.MakeDefaultOptions()` sets `Reorder` to `None`. The
  SDK returns manifests in accumulation order by default, unlike the
  `kustomize build` CLI which defaults to legacy sort order. To get sorted
  output from the SDK you must explicitly set `Reorder: ReorderOptionLegacy`.
- **Raw YAML** — no ordering whatsoever. Documents come back in file order.

Because each backend has different ordering behavior (or none), and because the
Kubernetes API applies one resource at a time, manifest ordering is a separate
concern handled by a dedicated ordering service before manifests reach the
provisioner. Orpheus's job ends at producing the correct set of manifests — it
makes no guarantees about their order.
