package appui

import (
	"context"
)

func (u *ui) showVerboseStatus() {
	u.runAsync("Loading verbose status…", func(ctx context.Context) error {
		out, err := u.client.StatusVerbose(ctx)
		if err != nil {
			return err
		}
		u.app.QueueUpdateDraw(func() {
			u.showModal("ufw status verbose", out, []string{"OK"}, func(string) {})
		})
		return nil
	})
}

