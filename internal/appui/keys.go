package appui

import "github.com/gdamore/tcell/v2"

func (u *ui) captureKeys(ev *tcell.EventKey) *tcell.EventKey {
	switch ev.Key() {
	case tcell.KeyEscape:
		if u.pages.HasPage("dialog_help") {
			u.pages.RemovePage("dialog_help")
			return nil
		}
		if u.pages.HasPage("dialog") {
			u.pages.RemovePage("dialog")
			u.clearDialogHelp()
			return nil
		}
		if u.pages.HasPage("modal") {
			u.pages.RemovePage("modal")
			u.clearModalHelp()
			return nil
		}
		if u.pages.HasPage("working") {
			u.pages.RemovePage("working")
			return nil
		}
		return nil
	case tcell.KeyTAB:
		// Simple focus cycle on main screen.
		if !u.pages.HasPage("modal") && !u.pages.HasPage("dialog") && !u.pages.HasPage("working") {
			if u.app.GetFocus() == u.table {
				u.app.SetFocus(u.details)
				u.flashHelp("Focus: Details (Shift-Tab to return)")
			} else {
				u.app.SetFocus(u.table)
				u.flashHelp("Focus: Rules (Tab for Details)")
			}
			return nil
		}
	case tcell.KeyBacktab:
		if !u.pages.HasPage("modal") && !u.pages.HasPage("dialog") && !u.pages.HasPage("working") {
			if u.app.GetFocus() == u.details {
				u.app.SetFocus(u.table)
				u.flashHelp("Focus: Rules (Tab for Details)")
			} else {
				u.app.SetFocus(u.details)
				u.flashHelp("Focus: Details (Shift-Tab to return)")
			}
			return nil
		}
	case tcell.KeyF1:
		if u.pages.HasPage("dialog") && u.dialogHelp != nil {
			u.dialogHelp()
		} else if u.pages.HasPage("modal") && u.modalHelp != nil {
			u.modalHelp()
		} else {
			u.showHelp()
		}
		return nil
	case tcell.KeyF2:
		u.showAddMenu()
		return nil
	case tcell.KeyF3:
		u.showFilterPrompt()
		return nil
	case tcell.KeyF5:
		u.refresh()
		return nil
	case tcell.KeyF6:
		u.showVerboseStatus()
		return nil
	case tcell.KeyF8:
		u.deleteSelected()
		return nil
	case tcell.KeyF10:
		u.confirmQuit()
		return nil
	case tcell.KeyDelete:
		u.deleteSelected()
		return nil
	}

	switch ev.Rune() {
	case 'q':
		u.confirmQuit()
		return nil
	case 'r':
		u.refresh()
		return nil
	case 'l':
		u.reload()
		return nil
	case 'e':
		u.toggleEnable()
		return nil
	case 'x', 'd':
		u.deleteSelected()
		return nil
	case 'a':
		u.showAddMenu()
		return nil
	case 'A':
		u.applyWizardAdd()
		return nil
	case 'v':
		u.showVerboseStatus()
		return nil
	case 'c':
		u.showSuggestedCommand()
		return nil
	case 'y':
		u.copySuggestedCommand()
		return nil
	case 'Y':
		u.copyRawRule()
		return nil
	case '?':
		if u.pages.HasPage("dialog") && u.dialogHelp != nil {
			u.dialogHelp()
		} else if u.pages.HasPage("modal") && u.modalHelp != nil {
			u.modalHelp()
		} else {
			u.showHelp()
		}
		return nil
	case '/':
		u.showFilterPrompt()
		return nil
	}
	return ev
}
