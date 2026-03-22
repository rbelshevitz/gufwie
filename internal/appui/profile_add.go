package appui

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/rbelshevitz/gufwie/internal/ufw"
)

const ufwApplicationsDir = "/etc/ufw/applications.d"

func (u *ui) addFromProfile() {
	u.runAsync("Loading profiles…", func(ctx context.Context) error {
		apps, err := u.client.AppList(ctx)
		if err != nil {
			return err
		}
		u.queueDraw(func() {
			if len(apps) == 0 {
				u.showModal(
					"No profiles",
					"No ufw application profiles found.\n\nProfiles usually come with packages like openssh-server, nginx, apache2.\nThey live in:\n  "+ufwApplicationsDir,
					[]string{"OK", "Profiles dir", "Freeform"},
					func(btn string) {
						switch btn {
						case "Profiles dir":
							u.showProfilesDir()
						case "Freeform":
							u.showAddMenu()
						}
					},
				)
				return
			}
			u.showProfilePicker(apps)
		})
		return nil
	})
}

func (u *ui) showProfilePicker(apps []string) {
	u.setModalHelp(func() {
		u.showDialogHelp("Help", profileHelpText())
	})
	action := ufw.ActionAllow
	query := ""
	selected := ""
	focusSlot := 1 // 0 search, 1 list, 2 actions, 3 info, 4 buttons
	buttonIdx := 0

	filtered := func() []string {
		q := strings.ToLower(strings.TrimSpace(query))
		if q == "" {
			return apps
		}
		out := make([]string, 0, len(apps))
		for _, a := range apps {
			if strings.Contains(strings.ToLower(a), q) {
				out = append(out, a)
			}
		}
		return out
	}

	list := tview.NewList()
	list.SetBorder(true).SetTitle("Profiles")
	list.ShowSecondaryText(false)
	list.SetMainTextColor(colorFG)
	list.SetSelectedTextColor(colorBG)
	list.SetSelectedBackgroundColor(colorAccent)
	list.SetBackgroundColor(colorBG)
	list.SetHighlightFullLine(true)

	info := tview.NewTextView().SetDynamicColors(true).SetWrap(true)
	info.SetBorder(true).SetTitle("Info")
	info.SetText(profileHint(apps))
	info.SetBackgroundColor(colorBG)
	info.SetTextColor(colorFG)

	search := tview.NewInputField().SetLabel("Search: ")
	search.SetFieldBackgroundColor(colorContrast2)
	search.SetFieldTextColor(colorFG)
	search.SetLabelColor(colorAccent)

	actions := tview.NewForm()
	actions.SetBorder(true).SetTitle("Action")
	applyFormTheme(actions)
	actions.AddDropDown("Action", []string{
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
	if dd, ok := actions.GetFormItem(0).(*tview.DropDown); ok {
		applyDropDownTheme(dd)
	}

	preview := tview.NewTextView().SetDynamicColors(true).SetWrap(false)
	preview.SetBorder(true).SetTitle("Preview")
	preview.SetBackgroundColor(colorBG)
	preview.SetTextColor(colorFG)

	setPreview := func() {
		if selected == "" {
			preview.SetText("Select a profile…")
			return
		}
		preview.SetText(ufw.FormatCommand(action, []string{selected}))
	}

	setProfilesTitle := func(visible int) {
		list.SetTitle(fmt.Sprintf("Profiles (%d/%d)", visible, len(apps)))
	}

	rebuildList := func() {
		items := filtered()
		list.Clear()
		for _, a := range items {
			list.AddItem(a, "", 0, func() {})
			// selection is handled via ChangedFunc
		}
		setProfilesTitle(len(items))
		if len(items) > 0 {
			list.SetCurrentItem(0)
			selected = items[0]
			setPreview()
		} else {
			selected = ""
			preview.SetText("No matches.")
		}
	}

	list.SetChangedFunc(func(_ int, mainText string, _ string, _ rune) {
		selected = mainText
		setPreview()
	})

	search.SetChangedFunc(func(text string) {
		query = text
		rebuildList()
	})

	btns := tview.NewForm()
	btns.SetButtonsAlign(tview.AlignRight)
	applyFormTheme(btns)

	refreshProfiles := func() {
		u.runAsync("Refreshing profiles…", func(ctx context.Context) error {
			newApps, err := u.client.AppList(ctx)
			if err != nil {
				return err
			}
			u.queueDraw(func() {
				apps = newApps
				info.SetText(profileHint(apps))
				rebuildList()
				setPreview()
			})
			return nil
		})
	}

	btns.AddButton("Close", func() {
		u.pages.RemovePage("modal")
		u.clearModalHelp()
		u.app.SetFocus(u.table)
	})
	btns.AddButton("Refresh", refreshProfiles)
	btns.AddButton("Profiles dir", func() {
		u.showProfilesDir()
	})
	btns.AddButton("Info", func() {
		if selected == "" {
			return
		}
		u.runAsync("Loading info…", func(ctx context.Context) error {
			out, err := u.client.AppInfo(ctx, selected)
			if err != nil {
				return err
			}
			u.queueDraw(func() {
				info.SetText(out)
			})
			return nil
		})
	})
	btns.AddButton("Apply", func() {
		if selected == "" {
			return
		}
		cmd := ufw.FormatCommand(action, []string{selected})
		u.showConfirm("Confirm", "Apply?\n\n"+cmd, "Apply", func() {
			u.runAsync("Applying…", func(ctx context.Context) error {
				return u.client.ApplyApp(ctx, action, selected)
			}, func(err error) {
				if err != nil {
					u.showError(err)
					return
				}
				u.pages.RemovePage("modal")
				u.refresh()
			})
		})
	})

	left := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(search, 1, 0, false).
		AddItem(list, 0, 1, true)
	right := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(actions, 3, 0, false).
		AddItem(preview, 3, 0, false).
		AddItem(info, 0, 1, false).
		AddItem(btns, 3, 0, false)

	root := tview.NewFlex().
		AddItem(left, 0, 2, true).
		AddItem(right, 0, 3, false)

	root.SetBorder(true).SetTitle(windowTitle("Add via Service/Profile"))

	focus := func() {
		switch focusSlot {
		case 0:
			u.app.SetFocus(search)
		case 1:
			u.app.SetFocus(list)
		case 2:
			u.app.SetFocus(actions)
			actions.SetFocus(0)
		case 3:
			u.app.SetFocus(info)
		case 4:
			u.app.SetFocus(btns)
			btns.SetFocus(buttonIdx)
		default:
			focusSlot = 1
			u.app.SetFocus(list)
		}
	}

	root.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyEscape:
			u.pages.RemovePage("modal")
			u.clearModalHelp()
			u.app.SetFocus(u.table)
			return nil
		case tcell.KeyTAB:
			if focusSlot == 4 && buttonIdx < btns.GetButtonCount()-1 {
				buttonIdx++
			} else {
				focusSlot = (focusSlot + 1) % 5
				if focusSlot == 4 {
					buttonIdx = 0
				}
			}
			focus()
			return nil
		case tcell.KeyBacktab:
			if focusSlot == 4 && buttonIdx > 0 {
				buttonIdx--
			} else {
				focusSlot = (focusSlot - 1 + 5) % 5
				if focusSlot == 4 {
					if n := btns.GetButtonCount(); n > 0 {
						buttonIdx = n - 1
					} else {
						buttonIdx = 0
					}
				}
			}
			focus()
			return nil
		case tcell.KeyF5:
			refreshProfiles()
			return nil
		}
		switch ev.Rune() {
		case '/':
			u.app.SetFocus(search)
			return nil
		case 'r':
			refreshProfiles()
			return nil
		}
		return ev
	})

	rebuildList()
	setPreview()

	u.pages.RemovePage("modal")
	u.pages.AddPage("modal", center(110, 28, root), true, true)
	focus()
}

func profileHint(apps []string) string {
	if len(apps) == 0 {
		return "No profiles found.\n\nProfiles are defined in:\n  " + ufwApplicationsDir + "\n\nInstall packages that ship profiles (e.g. openssh-server, nginx, apache2), then refresh and try again."
	}
	if len(apps) == 1 {
		return "Only one profile found.\n\nProfiles are defined in:\n  " + ufwApplicationsDir + "\n\nInstall packages that ship more profiles (e.g. nginx, apache2), then refresh."
	}
	return "Select a profile and press Info…\n\nProfiles are defined in:\n  " + ufwApplicationsDir
}

func (u *ui) showProfilesDir() {
	entries, err := os.ReadDir(ufwApplicationsDir)
	if err != nil {
		u.showError(err)
		return
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)
	if len(names) == 0 {
		u.showModal("Profiles dir", ufwApplicationsDir+"\n\n(no files)", []string{"OK"}, func(string) {})
		return
	}
	u.showModal("Profiles dir", ufwApplicationsDir+"\n\n"+strings.Join(names, "\n"), []string{"OK"}, func(string) {})
}
