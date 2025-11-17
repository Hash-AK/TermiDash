package main

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {
	v, _ := mem.VirtualMemory()
	text := fmt.Sprintf("%v, Free:%v", v.Total, v.Free)
	box := tview.NewBox().SetBorder(true).SetTitle(text)

	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}
