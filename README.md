# godirb

[![Go](https://img.shields.io/badge/Go-1.25.1-00ADD8?logo=go)](https://go.dev/) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE) [![Latest release](https://img.shields.io/github/v/release/MyCode83/godirb?sort=semver)](https://github.com/MyCode83/godirb/releases) 

`godirb` is a fast recursive directory/file brute-forcer written in Go.

It is for the moments where you want a modern `dirb`-like tool: run it, get useful results, tune the obvious flags, and avoid dragging a full fuzzing framework into a simple job.

## godirb vs DirSearch

DirSearch is a mature web path scanner. godirb is smaller on purpose.

| Feature | godirb | DirSearch |
| --- | --- | --- |
| Find files and folders | ✅ | ✅ |
| Recursive scan | ✅ | ✅ |
| Use custom wordlists | ✅ | ✅ |
| Made in Go | ✅ | ❌ |
| Works as a single binary | ✅ | ❌ |
| Baseline filter with heuristics | ✅ | ❌ |
| Embedded wordlists as Go slices | ✅ | ❌ |
| Default embedded `medium` wordlist | ✅ | ❌ |
| Basic scan without runtime wordlist files | ✅ | ❌ |
| Port fuzzing: `http://host:FUZZ` | ✅ | ❌ |

## Features

- Embedded wordlists: `small`, `common`, `medium`, `big`, `ports`, `payloads`, `xss`, `lfi`
- Default wordlist: `medium`
- Recursive mode with `-r, --recursive`
- Extensions with `-x, --ext`
- Threads with `-t, --threads` (default: `15`)
- Ignore status codes with `-i, --ignore` (default: `404,400,405,408`)
- Wildcard detection for directory scans
- Text, quiet text, JSON, CSV and file output

## Install

```bash
go install github.com/MyCode83/godirb@latest
```

Or download a binary from [Releases](https://github.com/MyCode83/godirb/releases), or build it:

```bash
git clone https://github.com/MyCode83/godirb.git
cd godirb
go build -o godirb .
```

## Usage

```bash
godirb -u http://localhost
godirb -u http://localhost -r
godirb -u http://localhost -w ./paths.txt
godirb -u http://localhost -t 30
godirb -u http://localhost -i 404,403,500
godirb -u http://localhost -x php,txt,bak
```

`FUZZ` in the URL switches to placeholder mode:

```bash
godirb -u "http://localhost/search?q=FUZZ" -w payloads
godirb -u http://localhost:FUZZ
```

## Output

```text
[DIR] http://localhost/admin ---> 200 | 1234
```

```bash
godirb -u http://localhost --json -o results.json
godirb -u http://localhost --csv -o results.csv
godirb -u http://localhost -q
```

## Disclaimer

Use godirb only for authorized testing, labs and CTFs. You are responsible for having permission to scan a target.

## License

MIT. See [LICENSE](LICENSE).
