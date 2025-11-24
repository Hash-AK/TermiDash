# TermiDash : the TERMInal DASHboard

Golang terminal dashboard that shows the current computer's specs, the available and used ressources and more. Heavily inspired by [btop](https://github.com/aristocratos/btop).

I made this project for [hackclub's Siege](https://siege.hackclub.com). It follow the framework theme because it uses two Golang  _frameworks_ to help display TUIs and get computer usage informations respectively, [TView](https://github.com/rivo/tview) and [Gopsutils](https://github.com/shirou/gopsutil).

It is still in developement so there might be bugs/missign features that I'd like to implement.

This program also displays your current distro's logo on the left panel in a neofetch/fastfetch style. Please do note that it's not 100% failproof, so from my testing sometimes Linux Mint is detected as Debian, or Windows versions may or may not display the correct logo.

## Instalation

To install, you can either download one of the prebuilt binary from the Release tab, or you can clone this repo locally :
```bash
git clone https://github.com/Hash-AK/TermiDash
cd TermiDash
go run .
```

## To Do:
- Adding a cpu scheduler parsing
- Adding battery informations
- Adding temperature information
- Fixing the cpu speed fetch