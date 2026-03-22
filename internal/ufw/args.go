package ufw

import (
	"fmt"
	"strings"
	"unicode"
)

// SplitArgs splits a shell-like argument line into arguments.
// Supports simple single and double quotes and backslash escapes.
func SplitArgs(line string) ([]string, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty args")
	}

	var args []string
	var b strings.Builder
	inSingle := false
	inDouble := false
	escaped := false

	flush := func() {
		if b.Len() == 0 {
			return
		}
		args = append(args, b.String())
		b.Reset()
	}

	for i, r := range line {
		if escaped {
			b.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' && !inSingle {
			escaped = true
			continue
		}
		if r == '\'' && !inDouble {
			inSingle = !inSingle
			continue
		}
		if r == '"' && !inSingle {
			inDouble = !inDouble
			continue
		}
		if !inSingle && !inDouble && unicode.IsSpace(r) {
			flush()
			continue
		}
		b.WriteRune(r)
		if i == len(line)-1 {
			// handled by final flush below
		}
	}
	if escaped {
		return nil, fmt.Errorf("dangling escape")
	}
	if inSingle || inDouble {
		return nil, fmt.Errorf("unterminated quote")
	}
	flush()
	return args, nil
}

