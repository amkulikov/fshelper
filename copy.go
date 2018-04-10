package fshelper

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/unix"
)

type FlagsCopy uint

const (
	CopyRecursive     FlagsCopy = 1 << iota
	CopyPreserveMode
	CopyPreserveOwner
	CopyParents
	CopyVerbose
)

const (
	DefaultDirMode  = 0755
	DefaultFileMode = 0644
)

func Copy(src, dest string, flags FlagsCopy) error {
	srcPrefix := string(os.PathSeparator)
	if CopyParents&flags != 0 {
		fullSrc, err := filepath.Abs(src)
		if err != nil {
			return err
		}
		src = fullSrc
	} else {
		srcPrefix = filepath.Dir(src)
	}
	if CopyRecursive&flags != 0 {
		return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return copyDirectory(srcPrefix, path, dest, info, flags)
			} else {
				return copyFile(srcPrefix, path, dest, info, flags)
			}
		})
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("src is directory, CopyRecursive is required")
	}
	return copyFile(srcPrefix, src, dest, info, flags)
}

func copyFile(srcPrefix, src, dest string, info os.FileInfo, flags FlagsCopy) error {
	destPart, err := filepath.Rel(srcPrefix, src)
	if err != nil {
		return err
	}

	destFile, err := os.OpenFile(filepath.Join(dest, destPart), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, DefaultFileMode)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if CopyPreserveMode&flags != 0 {
		if err = destFile.Chmod(info.Mode()); err != nil {
			return err
		}
	}
	if CopyPreserveOwner&flags != 0 && runtime.GOOS != "windows" {
		unixStat, ok := info.Sys().(*unix.Stat_t)
		if ok {
			if err := destFile.Chown(int(unixStat.Uid), int(unixStat.Gid)); err != nil {
				return err
			}
		}
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func copyDirectory(srcPrefix, src, dest string, info os.FileInfo, flags FlagsCopy) error {
	destPart, err := filepath.Rel(srcPrefix, src)
	if err != nil {
		return err
	}
	mode := os.ModeDir
	if CopyPreserveMode&flags != 0 {
		mode |= info.Mode()
	} else {
		mode |= DefaultDirMode
	}
	return os.MkdirAll(filepath.Join(dest, destPart), mode)
}
