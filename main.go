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
var currentTheme = &defaultTheme

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

	BarRed                tcell.Color
	BarYellow             tcell.Color
	BarGreen              tcell.Color
	Backgroundcolor       tcell.Color
	DropDownOptionStyle   tcell.Style
	DropDownSelectedStyle tcell.Style
}
type Config struct {
	BarFilledChar rune
	BarEmptyChar  rune
}
type StaticInfo struct {
	Logo          string
	OS            string
	OSFamily      string
	OSVersion     string
	KernelVersion string
	KernelArch    string
	Hostname      string
	CPUPhysCore   int
	CPULogCore    int
	CPUModel      string
}

var appConfig = Config{
	BarFilledChar: '❄',
	BarEmptyChar:  '-',
}
var defaultTheme = Theme{
	CPUPanel: PanelStyle{
		BorderColor:     tcell.ColorGreen,
		TitleColor:      tcell.ColorGreen,
		TextColor:       tcell.ColorWhite,
		BackGroundColor: tcell.GetColor("#000000"),
	},
	MemPanel: PanelStyle{
		BorderColor:     tcell.ColorBlue,
		TitleColor:      tcell.ColorBlue,
		TextColor:       tcell.ColorWhite,
		BackGroundColor: tcell.GetColor("#000000"),
	},
	InfoPanel: PanelStyle{
		BorderColor:     tcell.ColorOrange,
		TitleColor:      tcell.ColorOrange,
		TextColor:       tcell.ColorWhite,
		BackGroundColor: tcell.GetColor("#000000"),
	},
	TempPanel: PanelStyle{
		BorderColor:     tcell.ColorSteelBlue,
		TitleColor:      tcell.ColorSteelBlue,
		TextColor:       tcell.ColorWhite,
		BackGroundColor: tcell.GetColor("#000000"),
	},
	DiskPanel: PanelStyle{
		BorderColor:     tcell.ColorPurple,
		TitleColor:      tcell.ColorPurple,
		TextColor:       tcell.ColorWhite,
		BackGroundColor: tcell.GetColor("#000000"),
	},
	BarGreen:              tcell.ColorGreen,
	BarYellow:             tcell.ColorYellow,
	BarRed:                tcell.ColorRed,
	Backgroundcolor:       tcell.GetColor("#000000"),
	DropDownOptionStyle:   tcell.StyleDefault.Foreground(tcell.GetColor("#ffffff")).Background(tcell.GetColor("#000000")),
	DropDownSelectedStyle: tcell.StyleDefault.Foreground(tcell.GetColor("#000000")).Background(tcell.ColorLightGray),
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
	BarGreen:              tcell.GetColor("#a3be8c"),
	BarYellow:             tcell.GetColor("#ebcb8b"),
	BarRed:                tcell.GetColor("#bf616a"),
	Backgroundcolor:       tcell.GetColor("#2E3440"),
	DropDownOptionStyle:   tcell.StyleDefault.Foreground(tcell.GetColor("#eceff4")).Background(tcell.GetColor("#434c5e")),
	DropDownSelectedStyle: tcell.StyleDefault.Foreground(tcell.GetColor("#2e3440")).Background(tcell.GetColor("#88c0d0")),
}

var snowTheme = Theme{
	CPUPanel: PanelStyle{
		BorderColor:     tcell.GetColor("#d8dee9"),
		TitleColor:      tcell.GetColor("#5e81ac"),
		TextColor:       tcell.GetColor("#2e3440"),
		BackGroundColor: tcell.GetColor("#e5e9f0"),
	},
	MemPanel: PanelStyle{
		BorderColor:     tcell.GetColor("#d8dee9"),
		TitleColor:      tcell.GetColor("#5e81ac"),
		TextColor:       tcell.GetColor("#2e3440"),
		BackGroundColor: tcell.GetColor("#e5e9f0"),
	},
	InfoPanel: PanelStyle{
		BorderColor:     tcell.GetColor("#d8dee9"),
		TitleColor:      tcell.GetColor("#5e81ac"),
		TextColor:       tcell.GetColor("#2e3440"),
		BackGroundColor: tcell.GetColor("#e5e9f0"),
	},
	DiskPanel: PanelStyle{
		BorderColor:     tcell.GetColor("#d8dee9"),
		TitleColor:      tcell.GetColor("#5e81ac"),
		TextColor:       tcell.GetColor("#2e3440"),
		BackGroundColor: tcell.GetColor("#e5e9f0"),
	},
	TempPanel: PanelStyle{
		BorderColor:     tcell.GetColor("#d8dee9"),
		TitleColor:      tcell.GetColor("#5e81ac"),
		TextColor:       tcell.GetColor("#2e3440"),
		BackGroundColor: tcell.GetColor("#e5e9f0"),
	},
	BarRed:          tcell.GetColor("#bf616a"),
	BarYellow:       tcell.GetColor("#ebcb8b"),
	BarGreen:        tcell.GetColor("#a3be8c"),
	Backgroundcolor: tcell.GetColor("#eceff4"),

	DropDownOptionStyle:   tcell.StyleDefault.Foreground(tcell.GetColor("#2e3440")).Background(tcell.GetColor("#eceff4a4")),
	DropDownSelectedStyle: tcell.StyleDefault.Foreground(tcell.GetColor("#eceff4")).Background(tcell.GetColor("#5e81ac")),
}
var themesList []string

