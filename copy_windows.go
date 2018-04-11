// +build windows

package fshelper

import "os"

func preserveOwner(destFile *os.File, srcStat os.FileInfo) {
	return
}
