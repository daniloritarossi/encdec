//go:build windows

package main

// enforceLogPerm is a no-op on Windows: Unix file-mode/owner semantics
// (syscall.Stat_t, chmod) do not apply there.
func enforceLogPerm(_ string) {}
