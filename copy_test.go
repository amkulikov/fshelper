package fshelper

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const ControlData = "TestData"

var srcBaseDir = filepath.Join(os.TempDir(), "fspathsrc")
var destBaseDir = filepath.Join(os.TempDir(), "fspathdest")

func createSourceData() error {
	if err := os.MkdirAll(srcBaseDir, 0755); err != nil {
		return err
	}
	err := ioutil.WriteFile(filepath.Join(srcBaseDir, "1.txt"), []byte(ControlData), 0644)
	if err != nil {
		return err
	}
	subDir := filepath.Join(srcBaseDir, "sub")

	if err = os.MkdirAll(subDir, 0755); err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(subDir, "2.txt"), []byte(ControlData), 0644)
	if err != nil {
		return err
	}

	return nil
}

func TestCopy(t *testing.T) {
	if err := createSourceData(); err != nil {
		t.Fatalf("Can't create source data: %s", err)
	}
	if err := Copy(CopyRecursive|CopyPreserveMode|CopyVerbose, srcBaseDir, destBaseDir); err != nil {
		t.Fatalf("Can't copy: %s", err)
	}

	if err := Copy(CopyPreserveMode|CopyVerbose, filepath.Join(srcBaseDir, "1.txt"), filepath.Join(srcBaseDir, "5.txt")); err != nil {
		t.Fatalf("Can't copy: %s", err)
	}
}
