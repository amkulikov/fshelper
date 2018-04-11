// +build !windows

package fshelper

import (
	"os"
	"golang.org/x/sys/unix"
)

func preserveOwner(destFile *os.File, srcStat os.FileInfo) {
	unixStat, ok := srcStat.Sys().(*unix.Stat_t)
	if ok {
		destFile.Chown(int(unixStat.Uid), int(unixStat.Gid))
	}
	return
}
