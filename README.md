# lightdm-micro-greeter
A minimalist LightDM greeter written in Go, powered by [gotk3](https://github.com/gotk3/gotk3) and inspired by [lightdm-mini-greeter](https://github.com/prikhi/lightdm-mini-greeter).  

Many thanks to Matt Fischer for [his great blog post](http://www.mattfischer.com/blog/archives/5).

## Features
- log in (hopefully)
- single-user or multi-user mode
- wallpaper autoscaling on your primary monitor
- :bug:

## Installation

:rotating_light: Only manual installation is provided for now. :rotating_light:

You'll need the C shared libraries `lightdm-gobject-1`, `glib-2.0` and `gobject-2.0` installed. It should be the case after installing LightDM, but I can't tell for sure.  

Here's my main process on ArchLinux:
```bash
git clone https://github.com/nizil/lightdm-micro-greeter
cd lightdm-miro-greeter
go build
sudo cp lightdm-micro-greeter /usr/bin
sudo cp data/lightdm-micro-greeter.desktop /usr/share/xgreeters/
sudo mkdir /etc/lightdm/lightdm-micro-greeter
sudo cp data/config.json /etc/lightdm/lightdm-micro-greeter/
```

Alternatively, you could use `go install github.com/nizil/lightdm-micro-greeter@latest` to get the executable, but you'll have to put `GOBIN` in your path and still add the `.desktop` file to `/usr/share/xgreeters`.

Then, ensure LightDM is using this greeter (`greeter-session=lightdm-micro-greeter` in `/etc/lightdm/lightdm.conf`) and restart LightDM (`systemctl restart lightdm`).

Feel free to contact me through an issue if you want to try this greeter and need some help.  
If it works... Yay ! Don't forget to keep another greeter installed on your machine. You know, just in case ;)

Any help to create a packaging solution is greatly appreciated !

## Configuration

All the configuration is handled within the `/etc/lightdm/lightdm-micro-greeter/config.json` file.
```json
{
   "Username": "",
   "Wallpaper" : "",
   "Entry": {
       "WidthChars": 10,
       "Margin": 10,
       "XLocationRatio": 0.5,
       "YLocationRatio": 0.5
   }
}
```

| Parameters | Effect |
|------------|--------|
| `Username` | keep empty for multi-user mode, providing an username will switch to single-user. |
| `Wallpaper` | path to an image, `/etc/lightdm/lightdm-micro-greeter/` will be prepended. |
| `Entry.WidthChars` | entry width in chars. |
| `Entry.Margin` | margin between label and entry in pixel. |
| `Entry.XLocationRatio ` | control entry x position |
| `Entry.YLocationRatio` | control entry y position |

## Backlog 
- [ ] shutdown, reboot and so
- [ ] random wallpaper from a directory
- [ ] better wallpaper placement in case of different aspect ratio
- [ ] better handling of multihead setup
- [ ] hexcode in Wallpaper config to control background color
- [ ] HiDPI handling (entry auto scaling)
- [ ] error message feedback (wrong password, unknown user...)
- [ ] user list (?)
- [ ] sessions list (?)
