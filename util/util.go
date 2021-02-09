package util

import "syscall"

func DieOnErr(err error) {
	if err != nil && err != syscall.Errno(0) {
		panic(err)
	}
}
