package appui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/rbelshevitz/gufwie/internal/ufw"
)

func (u *ui) applyFilterAndRender() {
	prevNumber := u.selectedRuleNumber()
	u.model.applyFilter()
	u.renderStatus()
	u.renderTable()
	if !u.restoreSelection(prevNumber) && len(u.model.indexMap) > 0 {
		u.table.Select(1, 0)
	}
	row, _ := u.table.GetSelection()
	u.renderDetailsForRow(row)
}

func (u *ui) renderStatus() {
	state := "[#d75f5f]inactive[-]"
	if u.model.status.Active {
		state = "[#00d75f]active[-]"
	}
	filter := ""
	if u.model.filter != "" {
		filter = fmt.Sprintf("  Search: [#00afaf]%s[-]", tview.Escape(u.model.filter))
	}
	mode := ""
	if u.model.dryRun {
		mode = "  [#d7af00]DRY-RUN[-]"
	}
	u.statusBar.SetText(fmt.Sprintf("[#00afaf]gufwie[-]  UFW: %s%s  Rules: %d  Showing: %d%s", state, mode, len(u.model.rules), len(u.model.indexMap), filter))
}

func (u *ui) renderTable() {
	u.table.Clear()
	headers := []string{"#", "To", "Action", "From", "v6"}
	for c, h := range headers {
		u.table.SetCell(0, c, tview.NewTableCell("[::b]"+h).
			SetSelectable(false).
			SetTextColor(colorBG).
			SetBackgroundColor(colorAccent))
	}
	q := strings.ToLower(strings.TrimSpace(u.model.filter))
	for i, ridx := range u.model.indexMap {
		r := u.model.rules[ridx]
		row := i + 1
		bg := colorBG
		if row%2 == 0 {
			bg = colorContrast
		}
		u.table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", r.Number)).
			SetTextColor(colorDim).
			SetBackgroundColor(bg))

		toCell := tview.NewTableCell(r.To).SetTextColor(colorFG).SetBackgroundColor(bg)
		fromCell := tview.NewTableCell(r.From).SetTextColor(colorFG).SetBackgroundColor(bg)
		actionCell := tview.NewTableCell(r.Action).SetTextColor(colorForAction(r.Action)).SetBackgroundColor(bg)

		if q != "" {
			if strings.Contains(strings.ToLower(r.To), q) {
				toCell.SetTextColor(colorAccent).SetAttributes(tcell.AttrBold)
			}
			if strings.Contains(strings.ToLower(r.From), q) {
				fromCell.SetTextColor(colorAccent).SetAttributes(tcell.AttrBold)
			}
			if strings.Contains(strings.ToLower(r.Action), q) {
				actionCell.SetAttributes(tcell.AttrBold)
			}
		}

		u.table.SetCell(row, 1, toCell)
		u.table.SetCell(row, 2, actionCell)
		u.table.SetCell(row, 3, fromCell)
		v6 := ""
		if r.V6 {
			v6 = "v6"
		}
		u.table.SetCell(row, 4, tview.NewTableCell(v6).SetTextColor(colorMagenta).SetBackgroundColor(bg))
	}
}

func (u *ui) renderDetailsForRow(row int) {
	r, ok := u.model.ruleForRow(row)
	if !ok {
		u.details.SetText("Select a rule…")
		return
	}
	u.details.SetText(formatDetails(r))
}

func formatDetails(r ufw.Rule) string {
	v6 := "no"
	if r.V6 {
		v6 = "yes"
	}
	cmd := suggestCommand(r)
	return fmt.Sprintf("[::b]Rule %d[-]\n\nTo: [#e6e6e6]%s[-]\nFrom: [#e6e6e6]%s[-]\nAction: [#00afaf]%s[-]\nv6: %s\n\nSuggested:\n[#00afaf]%s[-]\n\nCopy: [::b]y[-] cmd, [::b]Y[-] raw\n\nRaw:\n%s",
		r.Number,
		tview.Escape(r.To),
		tview.Escape(r.From),
		tview.Escape(r.Action),
		v6,
		tview.Escape(cmd),
		tview.Escape(r.Raw),
	)
}

func colorForAction(action string) tcell.Color {
	up := strings.ToUpper(action)
	switch {
	case strings.HasPrefix(up, "ALLOW"):
		return colorGood
	case strings.HasPrefix(up, "DENY"):
		return colorBad
	case strings.HasPrefix(up, "REJECT"):
		return colorMagenta
	case strings.HasPrefix(up, "LIMIT"):
		return colorWarn
	default:
		return colorAccent
	}
}
