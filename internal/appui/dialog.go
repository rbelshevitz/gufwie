package appui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// showDialog renders a tabbable dialog: body + buttons.
// defaultButton is 0-based index into buttons slice.
func (u *ui) showDialog(title, body string, buttons []string, defaultButton int, onDone func(btn string)) {
	u.showDialogPage("dialog", title, body, buttons, defaultButton, onDone)
}

func (u *ui) showDialogHelp(title, body string) {
	u.showDialogPage("dialog_help", title, body, []string{"OK"}, 0, func(string) {})
}

func (u *ui) showDialogPage(pageName, title, body string, buttons []string, defaultButton int, onDone func(btn string)) {
	if defaultButton < 0 {
		defaultButton = 0
	}
	if defaultButton >= len(buttons) {
		defaultButton = 0
	}

	text := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetText(strings.TrimSpace(body))
	text.SetBorder(true).SetTitle(windowTitle(title)).SetTitleAlign(tview.AlignLeft)

	btns := tview.NewForm()
	btns.SetBorder(false)
	applyFormTheme(btns)

	for _, b := range buttons {
		b := b
		btns.AddButton(b, func() {
			u.pages.RemovePage(pageName)
			if pageName == "dialog" {
				u.clearDialogHelp()
			}
			onDone(b)
		})
	}

	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(text, 0, 1, false).
		AddItem(btns, 3, 0, true)
	root.SetBorder(true)

	root.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyEscape:
			u.pages.RemovePage(pageName)
			if pageName == "dialog" {
				u.clearDialogHelp()
			}
			if len(buttons) > 0 {
				onDone(buttons[0])
			} else {
				onDone("")
			}
			return nil
		}
		return ev
	})

	btns.SetFocus(defaultButton)

	u.pages.RemovePage(pageName)
	u.pages.AddPage(pageName, center(90, 20, root), true, true)
	u.app.SetFocus(btns)
}

func escCloses(u *ui) func(*tcell.EventKey) *tcell.EventKey {
	return func(ev *tcell.EventKey) *tcell.EventKey {
		if ev.Key() == tcell.KeyEscape {
			u.pages.RemovePage("modal")
			u.clearModalHelp()
			u.app.SetFocus(u.table)
			return nil
		}
		return ev
	}
}
