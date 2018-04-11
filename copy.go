package fshelper

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/unix"
)

type FlagsCopy uint

const (
	CopyContent       FlagsCopy = 1 << iota
	CopyRecursive
	CopyPreserveMode
	CopyPreserveOwner
	CopyParents
	CopyVerbose
	copyOneFile
)

const (
	DefaultDirMode  = 0755
	DefaultFileMode = 0644
)

func Copy(flags FlagsCopy, src, dest string) error {
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
			if path == src && CopyContent&flags != 0 {
				return nil
			}
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
	} else {
		return copyFile(srcPrefix, src, dest, info, flags|copyOneFile)
	}
}

func copyFile(srcPrefix, src, dest string, info os.FileInfo, flags FlagsCopy) (err error) {
	relPart, err := filepath.Rel(srcPrefix, src)
	if err != nil {
		return err
	}

	_, destFileName := filepath.Split(dest)
	if copyOneFile&flags == 0 || filepath.Dir(relPart) != "." || destFileName == "" {
		dest = filepath.Join(dest, relPart)
	}

	if CopyVerbose&flags != 0 {
		fmt.Printf("f: %s -> %s ", src, dest)
		defer func() {
			if err == nil {
				fmt.Println(": OK!")
			} else {
				fmt.Printf(": %s\n", err)
			}
		}()
	}

	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, DefaultFileMode)
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

func copyDirectory(srcPrefix, src, dest string, info os.FileInfo, flags FlagsCopy) (err error) {
	destPart, err := filepath.Rel(srcPrefix, src)
	if err != nil {
		return err
	}
	dest = filepath.Join(dest, destPart)

	if CopyVerbose&flags != 0 {
		fmt.Printf("d: %s -> %s ", src, dest)
		defer func() {
			if err == nil {
				fmt.Println(": OK!")
			} else {
				fmt.Printf(": %s\n", err)
			}
		}()
	}

	mode := os.ModeDir
	if CopyPreserveMode&flags != 0 {
		mode |= info.Mode()
	} else {
		mode |= DefaultDirMode
	}
	return os.MkdirAll(dest, mode)
}
