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
	"github.com/shirou/gopsutil/v4/sensors"
)

//go:embed logos
var logoFiles embed.FS

const barWidth = 20

type PanelStyle struct {
	BorderColor     tcell.Color
	TitleColor      tcell.Color
	TextColor       tcell.Color
	BackGroundColor tcell.Color
}
type Theme struct {
	CPUPanel  PanelStyle
	MemPanel  PanelStyle
	InfoPanel PanelStyle
	DiskPanel PanelStyle
	TempPanel PanelStyle

	BarRed          tcell.Color
	BarYellow       tcell.Color
	BarGreen        tcell.Color
	Backgroundcolor tcell.Color
}

var defaultTheme = Theme{
	CPUPanel: PanelStyle{
		BorderColor:     tcell.ColorGreen,
		TitleColor:      tcell.ColorGreen,
		TextColor:       tcell.ColorWhite,
		BackGroundColor: tcell.ColorBlack,
	},
	MemPanel: PanelStyle{
		BorderColor:     tcell.ColorBlue,
		TitleColor:      tcell.ColorBlue,
		TextColor:       tcell.ColorWhite,
		BackGroundColor: tcell.ColorBlack,
	},
	InfoPanel: PanelStyle{
		BorderColor:     tcell.ColorOrange,
		TitleColor:      tcell.ColorOrange,
		TextColor:       tcell.ColorWhite,
		BackGroundColor: tcell.ColorBlack,
	},
	TempPanel: PanelStyle{
		BorderColor:     tcell.ColorSteelBlue,
		TitleColor:      tcell.ColorSteelBlue,
		TextColor:       tcell.ColorWhite,
		BackGroundColor: tcell.ColorBlack,
	},
	DiskPanel: PanelStyle{
		BorderColor:     tcell.ColorPurple,
		TitleColor:      tcell.ColorPurple,
		TextColor:       tcell.ColorWhite,
		BackGroundColor: tcell.ColorBlack,
	},
	BarGreen:        tcell.ColorGreen,
	BarYellow:       tcell.ColorYellow,
	BarRed:          tcell.ColorRed,
	Backgroundcolor: tcell.ColorBlack,
}
var nordTheme = Theme{
	CPUPanel: PanelStyle{
		BorderColor:     tcell.GetColor("#3b4252"),
		TitleColor:      tcell.GetColor("#88c0d0"),
		TextColor:       tcell.GetColor("#eceff4"),
		BackGroundColor: tcell.GetColor("#2e3440"),
	},
	MemPanel: PanelStyle{
		BorderColor:     tcell.GetColor("#3b4252"),
		TitleColor:      tcell.GetColor("#81a1c1"),
		TextColor:       tcell.GetColor("#eceff4"),
		BackGroundColor: tcell.GetColor("#2e3440"),
	},
	InfoPanel: PanelStyle{
		BorderColor:     tcell.GetColor("#3b4252"),
		TitleColor:      tcell.GetColor("#b48ead"),
		TextColor:       tcell.GetColor("#D8DEE9"),
		BackGroundColor: tcell.GetColor("#2e3440"),
	},
	TempPanel: PanelStyle{
		BorderColor:     tcell.GetColor("#3b4252"),
		TitleColor:      tcell.GetColor("#5E81AC"),
		TextColor:       tcell.GetColor("#ECEFF4"),
		BackGroundColor: tcell.GetColor("#2e3440"),
	},
	DiskPanel: PanelStyle{
		BorderColor:     tcell.GetColor("#3b4252"),
		TitleColor:      tcell.GetColor("#8FBCBB"),
		TextColor:       tcell.GetColor("#eceff4"),
		BackGroundColor: tcell.GetColor("#2e3440"),
	},
	BarGreen:        tcell.GetColor("#a3be8c"),
	BarYellow:       tcell.GetColor("#ebcb8b"),
	BarRed:          tcell.GetColor("#bf616a"),
	Backgroundcolor: tcell.GetColor("#2E3440"),
}

