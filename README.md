<div align="center">

# godirb

**Fast, modern directory, file, port and FUZZ brute-forcer written in Go.**

Built for quick scans where you want a modern dirb-like tool: run it, get useful results, tune the obvious flags, and avoid dragging a full fuzzing framework into a simple job.

<p>

![License](https://img.shields.io/github/license/MyCode83/godirb?style=for-the-badge)
![Release](https://img.shields.io/github/v/release/MyCode83/godirb?style=for-the-badge)
![Go](https://img.shields.io/github/go-mod/go-version/MyCode83/godirb?style=for-the-badge)
![Stars](https://img.shields.io/github/stars/MyCode83/godirb?style=for-the-badge)

</p>

</div>

## 📦 Installation

### Go

```bash
go install github.com/MyCode83/godirb@latest
```
## Homebrew
```bash
brew install MyCode83/godirb/godirb
```

### Binary

Download the latest release for your platform from the
**Releases** page and add it to your `PATH`.

---

## 🚀 Quick Start

Basic scan

```bash
godirb -u https://example.com
```

Recursive

```bash
godirb -u https://example.com -r
```

Recursive with a depth limit

```bash
godirb -u https://example.com -r --depth 2
```

Custom wordlist

```bash
godirb -u https://example.com -w paths.txt
```

JSON output

```bash
godirb -u https://example.com --json -o results.json
```

---

## ✨ Why godirb?

godirb is designed for the common case: you want to enumerate directories and files quickly, without configuring a large fuzzing framework.

### Highlights

- ⚡ Fast native Go binary
- 📦 Single executable
- 📚 Embedded wordlists
- 🔄 Recursive scanning
- 📂 Directory and file discovery
- 🎯 Wildcard filtering / response calibration
- 🌐 Port fuzzing (`http://host:FUZZ`)
- 📄 JSON & CSV output
- 🧩 Simple CLI

---

## 📊 godirb vs DirSearch

DirSearch is a mature and feature-rich web path scanner.

godirb intentionally focuses on the most common workflow: install, run and get useful results with minimal setup.

| Feature | godirb | DirSearch |
| :--- | :---: | :---: |
| Find files and folders | ✅ | ✅ |
| Recursive scan | ✅ | ✅ |
| Custom wordlists | ✅ | ✅ |
| Written in Go | ✅ | ❌ |
| Single binary | ✅ | ❌ |
| Embedded default wordlists | ✅ | ❌ |
| Works without runtime wordlist files | ✅ | ❌ |
| Port fuzzing (`http://host:FUZZ`) | ✅ | ❌ |

---

## 📦 Features

### Scanning

- Directory and file brute-forcing
- Recursive mode (`-r`, `--recursive`)
- Recursive depth limit (`--depth`, default: 2)
- Extensions (`-x`, `--ext`)
- Custom wordlists (`-w`, `--wordlist`)
- FUZZ placeholder mode

### Embedded Wordlists

- small
- medium *(default)*
- big
- ports
- payloads
- xss
- lfi

### Output

- Standard text
- Quiet mode
- JSON
- CSV
- File output

### Control

- Threads (`-t`, `--threads`)
- Ignore status codes (`-i`, `--ignore`)
- Default ignored codes: `404,400,405,408`
- Wildcard filtering / response calibration

---

## 💻 Examples

Basic scan

```bash
godirb -u https://example.com
```

Recursive

```bash
godirb -u https://example.com -r
```

Recursive with one nested level

```bash
godirb -u https://example.com -r --depth 1
```

Custom wordlist

```bash
godirb -u https://example.com -w paths.txt
```

Extensions

```bash
godirb -u https://example.com -x php,txt,bak
```

FUZZ parameter

```bash
godirb -u "https://example.com/search?q=FUZZ" -w payloads
```

Port fuzzing

```bash
godirb -u https://example.com:FUZZ
```

Export JSON

```bash
godirb -u https://example.com --json -o results.json
```

Export CSV

```bash
godirb -u https://example.com --csv -o results.csv
```

---

## 📋 Example Output

```text
DIR       200      1234 B  https://example.com/admin
FILE      200       842 B  https://example.com/login.php
DIR       403       795 B  https://example.com/uploads
```

---

<details>
<summary><b>📖 Embedded wordlists</b></summary>

| Name | Purpose |
|------|---------|
| small | common.txt from SecLists |
| medium | Default raft-medium-directories |
| big | Larger enumeration DirBuster big |
| ports | Port fuzzing |
| payloads | Generic payloads |
| xss | XSS payloads |
| lfi | LFI payloads |

</details>

---

## ⚠️ Disclaimer

Use **godirb** only for authorized security testing, labs and CTFs.

You are responsible for obtaining permission before scanning any target.

---

## 📄 License

Licensed under the **MIT License**. See **LICENSE** for details.
