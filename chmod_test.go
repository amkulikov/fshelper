package fshelper

import (
	"testing"
)

func TestChmod(t *testing.T) {
	if err := createSourceData(); err != nil {
		t.Fatalf("Can't create source data: %s", err)
	}
	if err := Chmod(ChmodOnlyFiles|ChmodRecursive|ChmodVerbose, 0666, srcBaseDir); err != nil {
		t.Fatalf("Can't copy: %s", err)
	}
}