func applyTheme(theme *Theme, cpuPanel, memPanel, infoPanel, tempPanel, diskPanel *tview.TextView) {
	cpuPanel.SetBorderColor(theme.CPUPanel.BorderColor)
	cpuPanel.SetTitleColor(theme.CPUPanel.TitleColor)
	cpuPanel.SetTextColor(theme.CPUPanel.TextColor)
	cpuPanel.SetBackgroundColor(theme.CPUPanel.BackGroundColor)

	memPanel.SetBorderColor(theme.MemPanel.BorderColor)
	memPanel.SetTitleColor(theme.MemPanel.TitleColor)
	memPanel.SetTextColor(theme.MemPanel.TextColor)
	memPanel.SetBackgroundColor(theme.MemPanel.BackGroundColor)

	infoPanel.SetBorderColor(theme.InfoPanel.BorderColor)
	infoPanel.SetTitleColor(theme.InfoPanel.TitleColor)
	infoPanel.SetTextColor(theme.InfoPanel.TextColor)
	infoPanel.SetBackgroundColor(theme.InfoPanel.BackGroundColor)

	tempPanel.SetBorderColor(theme.TempPanel.BorderColor)
	tempPanel.SetTitleColor(theme.TempPanel.TitleColor)
	tempPanel.SetTextColor(theme.TempPanel.TextColor)
	tempPanel.SetBackgroundColor(theme.TempPanel.BackGroundColor)

	diskPanel.SetBorderColor(theme.DiskPanel.BorderColor)
	diskPanel.SetTitleColor(theme.DiskPanel.TitleColor)
	diskPanel.SetTextColor(theme.DiskPanel.TextColor)
	diskPanel.SetBackgroundColor(theme.DiskPanel.BackGroundColor)
}
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
func createBar(theme *Theme, percent float64) (string, string) {
	filledBlocks := int((percent / 100.0) * float64(barWidth))
	var colorCode string
	if percent >= 80 {
		colorCode = fmt.Sprintf("[%s]", theme.BarRed.TrueColor().String())
	} else if percent >= 50 {
		colorCode = fmt.Sprintf("[%s]", theme.BarYellow.TrueColor().String())
	} else {
		colorCode = fmt.Sprintf("[%s]", theme.BarGreen.TrueColor().String())
	}
	filledString := strings.Repeat("█", filledBlocks)
	emptyString := strings.Repeat("-", barWidth-filledBlocks)
	return colorCode + "[" + filledString + emptyString + "]" + "[-]", colorCode
}
func updateInfos(app *tview.Application, cpuPanel, memPanel, infoPanel, diskPanel, tempPanel *tview.TextView, theme *Theme) {
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
		colorCode := fmt.Sprintf("[%s]", theme.BarRed.TrueColor().String())
		usedMemPercentString = fmt.Sprintf("%s%.2f[-]", colorCode, usedMemPercent)
	} else if usedMemPercent >= 50 {
		colorCode := fmt.Sprintf("[%s]", theme.BarYellow.TrueColor().String())
		usedMemPercentString = fmt.Sprintf("%s%.2f[-]", colorCode, usedMemPercent)
	} else {
		colorCode := fmt.Sprintf("[%s]", theme.BarGreen.TrueColor().String())
		usedMemPercentString = fmt.Sprintf("%s%.2f[-]", colorCode, usedMemPercent)

	}
	memUsageBar, memColCode := createBar(theme, usedMemPercent)
	memBarString := fmt.Sprintf("Memory: %s", memUsageBar)
	memText := fmt.Sprintf("Total Memory: %s\nUsed Memory: %s (%s%%)\n%s%s[-]", totalMemString, usedMemString, usedMemPercentString, memColCode, memBarString)

	//CPU

	cpuCountPhys, _ := cpu.Counts(false)
	cpuCountLogical, _ := cpu.Counts(true)
	globalCpuUse, _ := cpu.Percent(0, false)
	globalCpuUseFloat := globalCpuUse[0]
	var globalCpuUseString string
	if globalCpuUseFloat >= 80 {
		colorCode := fmt.Sprintf("[%s]", theme.BarRed.TrueColor().String())
		globalCpuUseString = fmt.Sprintf("%s%.2f%%[-]", colorCode, globalCpuUseFloat)

	} else if globalCpuUseFloat >= 50 {
		colorCode := fmt.Sprintf("[%s]", theme.BarYellow.TrueColor().String())
		globalCpuUseString = fmt.Sprintf("%s%.2f%%[-]", colorCode, globalCpuUseFloat)
	} else {
		colorCode := fmt.Sprintf("[%s]", theme.BarGreen.TrueColor().String())
		globalCpuUseString = fmt.Sprintf("%s%.2f%%[-]", colorCode, globalCpuUseFloat)
	}

	var barStrings string
	allCoresUsage, _ := cpu.Percent(0, true)
	for i := range allCoresUsage {
		currentCorePercentBar, colorCode := createBar(theme, allCoresUsage[i])
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

		diskBar, _ := createBar(theme, usage.UsedPercent)
		diskText = fmt.Sprintf("%s%s: %s %.2f%% Used(%s/%s)\n", diskText, usage.Path, diskBar, usage.UsedPercent, usedSpaceString, totalSpaceString)
	}
	diskUsageText = diskText
	//Temperature

	temperatures, _ := sensors.SensorsTemperatures()
	var cpuText string
	for i := range temperatures {
		if strings.Contains(temperatures[i].SensorKey, "coretemp") || strings.Contains(temperatures[i].SensorKey, "k10temp") {
			cpuTemp := temperatures[i].Temperature
			cpuText = cpuText + fmt.Sprintf("%s : %.2fC\n", temperatures[i].SensorKey, cpuTemp)
		}
	}
	//Update
	app.QueueUpdateDraw(func() {
		infoPanel.SetText(OSInfoText)
		memPanel.SetText(memText)
		cpuPanel.SetText(cpuCountText)
		diskPanel.SetText(diskUsageText)
		tempPanel.SetText(cpuText)

	})
}
func main() {
	//CPU section

	//cpuManufacturer := cpuInfo[0].VendorID
	cpuPanel := tview.NewTextView()
	cpuPanel.SetScrollable(true)
	cpuPanel.SetBorder(true)
	cpuPanel.SetTitle("CPU")
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

	// Temperature section
	tempPanel := tview.NewTextView()
	tempPanel.SetBorder(true)
	tempPanel.SetTitle("Temperatures")
	tempPanel.SetDynamicColors(true)
	applyTheme(&defaultTheme, cpuPanel, memPanel, infoPanel, tempPanel, diskPanel)

	// General Layout
	rightColumnLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	rightColumnLayout.AddItem(cpuPanel, 0, 1, false)
	rightColumnLayout.AddItem(memPanel, 0, 1, false)
	rightColumnLayout.AddItem(tempPanel, 0, 1, false)
	app := tview.NewApplication()
	mainGrid := tview.NewGrid()
	mainGrid.SetRows(0, 0, 10)
	mainGrid.SetColumns(0, 0)
	mainGrid.SetBorder(true)
	mainGrid.AddItem(diskPanel, 2, 0, 1, 2, 0, 0, false)
	mainGrid.AddItem(infoPanel, 0, 0, 2, 1, 0, 0, false)
	mainGrid.AddItem(rightColumnLayout, 0, 1, 2, 1, 0, 0, false)
	pages := tview.NewPages()

	pages.AddPage("dashboard", mainGrid, true, true)

	app.SetRoot(mainGrid, true)
	go func() {
		updateInfos(app, cpuPanel, memPanel, infoPanel, diskPanel, tempPanel, &defaultTheme)

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			updateInfos(app, cpuPanel, memPanel, infoPanel, diskPanel, tempPanel, &defaultTheme)
		}
	}()
	app.Run()

}
