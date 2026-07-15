<p align="center">
  <img src="img/logo.svg" alt="ENCDEC logo" width="220"/>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.17%2B-00ADD8?logo=go&logoColor=white" alt="Go 1.17+"/>
  <img src="https://img.shields.io/badge/release-v2.0.6-0ea5e9" alt="Release v2.0.6"/>
  <img src="https://img.shields.io/badge/license-MIT-14b8a6" alt="License MIT"/>
  <img src="https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-1e293b" alt="Platform: Windows, macOS, Linux"/>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/cipher-AES--256--GCM-0ea5e9" alt="Cipher: AES-256-GCM"/>
  <img src="https://img.shields.io/badge/KDF-Argon2id%20%C2%B7%20SHA--256-6d28d9" alt="KDF: Argon2id and SHA-256"/>
  <img src="https://img.shields.io/badge/modes-machine--bound%20%C2%B7%20passphrase-334155" alt="Modes: machine-bound and passphrase"/>
  <img src="https://img.shields.io/badge/network-100%25%20offline-14b8a6" alt="100% offline"/>
</p>

# ENCDEC (encrypt/decrypt string)

CLI tool written in Go to encrypt/decrypt strings using **AES-GCM**.

It supports **two modes**:

1. **Machine-bound mode (default)**: the encryption key is deterministically derived from:
   - a secret prefix (configurable via env var)
   - Machine ID
   - OS  
   This means ciphertext is typically decryptable only on the *same machine*.

2. **Passphrase mode (portable, optional)**: the encryption key is derived from a passphrase using **Argon2id** with a random salt.
   Ciphertext is self-contained and can be decrypted on any machine as long as you have the passphrase.

> Security note: this is **not** a password hashing tool. If you need to store user passwords, use a dedicated password hashing algorithm (Argon2id / bcrypt / scrypt).

---

## Features

- Cross-platform (Windows / macOS / Linux)
- Authenticated encryption: **AES-GCM**
- Two modes:
  - Machine-bound key derivation (**SHA-256**, 32 bytes → AES-256)
  - Passphrase-based derivation (**Argon2id + random salt**, AES-256)
- Optional file logging (best-effort; failures do not stop the CLI)
- Input size limits (default **1 MiB**) to mitigate trivial local DoS

---

## Requirements

- Go **1.17+** (as per `go.mod`)

Dependencies are managed via Go modules (`go.mod`).

---

## Project layout

- `main.go`: CLI entrypoint and command dispatch
- `lib/lib.go`: crypto and key-derivation library (AES-GCM, machine-bound key, passphrase key)

---

## Usage

The CLI expects **2 arguments** after the program name:

```txt
usage: encdec ENC <plain-text> | encdec DEC <hex-ciphertext> | encdec ENCPASS <plain-text> | encdec DECPASS <p1:ciphertext>
```

Accepted operations:

- `ENC`     encrypt (machine-bound mode)
- `DEC`     decrypt (machine-bound mode)
- `ENCPASS` encrypt (passphrase mode, requires `ENCDEC_PASS`)
- `DECPASS` decrypt (passphrase mode, requires `ENCDEC_PASS`)

If an error occurs, the program prints the error on **stderr** and exits with **code 1**.

---

## Installation / Run

From the module folder (`encdec/`):

```sh
go mod tidy
```

### Run directly with `go run`

Examples below use `go run . ...` but the same arguments work with the built binary.

---

## Encrypt / Decrypt (machine-bound)

### Encrypt

```sh
go run . ENC 'PasswordToEncrypt'
```

Output:

```txt
encrypted : <hex-ciphertext>
```

Notes:
- The output ciphertext is **hex-encoded**.
- Internally the ciphertext contains the **nonce** prepended to the AES-GCM output (this is handled transparently by the tool).

### Decrypt

```sh
go run . DEC <hex-ciphertext>
```

Output:

```txt
PasswordToEncrypt
```

Notes:
- Decryption succeeds only if performed on the **same machine** (same Machine ID + OS) and with the same `ENCDEC_SECRET_PREFIX` used at encryption time.

---

## Passphrase mode (portable) (optional)

In this mode the passphrase is read from environment variable `ENCDEC_PASS`
(to avoid passing secrets through CLI args / history / process list).

### Set passphrase

```sh
export ENCDEC_PASS='your-strong-passphrase'
```

### Encrypt

```sh
go run . ENCPASS 'PasswordToEncrypt'
```

Output:

```txt
encrypted : p1:<salt_b64>:<hex-ciphertext>
```

### Decrypt

```sh
go run . DECPASS 'p1:<salt_b64>:<hex-ciphertext>'
```

Output:

```txt
PasswordToEncrypt
```

Passphrase ciphertext format details (as implemented in `lib/lib.go`):

- Prefix: `p1:` (format/version marker)
- Full format: `p1:<salt_b64>:<cipher_hex>`
- `<salt_b64>` is encoded with `base64.RawStdEncoding`
- Salt length is **16 bytes** (randomly generated on encryption)
- The derived key uses Argon2id with parameters:
  - iterations: `1`
  - memory: `64*1024` KiB (≈ **64 MiB**)
  - parallelism: `4`
  - output key length: `32` bytes (AES-256)

If the ciphertext does not start with `p1:` the tool returns:

```txt
invalid passphrase ciphertext format (missing p1: prefix)
```

---

