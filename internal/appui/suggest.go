package appui

import (
	"context"
	"strings"

	"github.com/rbelshevitz/gufwie/internal/ufw"
)

func (u *ui) showSuggestedCommand() {
	r, ok := u.selectedRule()
	if !ok {
		return
	}
	u.showModal("Suggested command", suggestCommand(r), []string{"OK"}, func(string) {})
}

func (u *ui) copySuggestedCommand() {
	r, ok := u.selectedRule()
	if !ok {
		return
	}
	suggest := suggestCommand(r)
	if err := CopyText(context.Background(), suggest+"\n"); err != nil {
		u.showError(err)
		return
	}
	u.flashHelp("Copied suggested command to clipboard")
}

func suggestCommand(r ufw.Rule) string {
	// Rough mapping from the status line columns back to an approximate ufw command.
	// Intended for copy/paste and tweaks, not guaranteed round-trippable.
	action := ufw.ActionAllow
	dir := ""

	parts := strings.Fields(strings.ToUpper(r.Action))
	if len(parts) > 0 {
		switch parts[0] {
		case "ALLOW":
			action = ufw.ActionAllow
		case "DENY":
			action = ufw.ActionDeny
		case "REJECT":
			action = ufw.ActionReject
		case "LIMIT":
			action = ufw.ActionLimit
		}
	}
	if len(parts) > 1 {
		switch parts[1] {
		case "IN":
			dir = "in"
		case "OUT":
			dir = "out"
		}
	}

	args := []string{}
	if dir != "" {
		args = append(args, dir)
	}
	if strings.TrimSpace(r.From) != "" && !strings.EqualFold(strings.TrimSpace(r.From), "Anywhere") {
		args = append(args, "from", strings.TrimSpace(r.From))
	}
	args = append(args, r.To)

	return ufw.FormatCommand(action, args)
}

