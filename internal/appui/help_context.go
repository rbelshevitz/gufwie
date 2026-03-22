package appui

import "strings"

func (u *ui) setModalHelp(fn func())  { u.modalHelp = fn }
func (u *ui) clearModalHelp()         { u.modalHelp = nil }
func (u *ui) setDialogHelp(fn func()) { u.dialogHelp = fn }
func (u *ui) clearDialogHelp()        { u.dialogHelp = nil }

func profileHelpText() string {
	return strings.Join([]string{
		"Service/Profile mode:",
		"  Adds a rule using UFW application profiles (`ufw app list/info`).",
		"",
		"Keys:",
		"  Tab / Shift-Tab  cycle focus",
		"  /  focus Search",
		"  r or F5  refresh profiles",
		"  Esc  close",
		"",
		"Buttons:",
		"  Info  runs `ufw app info <Profile>`",
		"  Apply applies `ufw <action> <Profile>` (confirm; default No)",
		"  Profiles dir shows `/etc/ufw/applications.d`",
	}, "\n")
}

func searchHelpText() string {
	return strings.Join([]string{
		"Search / Filter:",
		"  Filters visible rules by a substring match.",
		"",
		"Keys:",
		"  Ctrl+U  clear",
		"  Enter   apply & close",
		"  Esc     close",
	}, "\n")
}

func addMenuHelpText() string {
	return strings.Join([]string{
		"Add menu:",
		"  Freeform: type ufw args after the action (fast, flexible).",
		"  Wizard:   structured fields + preview.",
		"  Service/Profile: pick from `ufw app list` profiles.",
	}, "\n")
}

func dialogHelpText() string {
	return strings.Join([]string{
		"Dialog:",
		"  Tab / Shift-Tab moves focus between buttons.",
		"  Enter activates the focused button.",
		"  Esc closes the dialog.",
	}, "\n")
}
