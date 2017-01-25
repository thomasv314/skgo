package cli

import (
	ui "github.com/gizak/termui"
	"strconv"
)

type Window struct {
	Grid *ui.Grid

	evt       *SidekiqPollEvent
	host      string
	namespace string

	topBar *ui.Par
}

// Create the main window that gets rendered to UI
func NewWindow(host, namespace string) (window *Window) {
	window = &Window{
		host:      host,
		namespace: namespace,
		Grid:      ui.NewGrid(),
	}

	return window
}

func (w *Window) UpdateEvent(evt *SidekiqPollEvent) {
	w.evt = evt
	w.Draw()
}

func (w *Window) Setup() {
	w.topBar = ui.NewPar("skgo")

	w.Grid.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, w.topBar),
		),
	)
}

func (w *Window) Draw() {
	w.drawTopbar()
}

func (w *Window) drawTopbar() {
	var text = ""

	if w.evt != nil {
		numProcs := strconv.Itoa(len(w.evt.Processes))

		if numProcs == "0" {
			numProcs = "[" + numProcs + "](fg-red)"
		} else {
			numProcs = "[" + numProcs + "](fg-yellow)"
		}

		numJobs := strconv.Itoa(len(w.evt.Jobs()))

		text = "skgo - connected to [" + w.host + "](fg-green) on [" + w.namespace + "](fg-green)"
		text = text + " - processes: " + numProcs + " - jobs: " + numJobs
	} else {
		text = "skgo - not connected"
	}

	w.topBar.Text = text
	w.topBar.Width = 100
	w.topBar.Height = 1
	w.topBar.Border = false
}
