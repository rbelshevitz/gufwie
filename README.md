# gufwie

![CI](https://github.com/rbelshevitz/gufwie/actions/workflows/ci.yml/badge.svg)
![Security](https://github.com/rbelshevitz/gufwie/actions/workflows/security.yml/badge.svg)

`gufwie` is an `nmtui`-style TUI for **UFW**: browse rules, search, inspect details, and apply common changes with confirmations.

It’s meant for the “I just need to quickly understand what’s open and make one safe change” workflow, without remembering UFW’s exact syntax or dealing with numbered deletes in a shell.

## Requirements

- Linux with `ufw` (on Windows, run inside **WSL2/Ubuntu**)
- Run with privileges (usually `sudo`)
- Tested with `ufw` `0.36.2` (older versions may work)

## Run

From the repo:

```bash
sudo go run ./cmd/gufwie
```

Or install:

```bash
go install github.com/rbelshevitz/gufwie/cmd/gufwie@latest
sudo gufwie
```

Print version:

```bash
gufwie --version
```

## Service/Profile mode

UFW supports *application profiles* (e.g. `OpenSSH`, `Nginx Full`). `gufwie` can add rules by selecting a profile instead of typing ports.

To get more profiles, install software that ships them (commonly `nginx`, `apache2`, `openssh-server`). Profiles live in:

```text
/etc/ufw/applications.d
```

## Keys

- `F1` help, `F5` refresh, `F6` verbose status, `F8` delete, `F10` quit
- `/` search/filter, `F3` search dialog
- `F2` / `a` add (Freeform / Wizard / Service)
- `c` show suggested ufw command, `y` copy it, `Y` copy selected raw rule line
- `Tab` / `Shift-Tab` cycle focus (no mouse required)

Safety: potentially dangerous actions require confirmation and default to **No**.

## Security / audit

This tool needs `sudo` because `ufw` needs it. `gufwie` itself does not open network sockets and only shells out to `ufw` with arguments you approve in the UI.

If you want to audit:

```bash
rg -n "exec\\.Command|Run\\(ctx.*\\\"ufw\\\"" internal
go test ./...
go vet ./...
govulncheck ./...
gosec ./...
```

## Build (linux/amd64)

```bash
GOOS=linux GOARCH=amd64 go build -o gufwie ./cmd/gufwie
```
