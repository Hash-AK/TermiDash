package main

import (
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

//go:embed logos
var logoFiles embed.FS

const barWidth = 20

func createBar(percent float64) (string, string) {
	filledBlocks := int((percent / 100.0) * float64(barWidth))
	var colorCode string
	if percent >= 80 {
		colorCode = "[red]"
	} else if percent >= 50 {
		colorCode = "[yellow]"
	} else {
		colorCode = "[green]"
	}
	filledString := strings.Repeat("█", filledBlocks)
	emptyString := strings.Repeat("-", barWidth-filledBlocks)
	return colorCode + "[" + filledString + emptyString + "]" + "[-]", colorCode
}
func updateInfos(app *tview.Application, cpuPanel, memPanel, infoPanel, diskPanel *tview.TextView) {
	//General Info
	OSPlatform, _, _, _ := host.PlatformInformation()
	//debug OSPlatform = "windows10"
	logoBytes, err := logoFiles.ReadFile("logos/" + OSPlatform + ".ascii")
	var logo string
	if err == nil {
		translatedLogo := tview.TranslateANSI(string(logoBytes))
		logo = translatedLogo + "\n"
	} else {
		logo = ""
	}

	KernelVersion, _ := host.KernelVersion()
	cpuInfo, _ := cpu.Info()
	cpuModelName := cpuInfo[0].ModelName
	firstChar := strings.ToUpper(string(OSPlatform[0]))
	OSPlatform = firstChar + OSPlatform[1:]
	OSArch, _ := host.KernelArch()
	hostInfo, _ := host.Info()
	hostname := hostInfo.Hostname
	uptime := hostInfo.Uptime
	var uptimeInt int
	uptimeInt = int(uptime)
	uptimeString := time.Duration(uptimeInt) * time.Second
	OSInfoText := fmt.Sprintf("%sOS : %s %s\nKernel Version: %s\nHostname: %s\nUptime: %s\nCPU Model: %s\n", logo, OSPlatform, OSArch, KernelVersion, hostname, uptimeString, cpuModelName)

	//Memory
	v, _ := mem.VirtualMemory()
	totalMem := float64(v.Total)
	usedMem := float64(v.Used)
	usedMemPercent := v.UsedPercent
	var usedMemSign string
	var memSign string
	if (totalMem) >= float64(1000000000000) {
		totalMem = totalMem / float64(1000000000000)
		memSign = "TB"
	} else if (totalMem) >= float64(1000000000) {
		totalMem = totalMem / float64(1000000000)
		memSign = "GB"

	} else if (totalMem) >= float64(1000000) {
		totalMem = totalMem / float64(1000000)
		memSign = "MB"
	} else {
		totalMem = totalMem / float64(1000)
		memSign = "KB"
	}
	if (usedMem) >= float64(1000000000000) {
		usedMem = usedMem / float64(1000000000000)
		usedMemSign = "TB"

	} else if (usedMem) >= float64(1000000000) {
		usedMem = usedMem / float64(1000000000)
		usedMemSign = "GB"

	} else if (usedMem) >= float64(1000000) {
		usedMem = usedMem / float64(1000000)
		usedMemSign = "MB"

	} else {
		usedMem = usedMem / float64(1000)
		usedMemSign = "KB"
	}
	var usedMemPercentString string
	if usedMemPercent >= 80 {
		usedMemPercentString = fmt.Sprintf("[red]%.2f[-]", usedMemPercent)
	} else if usedMemPercent >= 50 {
		usedMemPercentString = fmt.Sprintf("[yellow]%.2f[-]", usedMemPercent)
	} else {
		usedMemPercentString = fmt.Sprintf("[green]%.2f[-]", usedMemPercent)

	}
	memUsageBar, memColCode := createBar(usedMemPercent)
	memBarString := fmt.Sprintf("Memory: %s", memUsageBar)
	memText := fmt.Sprintf("Total Memory: %.2f%s\nUsed Memory: %.2f%s (%s%%)\n%s%s[-]", totalMem, memSign, usedMem, usedMemSign, usedMemPercentString, memColCode, memBarString)

	//CPU

	cpuCountPhys, _ := cpu.Counts(false)
	cpuCountLogical, _ := cpu.Counts(true)
	globalCpuUse, _ := cpu.Percent(0, false)
	globalCpuUseFloat := globalCpuUse[0]
	var globalCpuUseString string
	if globalCpuUseFloat >= 80 {
		globalCpuUseString = fmt.Sprintf("[red]%.2f%%[-]", globalCpuUseFloat)

	} else if globalCpuUseFloat >= 50 {
		globalCpuUseString = fmt.Sprintf("[yellow]%.2f%%[-]", globalCpuUseFloat)
	} else {
		globalCpuUseString = fmt.Sprintf("[green]%.2f%%[-]", globalCpuUseFloat)
	}

	var barStrings string
	allCoresUsage, _ := cpu.Percent(0, true)
	for i := range allCoresUsage {
		currentCorePercentBar, colorCode := createBar(allCoresUsage[i])
		barStrings = fmt.Sprintf("%s%s\nCPU%d[-] %s %s%.0f%%[-]", barStrings, colorCode, i, currentCorePercentBar, colorCode, allCoresUsage[i])
	}
	cpuCountText := fmt.Sprintf("CPU count physical/logical: %v/%v\nTotal usage: %s%s", cpuCountPhys, cpuCountLogical, globalCpuUseString, barStrings)

	//Disk

	var diskUsageText string
	partitions, _ := disk.Partitions(false)
	var diskText string

	for i := range partitions {
		usage, _ := disk.Usage(partitions[i].Mountpoint)
		totalSpace := float64(usage.Total)
		usedSpace := float64(usage.Used)
		var totalSpaceSign string
		var usedSpaceSign string

		//TB or GB or MB or KB
		if totalSpace >= 1000000000000 {
			totalSpace = totalSpace / 1000000000000
			totalSpaceSign = "TB"
		} else if totalSpace >= 1000000000 {
			totalSpace = totalSpace / 1000000000
			totalSpaceSign = "GB"
		} else if totalSpace >= 1000000 {
			totalSpace = totalSpace / 1000000
			totalSpaceSign = "MB"
		} else if totalSpace >= 1000 {
			totalSpace = totalSpace / 1000
			totalSpaceSign = "KB"
		}

		//Again TB or GB or MB or KB
		if usedSpace >= 1000000000000 {
			usedSpace = usedSpace / 1000000000000
			usedSpaceSign = "TB"
		} else if usedSpace >= 1000000000 {
			usedSpace = usedSpace / 1000000000
			usedSpaceSign = "GB"
		} else if usedSpace >= 1000000 {
			usedSpace = usedSpace / 1000000
			usedSpaceSign = "MB"

		} else if usedSpace >= 1000 {
			usedSpace = usedSpace / 1000
			usedSpaceSign = "KB"
		}
		diskBar, _ := createBar(usage.UsedPercent)
		diskText = fmt.Sprintf("%s%s: %s %.2f%% Used(%.2f %s/%.2f %s)\n", diskText, usage.Path, diskBar, usage.UsedPercent, usedSpace, usedSpaceSign, totalSpace, totalSpaceSign)
	}
	diskUsageText = diskText
	//Update
	app.QueueUpdateDraw(func() {
		infoPanel.SetText(OSInfoText)
		memPanel.SetText(memText)
		cpuPanel.SetText(cpuCountText)
		diskPanel.SetText(diskUsageText)

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

	//Disk section

	diskPanel := tview.NewTextView()
	diskPanel.SetBorder(true)
	diskPanel.SetTitle("Disk Usage")
	diskPanel.SetDynamicColors(true)

	// General Layout
	rightColumnLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	rightColumnLayout.AddItem(cpuPanel, 0, 2, false)
	rightColumnLayout.AddItem(memPanel, 0, 1, false)
	app := tview.NewApplication()
	mainGrid := tview.NewGrid()
	mainGrid.SetRows(0, 0, 10)
	mainGrid.SetColumns(0, 0)
	mainGrid.SetBorder(true)
	mainGrid.AddItem(diskPanel, 2, 0, 1, 2, 0, 0, false)
	mainGrid.AddItem(infoPanel, 0, 0, 2, 1, 0, 0, false)
	mainGrid.AddItem(rightColumnLayout, 0, 1, 2, 1, 0, 0, false)
	app.SetRoot(mainGrid, true)
	go func() {
		updateInfos(app, cpuPanel, memPanel, infoPanel, diskPanel)

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			updateInfos(app, cpuPanel, memPanel, infoPanel, diskPanel)
		}
	}()
	app.Run()

}
