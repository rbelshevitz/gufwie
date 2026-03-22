package appui

func (u *ui) showAddMenu() {
	u.setDialogHelp(func() {
		u.showDialogHelp("Help", addMenuHelpText())
	})
	u.showModal("Add", "Choose how to add a rule:", []string{"Cancel", "Freeform", "Wizard", "Service/Profile"}, func(btn string) {
		u.clearDialogHelp()
		switch btn {
		case "Freeform":
			u.applyQuickAdd()
		case "Wizard":
			u.applyWizardAdd()
		case "Service/Profile":
			u.addFromProfile()
		}
	})
}
