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

func formatBytes(value uint64) string {
	const base = 1024
	var returnString string
	if value < base {
		returnString = fmt.Sprintf("%d B", value)
	}
	if value >= 1099511627776 {
		returnValue := float64(value) / float64(1099511627776)
		returnString = fmt.Sprintf("%.2f TiB", returnValue)
	} else if value >= 1073741824 {
		returnValue := float64(value) / float64(1073741824)
		returnString = fmt.Sprintf("%.2f GiB", returnValue)
	} else if value >= 1048576 {
		returnValue := float64(value) / float64(1048576)
		returnString = fmt.Sprintf("%.2f MiB", returnValue)
	} else if value >= base {
		returnValue := float64(value) / float64(base)
		returnString = fmt.Sprintf("%.2f KiB", returnValue)

	}
	return returnString
}
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
	OSPlatform, OSFamily, OSVersion, _ := host.PlatformInformation()
	logoToSearch := OSPlatform
	if strings.Contains(logoToSearch, "Microsoft Windows 10") {
		logoToSearch = "windows10"
	} else if strings.Contains(logoToSearch, "Microsoft Windows 11") {
		logoToSearch = "windows11"
	} else if strings.Contains(logoToSearch, "macOS") || OSFamily == "Darwin" {
		logoToSearch = "macos"
	} else if strings.Contains(OSVersion, "kali") {
		logoToSearch = "kali"
	}
	//OSPlatform = "fedora"
	logoBytes, err := logoFiles.ReadFile("logos/" + logoToSearch + ".ascii")
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
	OSInfoText := fmt.Sprintf("%sOS: %s %s\nOS family: %s\nOS version: %s\nKernel Version: %s\nHostname: %s\nUptime: %s\nCPU Model: %s\n", logo, OSPlatform, OSArch, OSFamily, OSVersion, KernelVersion, hostname, uptimeString, cpuModelName)

	//Memory
	v, _ := mem.VirtualMemory()
	usedMemPercent := v.UsedPercent
	totalMemString := formatBytes(v.Total)
	usedMemString := formatBytes(v.Used)

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
	memText := fmt.Sprintf("Total Memory: %s\nUsed Memory: %s (%s%%)\n%s%s[-]", totalMemString, usedMemString, usedMemPercentString, memColCode, memBarString)

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
		totalSpaceString := formatBytes(usage.Total)
		usedSpaceString := formatBytes(usage.Used)

		diskBar, _ := createBar(usage.UsedPercent)
		diskText = fmt.Sprintf("%s%s: %s %.2f%% Used(%s/%s)\n", diskText, usage.Path, diskBar, usage.UsedPercent, usedSpaceString, totalSpaceString)
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
	infoPanel.SetBorderColor(tcell.ColorBlueViolet)
	//Memory section
	memPanel := tview.NewTextView()
	memPanel.SetBorder(true)
	memPanel.SetTitle("Memory")
	memPanel.SetBorderColor(tcell.ColorBlue)
	memPanel.SetDynamicColors(true)

	//Disk section

	diskPanel := tview.NewTextView()
	diskPanel.SetBorder(true)
	diskPanel.SetTitle("Disk Usage")
	diskPanel.SetDynamicColors(true)
	diskPanel.SetBorderColor(tcell.ColorOrange)
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
