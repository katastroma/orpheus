package plain

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
)

func TestMatch(t *testing.T) {
	if !Match(memfs.New()) {
		t.Fatal("expected Match to return true")
	}
}
