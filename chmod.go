package fshelper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FlagsChmod uint

const (
	ChmodChanges   FlagsChmod = 1 << iota
	ChmodVerbose
	ChmodRecursive
	ChmodOnlyDirs
	ChmodOnlyFiles
)

var (
	ErrSkipFile          = errors.New("skip file")
	ErrSkipDir           = errors.New("skip dir")
	ErrIncompatibleFlags = errors.New("incompatible combination of flags")
)

func Chmod(flags FlagsChmod, mode os.FileMode, path string) error {
	if ChmodOnlyFiles&flags != 0 && ChmodOnlyDirs&flags != 0 {
		return ErrIncompatibleFlags
	}

	if ChmodVerbose&flags != 0 && ChmodChanges&flags != 0 {
		flags &= ^ChmodChanges
	}

	if ChmodRecursive&flags != 0 {
		return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			chmod(flags, mode, info, path)
			// TODO: may be errors count (except skips)
			return nil
		})
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	return chmod(flags, mode, info, path)
}

func chmod(flags FlagsChmod, mode os.FileMode, info os.FileInfo, path string) (err error) {
	defer func() {
		if err != nil && (ChmodVerbose|ChmodChanges)&flags == 0 {
			fmt.Printf("%s: %s\n", path, err)
		}
	}()

	if info.IsDir() {
		if ChmodOnlyFiles&flags != 0 {
			return ErrSkipDir
		}
	} else {
		if ChmodOnlyDirs&flags != 0 {
			return ErrSkipFile
		}
	}

	currentMode := info.Mode()
	if currentMode != mode {
		if ChmodChanges&flags != 0 || ChmodVerbose&flags != 0 {
			fmt.Printf("%s: %v -> %v\n", path, currentMode, mode)
		}
	} else if ChmodVerbose&flags != 0 {
		fmt.Printf("%s: already %v\n", path, currentMode)
	}

	return os.Chmod(path, mode)

}
