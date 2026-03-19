# Orpheus

Katastroma's reference resolver. Implements the
[keleustēs](https://github.com/katastroma/keleustes) interface.

Given a source repository and credentials, orpheus clones the repo, renders the
manifests, and produces the resource inventory — the set of Kubernetes resources
that should exist according to that source.

## What This Is

A Go module that implements the keleustēs resolver interface. Consumed as a
dependency by orchestrators like
[pedalion](https://github.com/katastroma/pedalion) that need to resolve git
sources into resource inventories.

## Ecosystem

Orpheus is one of two reference implementations provided by
[katastroma](https://github.com/katastroma):

- **Orpheus** (this) — reference resolver, implements
  [keleustēs](https://github.com/katastroma/keleustes)
- **[Histia](https://github.com/katastroma/histia)** — reference provisioner,
  implements [katartismos](https://github.com/katastroma/katartismos)
