package ufw

import (
	"fmt"
	"strings"
)

// NotPrivilegedError indicates ufw refused to run due to missing privileges.
// This typically happens when gufwie is launched without sudo.
type NotPrivilegedError struct {
	Args   []string
	Detail string
}

func (e *NotPrivilegedError) Error() string {
	d := strings.TrimSpace(e.Detail)
	if d == "" {
		return fmt.Sprintf("ufw %v failed: insufficient privileges", e.Args)
	}
	return fmt.Sprintf("ufw %v failed: insufficient privileges: %s", e.Args, d)
}

func isNotPrivilegedMessage(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return false
	}
	// Common ufw messages when run without sudo/root.
	if strings.Contains(s, "you need to be root to run this program") {
		return true
	}
	if strings.Contains(s, "need to be root to run this program") {
		return true
	}
	// Some distros/tools may surface a generic permission denial.
	if strings.Contains(s, "permission denied") {
		return true
	}
	return false
}

