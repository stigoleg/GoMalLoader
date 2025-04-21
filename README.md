



# 🛠️ GoMalLoader: Modular Malware Loader in Go

# ⚠️ WARNING

**This project is a modular malware loader capable of executing, injecting, and reflectively loading arbitrary code. Running this code on your main computer or any production system can be extremely dangerous and may harm your device, data, or network.**

- Only run GoMalLoader in a disposable, isolated virtual machine or test environment.
- Never use this tool on systems you do not own or have explicit permission to test.
- Use for educational, research, or red team purposes only, and always follow all applicable laws and ethical guidelines.

---

---

## Overview

**GoMalLoader** is a modular, cross-platform, and production-grade malware loader written in Go. It supports advanced payload delivery and evasion techniques, with a focus on modularity, extensibility, and cross-platform support (Windows, Linux, Mac).

---

## About This Project

This is a personal project created to better understand how advanced malware loaders work, including their architecture, evasion techniques, and cross-platform payload delivery. The codebase is intended for educational and research purposes, to study and experiment with loader design, not for malicious use.

---

## In-Depth Documentation

### Features & Capabilities

- **Shellcode Execution**
  - Load and execute raw shellcode in the current process.
  - Supports both local file and remote (HTTP) payload sources.
  - AES-encrypted payloads with optional obfuscation.
  - Suitable for running custom shellcode, C2 stagers, or proof-of-concept payloads.

- **Remote Process Injection**
  - Inject shellcode into a running target process (e.g., notepad.exe, or any PID on Linux/Mac).
  - Platform-specific injection techniques:
    - Windows: Uses OpenProcess, VirtualAllocEx, WriteProcessMemory, CreateRemoteThread.
    - Linux: Uses ptrace, remote mmap, process_vm_writev, and register hijacking.
    - Mac: Uses Mach APIs (task_for_pid, mach_vm_allocate, mach_vm_write, thread_create).
  - Suitable for EDR evasion, privilege escalation, or lateral movement scenarios.

- **Reflective Library Loader**
  - Loads a shared library (DLL/SO/dylib) into memory and executes its entry point.
  - Platform-specific strategies:
    - Windows: Manual PE parsing and in-memory DLL loading (true reflective loader).
    - Linux: In-memory ELF loading using memfd_create and dlopen (no disk writes).
    - Mac: Mach-O loading using a temp file and dlopen (temp file deleted after loading).
  - Suitable for fileless persistence, in-memory plugins, or advanced red team operations.

- **EDR Evasion & Anti-Analysis**
  - Mutex locking to prevent multiple instances.
  - Sleep skew detection and sandbox name checks.
  - Optional self-delete routine after execution (platform-specific implementation).

- **Encryption & Obfuscation**
  - AES-CBC encryption for all payloads.
  - Optional obfuscation (can be extended, e.g., AES + base64).

---

### When & How to Use Each Feature

- **Shellcode Mode**
  - Use when you want to execute custom shellcode in the current process context.
  - Suitable for initial access payloads, C2 stagers, or running PoC code.
  - Easiest to use and most reliable across all platforms.

- **Inject Remote Mode**
  - Use when you want to inject code into another process for stealth, privilege escalation, or evasion.
  - Requires knowledge of the target process (name or PID).
  - On Linux/Mac, you must run as root or with appropriate entitlements.
  - More likely to trigger EDR/AV alerts; use in controlled environments.

- **Reflective Loader Mode**
  - Use when you want to load a shared library in-memory without touching disk (Windows/Linux) or with minimal disk footprint (Mac).
  - Suitable for fileless persistence, in-memory plugins, or advanced red team operations.
  - On Mac, see the note below about in-memory Mach-O loading limitations.

---

### Step-by-Step Usage Instructions

#### 1. Build the Loader
```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o loader.exe
# Linux
GOOS=linux GOARCH=amd64 go build -o loader_linux
# Mac
GOOS=darwin GOARCH=amd64 go build -o loader_mac
```

#### 2. Generate a Payload
- For shellcode: Use msfvenom, Donut, or your own tool.
- For DLL/SO/dylib: Use your own compiled library.

Example (Windows shellcode):
```bash
msfvenom -p windows/x64/messagebox TEXT="hello" -f raw -o raw_shellcode.bin
```

#### 3. Encrypt the Payload
```bash
python build.py --in raw_shellcode.bin --out encrypted_shellcode.bin --key 0123456789abcdef
```

#### 4. Serve Over HTTP (Optional)
```bash
python -m http.server 8000
```

#### 5. Configure the Loader
Edit `config.json`:
```json
{
  "mode": "inject_remote",
  "source": "url",
  "url": "http://localhost:8000/encrypted_shellcode.bin",
  "path": "payloads/shellcode.bin",
  "aes_key": "0123456789abcdef",
  "target_process": "notepad.exe", // or PID for Linux/Mac
  "obfuscated": true,
  "self_delete": true
}
```

#### 6. Run the Loader
```bash
# Example for Linux remote injection (must be root)
sudo ./loader_linux
```

#### 7. Observe and Troubleshoot
- Check logs for step-by-step status and errors.
- Use in a controlled, isolated environment.
- Execution may trigger antivirus or EDR systems.

---

### Platform-Specific Notes

- **Windows**: All features are fully supported, including true in-memory reflective DLL loading.
- **Linux**: All features are supported. Reflective loader uses memfd_create and dlopen for in-memory SO loading.
- **Mac**: All features are supported. Reflective loader uses a temp file and dlopen for Mach-O dylib loading.

