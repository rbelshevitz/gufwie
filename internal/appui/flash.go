package appui

import "time"

func (u *ui) flashHelp(msg string) {
	u.helpBar.SetText(msg)
	time.AfterFunc(1200*time.Millisecond, func() {
		u.app.QueueUpdateDraw(func() {
			u.helpBar.SetText(u.helpText)
		})
	})
}

