package appui

import (
	"strings"

	"github.com/rivo/tview"

	"github.com/rbelshevitz/gufwie/internal/ufw"
)

func (u *ui) showError(err error) {
	u.logger.Printf("error: %v\n", err)
	u.showModal("Error", err.Error(), []string{"OK"}, func(string) {})
}

func (u *ui) showHelp() {
	body := strings.Join([]string{
		"Function keys:",
		"  F1 help   F2 add   F3 search   F5 refresh   F6 verbose   F8 delete   F10 quit",
		"",
		"UFW:",
		"  Tested with ufw 0.36.2 (older versions may work).",
		"",
		"Keys:",
		"  r  refresh (re-run `ufw status numbered`)",
		"  /  search/filter",
		"  a  add (choose: freeform / wizard / service)",
		"  Service/Profile: add via `ufw app` profiles",
		"  A  add rule (wizard + preview)",
		"  x/del  delete selected rule (uses `ufw --force delete N`)",
		"  e  enable/disable ufw (toggle)",
		"  l  reload ufw",
		"  v  `ufw status verbose`",
		"  c  show a suggested ufw command for selected rule",
		"  y  copy suggested command to clipboard",
		"  Y  copy selected rule line (raw) to clipboard",
		"  Tab / Shift-Tab  cycle focus",
		"  Esc closes dialogs",
		"  q  quit (confirm)",
		"",
		"Safety:",
		"  All potentially dangerous actions ask for confirmation; default is No.",
		"",
		"Freeform add examples (args):",
		"  22/tcp",
		"  in 80/tcp",
		"  from 192.168.1.0/24 to any port 22 proto tcp",
	}, "\n")
	u.showModal("Help", body, []string{"OK"}, func(string) {})
}

func (u *ui) showModal(title, body string, buttons []string, onDone func(btn string)) {
	// Default dialog help unless caller set a more specific one.
	if u.dialogHelp == nil {
		u.setDialogHelp(func() {
			u.showDialogHelp("Help", dialogHelpText())
		})
	}
	u.showDialog(title, body, buttons, 0, onDone)
}

func (u *ui) showConfirm(title, body, yesLabel string, onYes func()) {
	// "No" must be first: it becomes default.
	if u.dialogHelp == nil {
		u.setDialogHelp(func() {
			u.showDialogHelp("Help", dialogHelpText())
		})
	}
	u.showDialog(title, body, []string{"No", yesLabel}, 0, func(btn string) {
		if btn == yesLabel {
			onYes()
		}
	})
}

func (u *ui) showQuickAdd(onApply func(action ufw.Action, argsLine string)) {
	u.setModalHelp(func() {
		u.showDialogHelp("Freeform help", freeformHelpText())
	})
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(windowTitle("Add (freeform)")).SetTitleAlign(tview.AlignLeft)
	applyFormTheme(form)

	action := ufw.ActionAllow
	argsLine := ""
	form.AddDropDown("Action", []string{
		string(ufw.ActionAllow),
		string(ufw.ActionDeny),
		string(ufw.ActionReject),
		string(ufw.ActionLimit),
	}, 0, func(_ string, idx int) {
		switch idx {
		case 0:
			action = ufw.ActionAllow
		case 1:
			action = ufw.ActionDeny
		case 2:
			action = ufw.ActionReject
		case 3:
			action = ufw.ActionLimit
		}
	})
	if dd, ok := form.GetFormItem(0).(*tview.DropDown); ok {
		applyDropDownTheme(dd)
	}
	form.AddInputField("Args", "", 0, nil, func(v string) { argsLine = v })
	form.AddButton("Help", func() {
		u.showDialogHelp("Freeform help", freeformHelpText())
	})
	form.AddButton("Cancel", func() { u.pages.RemovePage("modal") })
	form.AddButton("Next", func() {
		u.pages.RemovePage("modal")
		u.clearModalHelp()
		onApply(action, argsLine)
	})

	form.SetInputCapture(escCloses(u))

	u.pages.RemovePage("modal")
	u.pages.AddPage("modal", center(90, 14, form), true, true)
	u.app.SetFocus(form)
}

func (u *ui) showWizardAdd(onApply func(action ufw.Action, args []string)) {
	u.setModalHelp(func() {
		u.showDialogHelp("Wizard help", wizardHelpText())
	})
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(windowTitle("Add (wizard)")).SetTitleAlign(tview.AlignLeft)
	applyFormTheme(form)

	action := ufw.ActionAllow
	direction := "in"
	iface := ""
	from := ""
	to := ""
	port := ""
	proto := "tcp"
	comment := ""

	form.AddDropDown("Action", []string{
		string(ufw.ActionAllow),
		string(ufw.ActionDeny),
		string(ufw.ActionReject),
		string(ufw.ActionLimit),
	}, 0, func(_ string, idx int) {
		switch idx {
		case 0:
			action = ufw.ActionAllow
		case 1:
			action = ufw.ActionDeny
		case 2:
			action = ufw.ActionReject
		case 3:
			action = ufw.ActionLimit
		}
	})
	if dd, ok := form.GetFormItem(0).(*tview.DropDown); ok {
		applyDropDownTheme(dd)
	}
	form.AddDropDown("Direction", []string{"in", "out", ""}, 0, func(opt string, _ int) { direction = opt })
	if dd, ok := form.GetFormItem(1).(*tview.DropDown); ok {
		applyDropDownTheme(dd)
	}
	form.AddInputField("Interface (optional)", "", 0, nil, func(v string) { iface = v })
	form.AddInputField("From (optional)", "", 0, nil, func(v string) { from = v })
	form.AddInputField("To (optional)", "", 0, nil, func(v string) { to = v })
	form.AddInputField("Port/service", "", 0, nil, func(v string) { port = v })
	form.AddDropDown("Proto", []string{"tcp", "udp", ""}, 0, func(opt string, _ int) { proto = opt })
	if dd, ok := form.GetFormItem(5).(*tview.DropDown); ok {
		applyDropDownTheme(dd)
	}
	form.AddInputField("Comment (optional)", "", 0, nil, func(v string) { comment = v })

	form.AddButton("Help", func() {
		u.showDialogHelp("Wizard help", wizardHelpText())
	})
	form.AddButton("Cancel", func() { u.pages.RemovePage("modal") })
	form.AddButton("Next", func() {
		args := buildWizardArgs(direction, iface, from, to, port, proto, comment)
		u.pages.RemovePage("modal")
		u.clearModalHelp()
		onApply(action, args)
	})

	form.SetInputCapture(escCloses(u))

	u.pages.RemovePage("modal")
	u.pages.AddPage("modal", center(90, 22, form), true, true)
	u.app.SetFocus(form)
}
