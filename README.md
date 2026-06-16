# ENCDEC (encrypt/decrypt string)

CLI tool written in Go to encrypt/decrypt strings using AES-GCM.

It supports **two modes**:

1. **Machine-bound mode (default)**: the encryption key is deterministically derived from:
   - a secret prefix (configurable)
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
- “Best-effort” logging (does not block encryption/decryption if logging fails)
- Input size limits (default ~1 MiB) to mitigate trivial local DoS

---

## Requirements

- Go **1.17+** (tested on Go **1.18**)

Dependencies are managed via Go modules (`go.mod`).

---

## Installation / Run

From the module folder (`encdec/`):

```sh
go mod tidy
```

### Encrypt / Decrypt (machine-bound)

Encrypt:

```sh
go run . ENC 'PasswordToEncrypt'
```

Output:

```txt
encrypted : <hex-ciphertext>
```

Decrypt:

```sh
go run . DEC <hex-ciphertext>
```

Output:

```txt
PasswordToEncrypt
```

---

## Passphrase mode (portable) (optional)

In this mode the passphrase is read from environment variable `ENCDEC_PASS`
(to avoid passing secrets through CLI args / history / process list).

Set passphrase:

```sh
export ENCDEC_PASS='your-strong-passphrase'
```

Encrypt:

```sh
go run . ENCPASS 'PasswordToEncrypt'
```

Output:

```txt
encrypted : p1:<salt_b64>:<hex-ciphertext>
```

Decrypt:

```sh
go run . DECPASS 'p1:<salt_b64>:<hex-ciphertext>'
```

Output:

```txt
PasswordToEncrypt
```

Notes:
- The `p1:` prefix identifies the passphrase ciphertext format version.
- The salt is encoded using `base64.RawStdEncoding`.

---

## Building

Build binary:

```sh
go build -o encdec .
```

Cross-compile example (Linux amd64):

```sh
GOOS=linux GOARCH=amd64 go build -o encdec .
```

---

## Configuration / Environment Variables

### `ENCDEC_SECRET_PREFIX` (machine-bound mode)

Override the default secret prefix used in machine-bound key derivation:

```sh
export ENCDEC_SECRET_PREFIX='my-secret-prefix'
```

If you keep the default value, encryption/decryption will still work, but the prefix is not secret if the code is public.

### `ENCDEC_PASS` (passphrase mode)

Required for `ENCPASS` / `DECPASS`:

```sh
export ENCDEC_PASS='your-strong-passphrase'
```

---

## Logging

The CLI tries to write logs to:

```go
const pathLog = "/opt/frm/writable/logs/"
```

Logging is **best-effort**:
- if the path is not writable or not suitable for the current OS, the program continues without file logging.

---

## Input size limits

To reduce accidental/malicious resource usage, inputs are limited:
- plaintext max: ~1 MiB
- ciphertext max: ~1 MiB (hex chars)

---

## Security model (important)

- **Machine-bound mode** prevents “copy & decrypt on another machine” but does **not** protect against an attacker already running on the same host.
- **Passphrase mode** provides portability and shifts security to the strength/handling of the passphrase.

---

## License

GNU General Public License v3.0

## Author

Danilo Ritarossi
