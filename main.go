package main

import (
	"github.com/rivo/tview"
)

func main() {
	//v, _ := mem.VirtualMemory()

	infoPanel := tview.NewBox().SetBorder(true).SetTitle("System Info")
	cpuPanel := tview.NewBox().SetBorder(true).SetTitle("CPU")
	memPanel := tview.NewBox().SetBorder(true).SetTitle("Memory")

	rightColumnLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	rightColumnLayout.AddItem(cpuPanel, 0, 1, false)
	rightColumnLayout.AddItem(memPanel, 0, 2, false)
	app := tview.NewApplication()
	mainGrid := tview.NewGrid()
	mainGrid.SetRows(0, 0, 10)
	mainGrid.SetColumns(0, 0)
	mainGrid.SetBorder(true)

	mainGrid.AddItem(infoPanel, 0, 0, 2, 1, 0, 0, false)
	mainGrid.AddItem(rightColumnLayout, 0, 1, 2, 1, 0, 0, false)
	app.SetRoot(mainGrid, true)
	app.Run()
}
