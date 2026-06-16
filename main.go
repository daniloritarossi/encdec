package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"encdenc/lib"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

const pathLog = "/opt/frm/writable/logs/"

func main() {
	if err := run(os.Args, os.Stdout, os.Stderr); err != nil {
		// keep non-zero exit code for CLI usage
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr *os.File) error {
	// Optional logging to file (best-effort). Avoid failing the tool if the log path is not available.
	setupLogging()

	if len(args) < 3 {
		return errors.New("usage: encdec ENC <plain-text> | encdec DEC <hex-ciphertext> | encdec ENCPASS <plain-text> | encdec DECPASS <p1:ciphertext>")
	}

	op := args[1]
	input := args[2]

	// Limit ciphertext size at CLI level as a second line of defense.
	// (lib layer also enforces bounds)
	const maxCiphertextHexLen = 1 << 20 // 1 MiB of hex
	if op == "DEC" && len(input) > maxCiphertextHexLen {
		return fmt.Errorf("ciphertext too large (max %d hex chars)", maxCiphertextHexLen)
	}

	switch op {
	case "ENC":
		key, err := lib.GenerateKey()
		if err != nil {
			return err
		}
		encrypted, err := lib.EncryptString(input, key)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "encrypted : %s\n", encrypted)
		return nil
	case "DEC":
		key, err := lib.GenerateKey()
		if err != nil {
			return err
		}
		decrypted, err := lib.DecryptString(input, key)
		if err != nil {
			return err
		}
		fmt.Fprintln(stdout, decrypted)
		return nil
	case "ENCPASS":
		pass := os.Getenv("ENCDEC_PASS")
		if pass == "" {
			return errors.New("ENCDEC_PASS env var is required for ENCPASS/DECPASS")
		}
		encrypted, err := lib.EncryptStringWithPassphrase(input, pass)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "encrypted : %s\n", encrypted)
		return nil
	case "DECPASS":
		pass := os.Getenv("ENCDEC_PASS")
		if pass == "" {
			return errors.New("ENCDEC_PASS env var is required for ENCPASS/DECPASS")
		}
		decrypted, err := lib.DecryptStringWithPassphrase(input, pass)
		if err != nil {
			return err
		}
		fmt.Fprintln(stdout, decrypted)
		return nil
	default:
		return errors.New("only ENC, DEC, ENCPASS and DECPASS values are accepted")
	}
}

func setupLogging() {
	// This project was designed to log under a fixed path (pathLog).
	// Make it safe/cross-platform and best-effort: if it fails, continue without file logging.
	if pathLog == "" {
		return
	}

	// On Windows, a Unix path like /opt/... will fail.
	if runtime.GOOS == "windows" && filepath.IsAbs(pathLog) && filepath.VolumeName(pathLog) == "" {
		return
	}

	if err := os.MkdirAll(pathLog, 0o775); err != nil {
		return
	}

	logPath := filepath.Join(pathLog, "enc_dec.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o664)
	if err != nil {
		return
	}
	log.SetOutput(file)

	// Permission/owner enforcement is Unix-specific (needs syscall.Stat_t).
	// Implemented per-platform via build tags so the program compiles on Windows too.
	enforceLogPerm(logPath)
}
