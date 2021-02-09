package oscompat

import (
	"github.com/dotslash/pagecacheutil/util"
	"os"
	"syscall"
)

func EvictFile(f *os.File, fstat os.FileInfo, mmapaddr uintptr) {
	// msync(mmapaddr, len, MS_INVALIDATE)
	_, _, err := syscall.Syscall(
		syscall.SYS_MSYNC,
		mmapaddr,
		uintptr(fstat.Size()), // full file => length = file size
		syscall.MS_INVALIDATE, // invalidate the cache
	)
	util.DieOnErr(err)
}
