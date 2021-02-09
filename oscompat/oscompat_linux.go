package oscompat

import (
	"github.com/dotslash/pagecacheutil/util"
	"os"
	"syscall"
)

// https://linux.die.net/man/2/fadvise and https://linux.die.net/man/2/madvise have the same
// advice arguments. golang for some reason does not give constants for FADV
const FADV_DONTNEED = syscall.MADV_DONTNEED // == 4

func EvictFile(f *os.File, fstat os.FileInfo, mmapaddr uintptr) {
	// posix_fadvise(fd, offset, len, POSIX_FADV_DONTNEED))
	_, _, err := syscall.Syscall6(
		syscall.SYS_FADVISE64,
		f.Fd(),
		0,                     // full file => 0 offset
		uintptr(fstat.Size()), // full file => length = file size
		FADV_DONTNEED,         // ask the os to discard this from page cache
		0, 0,                  // dummy args
	)
	util.DieOnErr(err)
}
