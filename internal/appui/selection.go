package appui

import (
	"context"

	"github.com/rbelshevitz/gufwie/internal/ufw"
)

func (u *ui) selectedRule() (ufw.Rule, bool) {
	row, _ := u.table.GetSelection()
	return u.model.ruleForRow(row)
}

func (u *ui) selectedRuleNumber() int {
	r, ok := u.selectedRule()
	if !ok {
		return 0
	}
	return r.Number
}

func (u *ui) restoreSelection(number int) bool {
	if number <= 0 {
		return false
	}
	for i, ridx := range u.model.indexMap {
		if u.model.rules[ridx].Number == number {
			u.table.Select(i+1, 0)
			return true
		}
	}
	return false
}

func (u *ui) copyRawRule() {
	r, ok := u.selectedRule()
	if !ok {
		return
	}
	if err := CopyText(context.Background(), r.Raw+"\n"); err != nil {
		u.showError(err)
		return
	}
	u.flashHelp("Copied rule line to clipboard")
}