## Building for source

### 1. (Recommended) change the secret prefix

For an extra layer of security in **machine-bound mode**, change the secret prefix.

- Preferred (no rebuild): set the env var at runtime
  ```sh
  export ENCDEC_SECRET_PREFIX='my-secret-prefix'
  ```
- Or compile a custom default into the binary by editing `defaultSecretKeyPrefix` in `lib/lib.go`:
  ```go
  const defaultSecretKeyPrefix = "jY-1"
  ```

### 2. (Optional) change the log path

If you want file logging at a different location, edit `pathLog` in `main.go`:

```go
const pathLog = "/opt/frm/writable/logs/"
```

### 3. Build the binary

> Note: the project is now a multi-file Go package (`main.go` + `lib/`), so build the
> **package** (`.`), not a single file. `go build main.go` from older versions no longer works.

```sh
go build -o encdec .
```

### Cross-compilation

General form:

```sh
env GOOS=<target-OS> GOARCH=<target-architecture> go build -o encdec .
```

> `env` is a Linux/Unix command. On Windows use **Git Bash**, or set the variables with
> PowerShell / `set` (see Windows example below).

**Linux example (amd64):**
```sh
env GOOS=linux GOARCH=amd64 go build -o encdec .
```

**macOS examples:**
```sh
# Intel Macs
env GOOS=darwin GOARCH=amd64 go build -o encdec .

# Apple Silicon (M1/M2/M3...)
env GOOS=darwin GOARCH=arm64 go build -o encdec .
```

**Windows examples (produces `encdec.exe`):**

From Git Bash / Linux / macOS:
```sh
env GOOS=windows GOARCH=amd64 go build -o encdec.exe .
```

From Windows PowerShell:
```powershell
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o encdec.exe .
```

From Windows `cmd.exe`:
```bat
set GOOS=windows
set GOARCH=amd64
go build -o encdec.exe .
```

> Reminder: on Windows the fixed Unix log path `/opt/...` is not used, so file logging is
> skipped automatically (see [Logging](#logging)).

### Table of config contents (GOOS / GOARCH)

Common combinations supported by the Go toolchain:

| GOOS (Target OS) | GOARCH (Target architecture) |
| ---------------- |:----------------------------:|
| android          | arm                          |
| darwin           | amd64                        |
| darwin           | arm64                        |
| dragonfly        | amd64                        |
| freebsd          | 386                          |
| freebsd          | amd64                        |
| freebsd          | arm                          |
| linux            | 386                          |
| linux            | amd64                        |
| linux            | arm                          |
| linux            | arm64                        |
| linux            | ppc64                        |
| linux            | ppc64le                      |
| linux            | mips                         |
| linux            | mipsle                       |
| linux            | mips64                       |
| linux            | mips64le                     |
| netbsd           | 386                          |
| netbsd           | amd64                        |
| netbsd           | arm                          |
| openbsd          | 386                          |
| openbsd          | amd64                        |
| openbsd          | arm                          |
| plan9            | 386                          |
| plan9            | amd64                        |
| solaris          | amd64                        |
| windows          | 386                          |
| windows          | amd64                        |

> Run `go tool dist list` to see the full, up-to-date list for your Go version.

---

## Configuration / Environment Variables

### `ENCDEC_SECRET_PREFIX` (machine-bound mode)

Overrides the default secret prefix used in machine-bound key derivation.

- If not set, the code uses the built-in default prefix (see `defaultSecretKeyPrefix` in `lib/lib.go`).
- If you keep the default value and the code is public, the prefix should be considered **not secret**.

Example:

```sh
export ENCDEC_SECRET_PREFIX='my-secret-prefix'
```

Important:
- To decrypt successfully in machine-bound mode, `ENCDEC_SECRET_PREFIX` must match the value used during encryption.

### `ENCDEC_PASS` (passphrase mode)

Required for `ENCPASS` / `DECPASS`.

Example:

```sh
export ENCDEC_PASS='your-strong-passphrase'
```

---

## Logging

The CLI tries to write logs to:

```go
const pathLog = "/opt/frm/writable/logs/"
```

Behavior (best-effort):

- The tool attempts to create the directory and append to:  
  `/opt/frm/writable/logs/enc_dec.log`
- If any step fails (directory not writable/missing permissions/path not valid), the program continues **without** file logging.

Platform notes (as implemented in `main.go`):

- On **Windows**, file logging is skipped because the fixed Unix path `/opt/...` is not suitable.
- On non-Windows systems, the tool attempts to enforce file permissions `0664` on the log file when possible.

---

## Input size limits

To reduce accidental/malicious resource usage, inputs are limited:

Library limits (`lib/lib.go`):

- plaintext max: **1 MiB**
- ciphertext max: **1 MiB** of hex chars

CLI-level additional check (`main.go`):

- for `DEC`, ciphertext max: **1 MiB** of hex chars (secondary defense)

If exceeded, the tool returns an error similar to:

```txt
plaintext too large (max 1048576 bytes)
```

or:

```txt
ciphertext too large (max 1048576 hex chars)
```

---

## Security model (important)

- **Machine-bound mode** prevents “copy & decrypt on another machine” but does **not** protect against an attacker already running on the same host.
- **Passphrase mode** provides portability and shifts security to the strength/handling of the passphrase.

---

## License

GNU General Public License v3.0

## Author

Danilo Ritarossi
