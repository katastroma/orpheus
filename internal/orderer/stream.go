//revive:disable:package-comments
package orderer

import "context"

// StreamFunc streams rendered manifests to the orderer.
type StreamFunc func(ctx context.Context, manifests [][]byte) error
