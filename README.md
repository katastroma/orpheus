# Orpheus

Katastroma's renderer. Implements the
[keleustēs](https://github.com/katastroma/keleustes) client API.

Given source content, orpheus renders Kubernetes manifests using helm,
kustomize, or raw YAML.

## Communication

The source handler streams a tar archive of source content via the keleustēs
`Render` RPC (client-streaming gRPC). The renderer type is passed as gRPC
metadata — the source handler detects it and the renderer dispatches to the
matching backend.

After receiving the full stream, orpheus responds to the source handler
immediately. Rendering and forwarding to the orderer happen asynchronously — the
source handler does not wait for rendering to complete. This decouples the
source handler's lifecycle from rendering time.

Orpheus opens a client-streaming connection to the orderer's diataxis `Order`
RPC and streams the rendered manifest blob. The orderer address is configured
via `ORDERER_ADDR`.

## Backend Dispatch

Orpheus registers rendering backends at startup: helm, kustomize, and plain
YAML. The renderer type from gRPC metadata selects which backend processes the
request.

Each backend implements `Receive` (buffer the tar stream) and `Render` (produce
a YAML manifest blob). Adding a new rendering backend means implementing that
interface and registering it.

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

## Assumptions

- Source content arrives as a tar archive with bare paths (no directory prefix).
  The helm backend handles repackaging into the format Helm's `LoadArchive`
  expects.
- The renderer type in gRPC metadata is correct. Orpheus trusts the source
  handler's detection.
- The orderer is reachable at `ORDERER_ADDR`. Forwarding failures are logged but
  not retried — pharos handles failure recovery via replay at the source handler.

## Design Direction

**Per-source-target renderer selection.** The renderer type is currently
detected by the source handler and applies to the entire request. In the future,
the renderer should be configurable at the source target (the ConfigMap), so
tenants can specify which renderer implementation processes their source.

**Override values from source target.** For helm rendering, override values
should be storable alongside the source target definition in the ConfigMap. The
helm backend would pull these values and apply them during rendering. This
enables platform-hosted overrides and allows platform-owned secrets to be
injected into rendered manifests without the tenant needing to manage them in
their source repository.
