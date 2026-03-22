package appui

import (
	"context"
	"log"
	"time"

	"github.com/rivo/tview"

	"github.com/rbelshevitz/gufwie/internal/ufw"
)

type Client interface {
	Status(ctx context.Context) (ufw.Status, error)
	StatusNumbered(ctx context.Context) (ufw.NumberedStatus, error)
	StatusVerbose(ctx context.Context) (string, error)
	AppList(ctx context.Context) ([]string, error)
	AppInfo(ctx context.Context, profile string) (string, error)
	ApplyApp(ctx context.Context, action ufw.Action, profile string) error
	Enable(ctx context.Context) error
	Disable(ctx context.Context) error
	Reload(ctx context.Context) error
	DeleteNumber(ctx context.Context, number int) error
	Apply(ctx context.Context, action ufw.Action, args []string) error
}

type ui struct {
	app    *tview.Application
	pages  *tview.Pages
	client Client
	logger *log.Logger

	statusBar *tview.TextView
	helpBar   *tview.TextView
	helpText  string
	table     *tview.Table
	details   *tview.TextView

	model model

	modalHelp  func()
	dialogHelp func()
}

func Run(client Client, logger *log.Logger) error {
	applyTheme()

	u := &ui{
		app:    tview.NewApplication(),
		pages:  tview.NewPages(),
		client: client,
		logger: logger,
		model:  model{filter: ""},
	}

	u.statusBar = tview.NewTextView().SetDynamicColors(true).SetWrap(false)
	u.statusBar.SetBorder(false)
	u.statusBar.SetBackgroundColor(colorContrast2)
	u.statusBar.SetTextColor(colorFG)
	u.statusBar.SetTextAlign(tview.AlignLeft)

	u.helpBar = tview.NewTextView().SetDynamicColors(true).SetWrap(false)
	u.helpBar.SetBorder(false)
	u.helpBar.SetBackgroundColor(colorContrast2)
	u.helpBar.SetTextColor(colorFG)
	u.helpBar.SetTextAlign(tview.AlignLeft)
	u.helpText = "[#000000:#00afaf]F1[-:-] Help  [#000000:#00afaf]F2[-:-] Add  [#000000:#00afaf]F3[-:-] Search  [#000000:#00afaf]F5[-:-] Refresh  [#000000:#00afaf]F6[-:-] Verbose  [#000000:#00afaf]F8[-:-] Delete  [#000000:#00afaf]F10[-:-] Quit"
	u.helpBar.SetText(u.helpText)

	u.table = tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)
	u.table.SetBorder(true).SetTitle("UFW rules")
	u.table.SetSelectedStyle(selectedTableStyle())
	u.table.SetSelectionChangedFunc(func(row, _ int) {
		u.renderDetailsForRow(row)
	})
	u.table.SetSelectedFunc(func(row, _ int) {
		r, ok := u.model.ruleForRow(row)
		if !ok {
			return
		}
		u.showModal("Rule", r.Raw, []string{"OK"}, func(string) {})
	})

	u.details = tview.NewTextView().SetDynamicColors(true).SetWrap(true)
	u.details.SetBorder(true).SetTitle("Details")
	u.details.SetText("Select a rule…")
	u.details.SetBackgroundColor(colorBG)

	content := tview.NewFlex().
		AddItem(u.table, 0, 3, true).
		AddItem(u.details, 0, 2, false)
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(u.statusBar, 1, 0, false).
		AddItem(content, 0, 1, true).
		AddItem(u.helpBar, 1, 0, false)

	u.pages.AddPage("main", mainFlex, true, true)

	u.app.SetInputCapture(u.captureKeys)

	u.app.SetRoot(u.pages, true).EnableMouse(true)
	go u.refresh()
	return u.app.Run()
}

func (u *ui) runAsync(label string, fn func(ctx context.Context) error, onDone ...func(err error)) {
	m := tview.NewModal().SetText(label).AddButtons([]string{})
	m.SetBorder(true).SetTitle("Working")
	u.queueDraw(func() {
		u.pages.RemovePage("working")
		u.pages.AddPage("working", center(70, 7, m), true, true)
	})

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err := fn(ctx)
		u.queueDraw(func() {
			u.pages.RemovePage("working")
		})
		if len(onDone) > 0 {
			u.queueDraw(func() { onDone[0](err) })
		} else if err != nil {
			u.queueDraw(func() {
				u.showError(err)
			})
		}
	}()
}
