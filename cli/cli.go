package cli

import (
	ui "github.com/gizak/termui"
	"github.com/thomasv314/skgo/sidekiq"
	"log"
	//	"os"
	"strconv"
)

var (
	client sidekiq.Client
)

func Start(redis_addr, redis_namespace string, poll_secs int) {
	client = sidekiq.NewClient(redis_addr, redis_namespace)
	defer client.Close()

	err := ui.Init()
	if err != nil {
		log.Fatal(err)
	}

	hdr := newSkgoHeader(redis_addr, redis_namespace)
	procLs := newProcessList()
	jobTb := newJobTable()
	queueLs := newQueueList()

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, hdr, procLs, jobTb),
		),
	)

	drawEvent := func(evt SidekiqPollEvent) {
		procLs.Items = buildProcessesList(&evt)
		procLs.Height = len(evt.Info.Processes) + 2

		jobsTableData := buildJobsTable(&evt)
		jobTb.Rows = jobsTableData
		jobTb.Height = (len(jobsTableData) * 2) + 1

		queueLs.Items = evt.Info.Queues
		queueLs.Height = len(evt.Info.Queues) + 2

		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	}

	drawEvent(pollSidekiq())

	// calculate layout
	//pollTimer := "/timer/" + strconv.Itoa(poll_secs) + "s"
	ui.Handle("/timer/1s", func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		// t is a EvtTimer
		if t.Count%2 == 0 {
			event := pollSidekiq()
			drawEvent(event)
		}
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()

	defer ui.Close()
}

func newSkgoHeader(addr, ns string) (header *ui.Par) {
	header = ui.NewPar("skgo - [connected to " + addr + "](fg-green) - namespace: " + ns)
	header.Y = 1
	header.Border = false
	return
}

func newProcessList() (processList *ui.List) {
	processList = ui.NewList()
	processList.ItemFgColor = ui.ColorYellow
	processList.BorderLabel = "Processes"
	return
}

func newJobTable() (jobTable *ui.Table) {
	jobTable = ui.NewTable()
	jobTable.FgColor = ui.ColorWhite
	jobTable.BgColor = ui.ColorDefault
	return
}

func newQueueList() (queueList *ui.List) {
	queueList = ui.NewList()
	queueList.ItemFgColor = ui.ColorWhite
	queueList.BorderLabel = "Queues"
	queueList.Y = 9
	return
}

func buildProcessesList(evt *SidekiqPollEvent) (processList []string) {
	processList = make([]string, 0)

	i := 0

	for proc := range evt.Processes {
		jobN := len(evt.Processes[proc])
		concrN := proc.Info.Concurrency

		clrSuffix := "](fg-green)"

		if jobN == concrN {
			clrSuffix = "](fg-red)"
		} else if jobN == 0 {
			clrSuffix = "](fg-blue)"
		}

		idx := "[" + strconv.Itoa(i) + "]"
		threadStr := "[" + strconv.Itoa(jobN) + "/" + strconv.Itoa(concrN) + clrSuffix
		name := idx + " " + proc.Name + " " + threadStr

		processList = append(processList, name)

		i++
	}

	return
}

func buildJobsTable(evt *SidekiqPollEvent) (jobTable [][]string) {
	jobTable = make([][]string, 0)

	jobTable = append(jobTable, []string{"Id", "Class", "Queue", "Arguments"})

	for proc := range evt.Processes {
		jobs := evt.Processes[proc]

		for i := range jobs {
			job := jobs[i]
			jobTable = append(jobTable, []string{job.Payload.Id, job.Payload.Class, job.Payload.Queue, ""})
		}
	}

	return
}
