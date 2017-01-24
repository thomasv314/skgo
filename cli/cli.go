package cli

import (
	ui "github.com/gizak/termui"
	"github.com/thomasv314/skgo/sidekiq"
	"log"
)

var (
	client sidekiq.Client
	window *Window
)

func Start(redis_addr, redis_namespace string, poll_secs int) {
	client = sidekiq.NewClient(redis_addr, redis_namespace)
	defer client.Close()

	err := ui.Init()
	if err != nil {
		log.Fatal(err)
	}

	window = NewWindow(redis_addr, redis_namespace)
	window.Setup()
	window.Draw()

	ui.Render(window.Grid)

	// calculate layout
	//pollTimer := "/timer/" + strconv.Itoa(poll_secs) + "s"
	ui.Handle("/timer/1s", func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		// t is a EvtTimer
		if t.Count%2 == 0 {
			event := pollSidekiq()
			window.UpdateEvent(&event)
			ui.Render(window.Grid)
		}
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()

	defer ui.Close()
}
