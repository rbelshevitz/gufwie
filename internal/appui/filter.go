package appui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (u *ui) showFilterPrompt() {
	u.setModalHelp(func() {
		u.showDialogHelp("Help", searchHelpText())
	})
	input := tview.NewInputField().
		SetLabel("Search: ").
		SetText(u.model.filter)
	input.SetFieldBackgroundColor(colorContrast2)
	input.SetFieldTextColor(colorFG)
	input.SetLabelColor(colorAccent)

	update := func(text string) {
		u.model.filter = text
		u.applyFilterAndRender()
	}
	input.SetChangedFunc(update)

	form := tview.NewForm()
	form.SetBorder(true).SetTitle(windowTitle("Search / Filter")).SetTitleAlign(tview.AlignLeft)
	applyFormTheme(form)
	form.AddFormItem(input)
	form.AddButton("Close", func() {
		u.pages.RemovePage("modal")
		u.clearModalHelp()
		u.app.SetFocus(u.table)
	})
	form.AddButton("Clear", func() {
		input.SetText("")
		update("")
	})

	input.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
			if ev.Key() == tcell.KeyEnter {
				update(strings.TrimSpace(input.GetText()))
			}
			u.pages.RemovePage("modal")
			u.clearModalHelp()
			u.app.SetFocus(u.table)
			return nil
		case tcell.KeyCtrlU:
			input.SetText("")
			update("")
			return nil
		}
		return ev
	})

	form.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if ev.Key() == tcell.KeyEscape {
			u.pages.RemovePage("modal")
			u.app.SetFocus(u.table)
			return nil
		}
		return ev
	})

	u.pages.RemovePage("modal")
	u.pages.AddPage("modal", center(80, 9, form), true, true)
	u.app.SetFocus(input)
}
