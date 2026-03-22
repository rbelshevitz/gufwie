package appui

import (
	"context"
	"fmt"

	"github.com/rbelshevitz/gufwie/internal/ufw"
)

func (u *ui) refresh() {
	u.runAsync("Refreshing…", func(ctx context.Context) error {
		ns, err := u.client.StatusNumbered(ctx)
		if err != nil {
			return err
		}
		u.app.QueueUpdateDraw(func() {
			u.model.setStatus(ns)
			u.applyFilterAndRender()
		})
		return nil
	})
}

func (u *ui) reload() {
	u.runAsync("Reloading…", func(ctx context.Context) error {
		return u.client.Reload(ctx)
	}, func(err error) {
		if err != nil {
			u.showError(err)
			return
		}
		u.refresh()
	})
}

func (u *ui) toggleEnable() {
	target := "enable"
	if u.model.status.Active {
		target = "disable"
	}
	u.showConfirm("Confirm", fmt.Sprintf("Do you want to %s ufw?", target), "Yes", func() {
		label := "Enabling…"
		if u.model.status.Active {
			label = "Disabling…"
		}
		u.runAsync(label, func(ctx context.Context) error {
			if u.model.status.Active {
				return u.client.Disable(ctx)
			}
			return u.client.Enable(ctx)
		}, func(err error) {
			if err != nil {
				u.showError(err)
				return
			}
			u.refresh()
		})
	})
}

func (u *ui) deleteSelected() {
	r, ok := u.selectedRule()
	if !ok {
		return
	}
	cmd := ufw.FormatUfwArgs([]string{"--force", "delete", fmt.Sprintf("%d", r.Number)})
	u.showConfirm("Confirm", "Delete rule?\n\n"+r.Raw+"\n\n"+cmd, "Delete", func() {
		u.runAsync("Deleting…", func(ctx context.Context) error {
			return u.client.DeleteNumber(ctx, r.Number)
		}, func(err error) {
			if err != nil {
				u.showError(err)
				return
			}
			u.refresh()
		})
	})
}

func (u *ui) confirmQuit() {
	u.showConfirm("Quit", "Exit gufwie?", "Quit", func() {
		u.app.Stop()
	})
}
