package appui

// queueDraw schedules a UI update safely from any goroutine.
// Important: tview's QueueUpdateDraw can deadlock if called from the UI event loop;
// this helper always calls it from a separate goroutine.
func (u *ui) queueDraw(fn func()) {
	go u.app.QueueUpdateDraw(fn)
}

