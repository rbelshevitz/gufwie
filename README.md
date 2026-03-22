# gufwie

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