func applyTheme(theme *Theme, cpuPanel, memPanel, infoPanel, tempPanel, diskPanel *tview.TextView, grid *tview.Grid, themeSelector *tview.DropDown, settings *tview.Form) {
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
	tview.Styles.PrimitiveBackgroundColor = theme.Backgroundcolor
	grid.SetBackgroundColor(theme.Backgroundcolor)

	themeSelector.SetLabelColor(theme.InfoPanel.TitleColor)
	themeSelector.SetFieldTextColor(theme.InfoPanel.TextColor)
	themeSelector.SetFieldBackgroundColor(theme.InfoPanel.BackGroundColor)
	themeSelector.SetListStyles(theme.DropDownOptionStyle, theme.DropDownSelectedStyle)
	settings.SetBackgroundColor(theme.Backgroundcolor)

}
func formatBytes(value uint64) string {
	const base = 1024
	var returnString string

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

	} else if value < base {
		returnString = fmt.Sprintf("%d B", value)

	}
	return returnString
}
func createBar(theme *Theme, percent float64, filledChar, emptyChar rune) (string, string) {
	filledBlocks := int((percent / 100.0) * float64(barWidth))
	var colorCode string
	if percent >= 80 {
		colorCode = fmt.Sprintf("[%s]", theme.BarRed.TrueColor().String())
	} else if percent >= 50 {
		colorCode = fmt.Sprintf("[%s]", theme.BarYellow.TrueColor().String())
	} else {
		colorCode = fmt.Sprintf("[%s]", theme.BarGreen.TrueColor().String())
	}
	filledString := strings.Repeat(string(filledChar), filledBlocks)
	emptyString := strings.Repeat(string(emptyChar), barWidth-filledBlocks)
	return colorCode + "[" + filledString + emptyString + "]" + "[-]", colorCode
}
func updateInfos(app *tview.Application, cpuPanel, memPanel, infoPanel, diskPanel, tempPanel *tview.TextView, theme *Theme, staticInfo *StaticInfo) {
	//General Info
	OSPlatform := staticInfo.OS
	OSFamily := staticInfo.OSFamily
	OSVersion := staticInfo.OSVersion
	KernelVersion := staticInfo.KernelVersion
	cpuModelName := staticInfo.CPUModel

	OSArch := staticInfo.KernelArch
	hostInfo, _ := host.Info()
	hostname := staticInfo.Hostname
	uptime := hostInfo.Uptime
	var uptimeInt int
	uptimeInt = int(uptime)
	uptimeString := time.Duration(uptimeInt) * time.Second
	logo := staticInfo.Logo
	OSInfoText := fmt.Sprintf("%s❄ OS: %s %s\n❄ OS family: %s\n❄ OS version: %s\n❄ Kernel Version: %s\n❄ Hostname: %s\n❄ Uptime: %s\n❄ CPU Model: %s\n", logo, OSPlatform, OSArch, OSFamily, OSVersion, KernelVersion, hostname, uptimeString, cpuModelName)

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
	memUsageBar, memColCode := createBar(theme, usedMemPercent, appConfig.BarFilledChar, appConfig.BarEmptyChar)
	memBarString := fmt.Sprintf("Memory: %s", memUsageBar)
	memText := fmt.Sprintf("Total Memory: %s\nUsed Memory: %s (%s%%)\n%s%s[-]", totalMemString, usedMemString, usedMemPercentString, memColCode, memBarString)

	//CPU

	cpuCountPhys := staticInfo.CPUPhysCore
	cpuCountLogical := staticInfo.CPULogCore
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
		currentCorePercentBar, colorCode := createBar(theme, allCoresUsage[i], appConfig.BarFilledChar, appConfig.BarEmptyChar)
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

		diskBar, _ := createBar(theme, usage.UsedPercent, appConfig.BarFilledChar, appConfig.BarEmptyChar)
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
	staticPlatform, staticFam, staticVersion, _ := host.PlatformInformation()
	logoToSearch := staticPlatform
	if strings.Contains(logoToSearch, "Microsoft Windows 10") {
		logoToSearch = "windows10"
	} else if strings.Contains(logoToSearch, "Microsoft Windows 11") {
		logoToSearch = "windows11"
	} else if strings.Contains(logoToSearch, "macOS") || staticFam == "Darwin" {
		logoToSearch = "macos"
	} else if strings.Contains(staticVersion, "kali") {
		logoToSearch = "kali"
	}
	logoBytes, err := logoFiles.ReadFile("logos/" + logoToSearch + ".ascii")
	var logo string
	if err == nil {
		translatedLogo := tview.TranslateANSI(string(logoBytes))
		logo = translatedLogo + "\n"
	} else {
		logo = ""
	}
	cpuInfo, _ := cpu.Info()
	cpuPhys, _ := cpu.Counts(false)
	cpuLog, _ := cpu.Counts(true)
	cpuModelName := cpuInfo[0].ModelName
	hostInfo, _ := host.Info()

	firstChar := strings.ToUpper(string(staticPlatform[0]))
	staticPlatform = firstChar + staticPlatform[1:]
	kernelVersion, _ := host.KernelVersion()
	hostname := hostInfo.Hostname
	kernelArch, _ := host.KernelArch()
	staticInfo := StaticInfo{
		Logo:          logo,
		OS:            staticPlatform,
		OSFamily:      staticFam,
		OSVersion:     staticVersion,
		KernelVersion: kernelVersion,
		KernelArch:    kernelArch,
		Hostname:      hostname,
		CPUPhysCore:   cpuPhys,
		CPULogCore:    cpuLog,
		CPUModel:      cpuModelName,
	}
	themesList = append(themesList, "Default", "Nord", "Snow Day")
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
	settings := tview.NewForm()
	settings.SetBorder(true)
	settings.SetTitle("Settings - ESC to go back")
	themeSelector := tview.NewDropDown()
	themeSelector.SetLabel("Select a theme (hit Enter): ")
	themeSelector.SetOptions(themesList, nil)
	themeSelector.SetCurrentOption(0)
	settings.AddFormItem(themeSelector)

	applyTheme(currentTheme, cpuPanel, memPanel, infoPanel, tempPanel, diskPanel, mainGrid, themeSelector, settings)

	pages := tview.NewPages()
	pages.AddPage("settings", settings, true, false)

	pages.AddPage("dashboard", mainGrid, true, true)
	app.SetRoot(pages, true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
			return nil
		}
		curentPage, _ := pages.GetFrontPage()
		if event.Rune() == 's' {
			if curentPage == "dashboard" {
				pages.SwitchToPage("settings")
			} else {
				pages.SwitchToPage("dashboard")
			}
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			if curentPage == "settings" {
				pages.SwitchToPage("dashboard")
			}
		}
		return event
	})
	settings.AddButton("Save and close", func() {
		_, selection := themeSelector.GetCurrentOption()
		switch selection {
		case "Default":
			currentTheme = &defaultTheme
			applyTheme(currentTheme, cpuPanel, memPanel, infoPanel, tempPanel, diskPanel, mainGrid, themeSelector, settings)

		case "Nord":
			currentTheme = &nordTheme
			applyTheme(currentTheme, cpuPanel, memPanel, infoPanel, tempPanel, diskPanel, mainGrid, themeSelector, settings)
		case "Snow Day":
			currentTheme = &snowTheme
			applyTheme(currentTheme, cpuPanel, memPanel, infoPanel, tempPanel, diskPanel, mainGrid, themeSelector, settings)
		}
		pages.SwitchToPage("dashboard")

	})
	go func() {

		updateInfos(app, cpuPanel, memPanel, infoPanel, diskPanel, tempPanel, currentTheme, &staticInfo)

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			updateInfos(app, cpuPanel, memPanel, infoPanel, diskPanel, tempPanel, currentTheme, &staticInfo)
		}
	}()
	app.Run()

}
