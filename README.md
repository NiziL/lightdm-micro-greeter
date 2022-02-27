# lightdm-micro-greeter

A minimalist LightDM greeter written in Go, powered by [gotk3](https://github.com/gotk3/gotk3) and inspired by [lightdm-mini-greeter](https://github.com/prikhi/lightdm-mini-greeter).  

Many thanks to Matt Fischer for [his great blog post](http://www.mattfischer.com/blog/archives/5).


## Features

- log in
- single-user or multi-user mode
- wallpaper customization
    - autoscaling on your primary monitor
    - random selection from a directory

#### Backlog

- [ ] shutdown, reboot and so
- [x] random wallpaper from a directory
- [ ] better wallpaper placement in case of different aspect ratio
- [ ] better handling of multihead setup
- [ ] hexcode in Wallpaper config to control background color
- [ ] HiDPI handling (e.g. entry auto scaling)
- [ ] error message feedback (wrong password, unknown user...)
- [ ] packaging


## Installation

### Package manager

:rotating_light: Only manual installation is provided for now :rotating_light:  
Any help to package lightdm-micro-greeter for your favorite distribution is greatly appreciated ! 

### Manual 

#### Requirements 

You need the C shared libraries lightdm-gobject-1, gobject-2.0 and glib-2.0. They might be shipped with LightDM, but I can't tell for sure. Depending on your distro, you might have to install some `-dev` or `-devel` packages.

Obviously, you also need [Go](https://go.dev/doc/install).

#### Get the binary

```bash
go install github.com/nizil/lightdm-micro-greeter@latest
```
or 
```bash
git clone github.com/nizil/lightdm-micro-greeter
cd lightdm-micro-greeter
go build
```

#### Setup the greeter

Now, you have to tell LightDM to use this greeter, and this is done in two simples steps:
- Create a [desktop entry](https://wiki.archlinux.org/title/desktop_entries) at `/usr/share/xgreeters` which execute `lightdm-micro-greeter`. 
- Change the LightDM config to use the newly created `.desktop`, it could be done through the `greeter-session` parameter of `/etc/lightdm/lightdm.conf`.

If `lightdm-micro-greeter` binary is accessible from your `PATH`, you could use [the desktop file](https://github.com/NiziL/lightdm-micro-greeter/blob/main/data/lightdm-micro-greeter.desktop) in the data directory and set `greeter-session=lightdm-micro-greeter`.

If you've got the binary using a `git clone`/`go build`, you could just run the following commands **as root**.
```bash
cp lightdm-micro-greeter /usr/bin
mkdir /etc/lightdm/lightdm-micro-greeter
cp lightdm-micro-greeter /usr/bin
cp data/lightdm-micro-greeter.desktop /usr/share/xgreeters/
cp data/config.json /etc/lightdm/lightdm-micro-greeter/
sed -i "s/^greeter-session=.*$/greeter-session=lightdm-micro-greeter/g /etc/lightdm/lightdm.conf"
```

## Configuration

All the configuration is handled within the `/etc/lightdm/lightdm-micro-greeter/config.json` file.
If the file does not exist, the following configuration will be used:
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
| `Wallpaper` | path to an image or a directory, `/etc/lightdm/lightdm-micro-greeter/` will be prepended. |
| `Entry.WidthChars` | entry width in chars. |
| `Entry.Margin` | margin between label and entry in pixel. |
| `Entry.XLocationRatio ` | control entry x position |
| `Entry.YLocationRatio` | control entry y position |

### Tips & Tricks

- If `Wallpaper` is a directory, it must only contain images, as the greeter will randomly chose a file from this directory. 
- LightDM must have access to `Wallpaper` (using `/etc/lightdm` is pretty convenient).
- (`XLocationRatio`, `YLocationRatio`) define the location of the entry box center. `(0, 0)` is the top left corner, `(1, 1)` is the bottom right corner and `(0.5, 0.5)` the screen center.