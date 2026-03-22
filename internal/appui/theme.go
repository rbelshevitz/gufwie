package appui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// htop-ish palette.
var (
	colorBG        = tcell.NewRGBColor(0x0b, 0x0f, 0x14)
	colorFG        = tcell.NewRGBColor(0xe6, 0xe6, 0xe6)
	colorDim       = tcell.NewRGBColor(0x9a, 0xa4, 0xad)
	colorAccent    = tcell.NewRGBColor(0x00, 0xaf, 0xaf) // cyan
	colorGood      = tcell.NewRGBColor(0x00, 0xd7, 0x5f) // green
	colorWarn      = tcell.NewRGBColor(0xd7, 0xaf, 0x00) // yellow
	colorBad       = tcell.NewRGBColor(0xd7, 0x5f, 0x5f) // red-ish
	colorMagenta   = tcell.NewRGBColor(0xaf, 0x5f, 0xd7)
	colorContrast  = tcell.NewRGBColor(0x16, 0x22, 0x30)
	colorContrast2 = tcell.NewRGBColor(0x20, 0x2f, 0x40)
)

func applyTheme() {
	tview.Styles.PrimitiveBackgroundColor = colorBG
	tview.Styles.ContrastBackgroundColor = colorContrast
	tview.Styles.MoreContrastBackgroundColor = colorContrast2
	tview.Styles.BorderColor = colorAccent
	tview.Styles.TitleColor = colorAccent
	tview.Styles.GraphicsColor = colorAccent
	tview.Styles.PrimaryTextColor = colorFG
	tview.Styles.SecondaryTextColor = colorDim
	tview.Styles.TertiaryTextColor = colorDim
	tview.Styles.InverseTextColor = colorBG
}

func selectedTableStyle() tcell.Style {
	// htop-like selection: bright bar with dark text.
	return tcell.StyleDefault.Background(colorAccent).Foreground(colorBG).Bold(true)
}

func applyFormTheme(f *tview.Form) {
	f.SetFieldBackgroundColor(colorContrast)
	f.SetFieldTextColor(colorFG)
	f.SetLabelColor(colorAccent)
	// Keep buttons readable on terminals without truecolor: neutral button background,
	// focus highlight is handled by tview via contrast colors.
	f.SetButtonBackgroundColor(colorContrast2)
	f.SetButtonTextColor(colorFG)
}

func applyDropDownTheme(d *tview.DropDown) {
	d.SetFieldBackgroundColor(colorContrast)
	d.SetFieldTextColor(colorFG)
	d.SetLabelColor(colorAccent)
	normal, selected := dropDownListStyles()
	d.SetListStyles(normal, selected)
}

func dropDownListStyles() (normal tcell.Style, selected tcell.Style) {
	normal = tcell.StyleDefault.Background(colorBG).Foreground(colorFG)
	selected = tcell.StyleDefault.Background(colorAccent).Foreground(colorBG).Bold(true)
	return normal, selected
}
