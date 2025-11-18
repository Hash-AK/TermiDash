package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v4/cpu"
)

func main() {
	//CPU section
	cpuCountPhys, _ := cpu.Counts(false)
	cpuCountLogical, _ := cpu.Counts(true)
	cpuCountText := fmt.Sprintf("CPU count physical/logical: %v/%v", cpuCountPhys, cpuCountLogical)
	cpuPanel := tview.NewTextView()
	cpuPanel.SetText(cpuCountText)
	cpuPanel.SetBorder(true)
	cpuPanel.SetTitle("CPU")
	cpuPanel.SetBorderColor(tcell.ColorGreen)
	cpuPanel.SetDynamicColors(true)

	//FastFetch-style section
	infoPanel := tview.NewTextView()
	infoPanel.SetBorder(true)
	infoPanel.SetTitle("System Information")
	infoPanel.SetDynamicColors(true)

	// Memory Section
	memPanel := tview.NewTextView()
	memPanel.SetBorder(true)
	memPanel.SetTitle("Memory")
	memPanel.SetDynamicColors(true)

	// General Layout
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
