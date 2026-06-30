<div align="center">

# godirb

**Fast, modern recursive directory and file brute-forcer written in Go.**

Built for quick scans where you want a modern dirb-like tool: run it, get useful results, tune the obvious flags, and avoid dragging a full fuzzing framework into a simple job.

<p>

![License](https://img.shields.io/github/license/MyCode83/godirb?style=for-the-badge)
![Release](https://img.shields.io/github/v/release/MyCode83/godirb?style=for-the-badge)
![Go](https://img.shields.io/github/go-mod/go-version/MyCode83/godirb?style=for-the-badge)
![Stars](https://img.shields.io/github/stars/MyCode83/godirb?style=for-the-badge)

</p>

</div>

---

## 🚀 Quick Start

Install:

```bash
go install github.com/MyCode83/godirb@latest
```

Basic scan:

```bash
godirb -u https://example.com
```

Recursive:

```bash
godirb -u https://example.com -r
```

Custom wordlist:

```bash
godirb -u https://example.com -w paths.txt
```

JSON output:

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
- 🎯 Wildcard detection
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
- Extensions (`-x`, `--ext`)
- Custom wordlists (`-w`, `--wordlist`)
- FUZZ placeholder mode

### Embedded Wordlists

- small
- common
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
- Wildcard detection

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
[DIR]  https://example.com/admin       ---> 200 | 1234
[FILE] https://example.com/login.php   ---> 200 | 842
[DIR]  https://example.com/uploads     ---> 403 | 795
```

---

<details>
<summary><b>📖 Embedded wordlists</b></summary>

| Name | Purpose |
|------|---------|
| small | Tiny scans |
| common | medium alias |
| medium | Default |
| big | Larger enumeration |
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
