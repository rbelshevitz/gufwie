package appui

import (
	"context"

	"github.com/rbelshevitz/gufwie/internal/ufw"
)

func (u *ui) applyQuickAdd() {
	u.showQuickAdd(func(action ufw.Action, argsLine string) {
		args, err := ufw.SplitArgs(argsLine)
		if err != nil {
			u.showError(err)
			return
		}
		cmd := ufw.FormatCommand(action, args)
		u.showConfirm("Confirm", "Apply?\n\n"+cmd, "Apply", func() {
			u.runAsync("Applying…", func(ctx context.Context) error {
				return u.client.Apply(ctx, action, args)
			}, func(err error) {
				if err != nil {
					u.showError(err)
					return
				}
				u.refresh()
			})
		})
	})
}

func (u *ui) applyWizardAdd() {
	u.showWizardAdd(func(action ufw.Action, args []string) {
		cmd := ufw.FormatCommand(action, args)
		u.showConfirm("Confirm", "Apply?\n\n"+cmd, "Apply", func() {
			u.runAsync("Applying…", func(ctx context.Context) error {
				return u.client.Apply(ctx, action, args)
			}, func(err error) {
				if err != nil {
					u.showError(err)
					return
				}
				u.refresh()
			})
		})
	})
}
