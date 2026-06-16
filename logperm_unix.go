//go:build !windows

package main

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
)

// enforceLogPerm best-effort fixes the log file mode to 0664 when the current
// user owns the file. Unix-only: it relies on syscall.Stat_t for the owner UID.
// Any failure is ignored so logging never breaks the CLI.
func enforceLogPerm(logPath string) {
	fileInfo, err := os.Stat(logPath)
	if err != nil {
		return
	}

	if fileInfo.Mode().Perm() == 0o664 {
		return
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return
	}

	owner, err := user.LookupId(fmt.Sprint(stat.Uid))
	if err != nil {
		return
	}
	userCurrent, err := user.Current()
	if err != nil {
		return
	}
	if owner.Username == userCurrent.Username {
		_ = os.Chmod(logPath, 0o664)
	}
}
