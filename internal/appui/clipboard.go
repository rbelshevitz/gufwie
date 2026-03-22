package appui

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func CopyText(ctx context.Context, text string) error {
	// Prefer common Linux clipboard tools first, then WSL interop.
	candidates := [][]string{
		{"wl-copy"},
		{"xclip", "-selection", "clipboard"},
		{"clip.exe"},
	}

	for _, c := range candidates {
		if len(c) == 0 {
			continue
		}
		path, err := exec.LookPath(c[0])
		if err != nil && c[0] == "clip.exe" {
			if p, ok := wslClipPath(); ok {
				path = p
				err = nil
			}
		}
		if err != nil {
			continue
		}
		if err := runClipboard(ctx, path, c[1:], text); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no clipboard helper found (install `wl-clipboard` or `xclip`, or enable WSL interop for clip.exe)")
}

func runClipboard(ctx context.Context, exe string, args []string, text string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Stdin = bytes.NewBufferString(text)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func wslClipPath() (string, bool) {
	// Typical WSL path for Windows' clip.exe
	p := "/mnt/c/Windows/System32/clip.exe"
	if _, err := os.Stat(p); err == nil {
		return p, true
	}
	// Some setups mount drives differently; as a fallback try a best-effort relative check.
	if vol := os.Getenv("SYSTEMROOT"); vol != "" {
		p2 := filepath.Join(vol, "System32", "clip.exe")
		if _, err := os.Stat(p2); err == nil {
			return p2, true
		}
	}
	return "", false
}

