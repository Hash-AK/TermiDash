package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func updateInfos(app *tview.Application, cpuPanel, memPanel, infoPanel *tview.TextView) {
	KernelVersion, _ := host.KernelVersion()
	cpuInfo, _ := cpu.Info()

	cpuModelName := cpuInfo[0].ModelName

	OSPlatform, _, _, _ := host.PlatformInformation()
	firstChar := strings.ToUpper(string(OSPlatform[0]))
	OSPlatform = firstChar + OSPlatform[1:]
	OSArch, _ := host.KernelArch()
	hostInfo, _ := host.Info()
	hostname := hostInfo.Hostname
	uptime := hostInfo.Uptime
	var uptimeInt int
	uptimeInt = int(uptime)
	uptimeString := time.Duration(uptimeInt) * time.Second
	OSInfoText := fmt.Sprintf("OS : %s %s\nKernel Version: %s\nHostname: %s\nUptime: %s\nCPU Model: %s\n", OSPlatform, OSArch, KernelVersion, hostname, uptimeString, cpuModelName)

	v, _ := mem.VirtualMemory()
	totalMem := v.Total
	usedMem := v.Used
	usedMemPercent := v.UsedPercent
	var usedMemSign string
	var memSign string
	if (totalMem) >= 10000000000000 {
		totalMem = totalMem / 10000000000000
		memSign = "TB"
	} else if (totalMem) >= 1000000000 {
		totalMem = totalMem / 1000000000
		memSign = "GB"

	} else if (totalMem) >= 1000000 {
		totalMem = totalMem / 1000000
		memSign = "MB"
	}
	if (usedMem) >= 10000000000000 {
		usedMem = usedMem / 10000000000000
		usedMemSign = "TB"

	} else if (usedMem) >= 1000000000 {
		usedMem = usedMem / 1000000000
		usedMemSign = "GB"

	} else if (usedMem) >= 1000000 {
		usedMem = usedMem / 1000000
		usedMemSign = "MB"

	}
	var usedMemPercentString string
	if usedMemPercent >= 80 {
		usedMemPercentString = fmt.Sprintf("[red]%.2f[::-]", usedMemPercent)
	} else if usedMemPercent >= 50 {
		usedMemPercentString = fmt.Sprintf("[yellow]%.2f[::-]", usedMemPercent)
	} else {
		usedMemPercentString = fmt.Sprintf("%.2f", usedMemPercent)

	}
	memText := fmt.Sprintf("Total Memory: %v%s\nUsed Memory: %v%s (%s%%)\n", totalMem, memSign, usedMem, usedMemSign, usedMemPercentString)
	cpuCountPhys, _ := cpu.Counts(false)
	cpuCountLogical, _ := cpu.Counts(true)
	cpuFreq := cpuInfo[0].Mhz
	globalCpuUse, _ := cpu.Percent(1*time.Second, false)

	var cpuFreqSign string
	if cpuFreq >= 1000 {
		cpuFreq = cpuFreq / 1000
		cpuFreqSign = "GHz"
	} else {
		cpuFreqSign = "MHz"
	}
	cpuCountText := fmt.Sprintf("CPU count physical/logical: %v/%v\nMax frequency: %.2f%s\nTotal CPU usage: %.2f%%", cpuCountPhys, cpuCountLogical, cpuFreq, cpuFreqSign, globalCpuUse)
	app.QueueUpdateDraw(func() {
		infoPanel.SetText(OSInfoText)
		memPanel.SetText(memText)
		cpuPanel.SetText(cpuCountText)

	})
}
func main() {
	//CPU section

	//cpuManufacturer := cpuInfo[0].VendorID

	cpuPanel := tview.NewTextView()
	cpuPanel.SetBorder(true)
	cpuPanel.SetTitle("CPU")
	cpuPanel.SetBorderColor(tcell.ColorGreen)
	cpuPanel.SetDynamicColors(true)

	//FastFetch-style section

	infoPanel := tview.NewTextView()
	infoPanel.SetBorder(true)
	infoPanel.SetTitle("System Information")
	infoPanel.SetDynamicColors(true)
	//Memory section
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
	go func() {
		updateInfos(app, cpuPanel, memPanel, infoPanel)

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			updateInfos(app, cpuPanel, memPanel, infoPanel)
		}
	}()
	app.Run()

}