#### ⚠️ Note on Pure In-Memory Mach-O Loading (Mac)

> **macOS does not natively support pure in-memory Mach-O loading.**
>
> - The system dynamic loader (`dyld`) and `dlopen` require a file path.
> - True in-memory Mach-O loading would require a custom loader to parse, map, and resolve Mach-O binaries in memory, which is not supported by public APIs and is a significant engineering effort.
> - All known "in-memory Mach-O loader" projects are experimental, fragile, and may break with new macOS releases.
> - For most use cases, writing to a temp file and using `dlopen` is the most robust and portable solution.
>
> **Future Work:**
> - If pure in-memory Mach-O loading is required, a custom loader must be implemented in Go or C, which parses Mach-O headers, maps segments, resolves symbols, and jumps to the entry point. This is not currently implemented in GoMalLoader.

---

## Features

| Feature                          | Description                                                  |
|----------------------------------|--------------------------------------------------------------|
| `shellcode`                      | Load and execute raw shellcode in current process           |
| `inject_remote`                  | Inject shellcode into a target process (Windows, Linux, Mac) |
| `dll_reflective`                 | Reflectively load a library (DLL/SO/dylib) in memory        |
| AES decryption                   | Payloads can be encrypted via AES-CBC                       |
| Configurable source              | Payloads can be loaded from disk or URL                     |
| EDR evasion                      | Includes mutex locking, sleep skew detection, sandbox name checks |
| Optional self-delete             | Removes executable after execution (platform-specific)       |

---

## Architecture

- **Interface-driven, modular design**: All loader modes and utilities are abstracted via interfaces, with platform-specific implementations selected at build time.
- **Platform-specific files**: Windows, Linux, and Mac implementations are separated using Go build tags and file suffixes.
- **Loader Modes**:
  - **Shellcode**: Executes shellcode in the current process (all platforms).
  - **Remote Injection**: Injects shellcode into a remote process (all platforms, with platform-specific techniques).
  - **Reflective Loader**: Loads a shared library (DLL/SO/dylib) in memory (see below for platform details).
- **Utilities**: AES decryption, evasion, memory allocation, and process/thread manipulation are all abstracted and implemented per platform.

---

## Cross-Platform Reflective Loader Strategies

- **Windows**: Manual PE parsing and in-memory DLL loading (true reflective loader).
- **Linux**: Uses `memfd_create` to create an anonymous in-memory file, writes the ELF shared object (SO) payload, and loads it using `dlopen` via cgo. This achieves in-memory loading without writing to disk.
- **Mac**: Writes the Mach-O dylib payload to a temp file, loads it using `dlopen` via cgo, and deletes the temp file after loading. This is the most robust and compatible approach for now.

### ⚠️ Note on Pure In-Memory Mach-O Loading

> **Pure in-memory Mach-O loading (without writing to disk or using a temp file) is not natively supported by macOS.**
>
> - The system dynamic loader (`dyld`) and `dlopen` require a file path.
> - True in-memory Mach-O loading would require manual parsing, mapping, relocation, and symbol resolution, and is not supported by public APIs.
> - All known "in-memory Mach-O loader" projects are experimental, fragile, and may break with new macOS releases.
> - For most use cases, writing to a temp file and using `dlopen` is the most robust and portable solution.
>
> **Future Work:**
> - If pure in-memory Mach-O loading is required, a custom loader must be implemented in Go or C, which parses Mach-O headers, maps segments, resolves symbols, and jumps to the entry point. This is a significant engineering effort and is not currently implemented in GoMalLoader.

---

## Usage

### 1. Build the loader
```bash
go build -o loader.exe # Windows
go build -o loader_linux # Linux
go build -o loader_mac # Mac
```

### 2. Generate a payload
Using `msfvenom` or another tool:
```bash
msfvenom -p windows/x64/messagebox TEXT="hello" -f raw -o raw_shellcode.bin
```

### 3. Encrypt the payload
Use `build.py` to AES-encrypt the raw shellcode or DLL:
```bash
python build.py --in raw_shellcode.bin --out encrypted_shellcode.bin --key 0123456789abcdef
```

### 4. Serve over HTTP (optional)
```bash
python -m http.server 8000
```

### 5. Configure
Edit `config.json` to point to the file or URL, and choose the mode of operation.

Example:
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

---

## Modes

- `shellcode` – Runs shellcode directly inside the current process.
- `inject_remote` – Injects shellcode into a running target process.
- `dll_reflective` – Reflectively loads and executes a DLL/SO/dylib using platform-specific techniques.

---

## Self-delete

If `"self_delete": true`, the loader deletes itself from disk using a platform-specific method (e.g., `cmd.exe` on Windows, shell script on Linux/Mac).

---

## Notes

- The loader assumes 64-bit payloads.
- Reflective DLLs/SOs/dylibs must contain a known entry point (`RunPayload` or specified RVA).
- Obfuscation is optional and can be extended (e.g., AES + base64).
- Use in isolated test environments. Execution may trigger antivirus or EDR systems.
- **Linux/Mac support for remote injection and reflective loader is now implemented, with platform-specific strategies.**
- **Pure in-memory Mach-O loading is not implemented; see above for details and future work.**


---

## Legal

This project is for **educational and red team** purposes only.  
Do **not** deploy in production or on systems without **explicit permission**.
