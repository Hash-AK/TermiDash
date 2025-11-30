# TermiDash : the TERMInal DASHboard

Golang terminal dashboard that shows the current computer's specs, the available and used ressources and more. Heavily inspired by [btop](https://github.com/aristocratos/btop).

I made this project for [hackclub's Siege](https://siege.hackclub.com). It follows the Week's theme, 'Winter', because I added a 'Snow Day' theme (caution, it's really blinding...), a Nord theme, and the bar's characteres are now snowflakes by default. It also follows the 8th Week's framework theme because it uses two Golang  _frameworks_ to help display TUIs and get computer usage informations respectively, [TView](https://github.com/rivo/tview) and [Gopsutils](https://github.com/shirou/gopsutil).

To quit press CTRL+C or 'q'.
To open the settings (basically just to change the current theme for now), press 's'.
While in settigns you can press either 's' again, or ESC to go back to the dashboard.

It is still in developement so there might be bugs/missing features that I'd like to implement.

The selected theme and the characteres for the bars are stored in _TOML_ in the following directory :

### Linux
```~/.config/TermiDash/config.toml``` or ```$XDG_CONFIG_HOME/TermiDash/config.toml```

### macOS
```~/Library/Application Support/TermiDash/config.toml```
(To be confirmed)

### Windows
``` %APPDATA%\Local\TermiDash\config.toml```



This program also displays your current distro's logo on the left panel in a neofetch/fastfetch style. Please do note that it's not 100% failproof, for example Zorin is detected as Debian. Please also note that if your terminal doesn't support correctly all the colors some text may appear weirdly/not appear at all.
## Screenshots/Demo  
### V1.0.0
![testing on arch](/assets/TermiDashOnArch.png)  
_Testing on Arch Linux, pelase note that this was a testing version so the OS Family and OS Version fields are missing_

![testing on windows 10](/assets/TermiDashOnWindows10.png)
_Testing on Windows 10 Home in a VM, please excuse the 100% CPU usage on the only core_

![testing on windows 11](/assets/TermiDashOnWin11.png)
_Testing on Windows 11_

### V2.0.0
![testing on LM XFCE](/assets/TermiDashV2-LM-XFCE.png)
_Testing the v2.0.0 on Linux Mint, with the Nord theme applied from teh settings_

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
- Fixing the cpu speed fetch
