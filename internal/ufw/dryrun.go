package ufw

type DryRunError struct {
	Cmd string
}

func (e *DryRunError) Error() string {
	if e == nil {
		return "dry-run"
	}
	if e.Cmd == "" {
		return "dry-run"
	}
	return "dry-run: " + e.Cmd
}

