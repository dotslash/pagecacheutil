package main

import (
	"flag"
	"fmt"
	"github.com/dotslash/pagecacheutil/oscompat"
	. "github.com/dotslash/pagecacheutil/util"
	"github.com/fatih/color"
	"os"
	"syscall"
	"unsafe"
)

var verbose = flag.Bool("v", false, "Verbose mode")
var evict = flag.Bool("e", false, "Evict the file")
var touch = flag.Bool("t", false, "Touch the file")

func handleFile(fname string) {
	file, err, fstat := openFile(fname)
	DieOnErr(err)
	defer tryClose(file)

	memAddr, err := mmapSyscall(fstat, file)
	DieOnErr(err)

	defer tryMunmapSyscall(memAddr, fstat, fname)

	isPageCached, err := mincoreSyscall(memAddr, fstat)
	DieOnErr(err)
	printMincoreRes(isPageCached, fstat)

	if *touch {
		touchFile(fstat, memAddr)
	}
	if *evict {
		oscompat.EvictFile(file, fstat, memAddr)
	}

	if *touch || *evict {
		isPageCached, err = mincoreSyscall(memAddr, fstat)
		DieOnErr(err)
		printMincoreRes(isPageCached, fstat)
	}

}

func printMincoreRes(isPageCached []bool, fstat os.FileInfo) {
	inCache := 0
	greenBg := color.New(color.BgHiGreen, color.FgBlack)
	yellowBg := color.New(color.BgYellow, color.FgBlack)
	redBg := color.New(color.BgWhite, color.FgBlack)
	reset := color.New(color.Reset)
	width := 50
	if len(isPageCached) < width {
		width = len(isPageCached)
	}
	groupSize := len(isPageCached) / width
	for i := 0; i < len(isPageCached); i += groupSize {
		totInGroup := 0
		hitsInGroup := 0
		for ; totInGroup < groupSize && i+totInGroup < len(isPageCached); totInGroup++ {
			if isPageCached[i+totInGroup] {
				hitsInGroup++
			}
		}
		inCache += hitsInGroup
		if *verbose {
			if totInGroup == hitsInGroup {
				greenBg.Print("O")
			} else if hitsInGroup != 0 {
				yellowBg.Print("o")
			} else {
				redBg.Print("x")
			}
		}
	}
	if *verbose {
		reset.Println()
	}
	fmt.Printf("Pages in cache for %v: %v/%v\n", fstat.Name(), inCache, len(isPageCached))
}

func tryMunmapSyscall(memAddr uintptr, fstat os.FileInfo, fname string) {
	_, _, err := syscall.Syscall(
		syscall.SYS_MUNMAP,
		memAddr,
		uintptr(fstat.Size()),
		0,
	)
	if err != 0 {
		fmt.Printf("MUNMAP failed for %v with err:%v \n", fname, err)
	}
}

func mmapSyscall(fstat os.FileInfo, file *os.File) (uintptr, error) {
	memAddr, _, err := syscall.Syscall6(
		syscall.SYS_MMAP, 0,
		uintptr(fstat.Size()),
		syscall.PROT_READ,
		syscall.MAP_SHARED,
		file.Fd(),
		0,
	)
	if err != 0 {
		return 0, err
	}
	return memAddr, nil
}

func mincoreSyscall(mmapAddr uintptr, fstat os.FileInfo) ([]bool, error) {
	pageSize := int64(syscall.Getpagesize())
	numPages := (fstat.Size() + pageSize - 1) / pageSize
	mincoreRes := make([]byte, numPages)
	_, _, err := syscall.Syscall(
		syscall.SYS_MINCORE,
		mmapAddr,
		uintptr(fstat.Size()),
		uintptr(unsafe.Pointer(&mincoreRes[0])),
	)
	if err != 0 {
		return nil, err
	}
	isPageCached := make([]bool, numPages)
	for i := range mincoreRes {
		isPageCached[i] = mincoreRes[i] != 0
	}
	return isPageCached, nil

}

func openFile(fname string) (*os.File, error, os.FileInfo) {
	f, err := os.Open(fname)
	DieOnErr(err)

	stat, err := f.Stat()
	DieOnErr(err)
	return f, err, stat
}

func tryClose(f *os.File) {
	if err := f.Close(); err != nil {
		fmt.Printf("Failed to close %v err:%v\n", f.Name(), err)
	}
}

func touchFile(fstat os.FileInfo, memAddr uintptr) {
	dummy, mod := 0, (10000*10000)+7

	// Slice memory layout
	// Taken from syscall/syscall_unix.go
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{memAddr, int(fstat.Size()), int(fstat.Size())}
	membytes := *(*[]byte)(unsafe.Pointer(&sl))

	for _, membyte := range membytes {
		dummy = (dummy + int(membyte)) % mod
	}
	fmt.Printf("touched %v contentHash:%v\n", fstat.Name(), dummy)
}

func main() {
	flag.Parse()
	files := flag.Args()
	for _, f := range files {
		handleFile(f)
	}
}
