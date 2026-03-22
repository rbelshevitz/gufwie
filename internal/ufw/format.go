package ufw

import (
	"strings"
)

func FormatCommand(action Action, args []string) string {
	argv := append([]string{"ufw", string(action)}, args...)
	return joinQuoted(argv)
}

func FormatUfwArgs(args []string) string {
	argv := append([]string{"ufw"}, args...)
	return joinQuoted(argv)
}

func joinQuoted(argv []string) string {
	var b strings.Builder
	for i, a := range argv {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(quoteArg(a))
	}
	return b.String()
}

func quoteArg(s string) string {
	if s == "" {
		return "''"
	}
	needs := false
	for _, r := range s {
		if r == '\'' || r == '"' || r == '\\' || r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			needs = true
			break
		}
	}
	if !needs {
		return s
	}
	// POSIX-ish single-quote escaping: ' becomes '\''.
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

