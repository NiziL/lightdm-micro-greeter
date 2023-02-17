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
- [ ] better handling of multihead setup
- [ ] hexcode in Wallpaper config to control background color
- [x] HiDPI handling
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

If `lightdm-micro-greeter` binary is accessible from your `PATH` (i.e. `GOBIN` is in it), you could use [the desktop file](https://github.com/NiziL/lightdm-micro-greeter/blob/main/data/lightdm-micro-greeter.desktop) in the data directory and set `greeter-session=lightdm-micro-greeter`.

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
    "DPI": 96,
    "Entry": {
        "WidthChars": 10,
    },
    "Label": {
        "Margin": 10,
        "Color": "#ffffff",
        "UsernameText": "username:",
        "PasswordText": "password:",
    },
    "Box": {
        ""
    }
}
```

| Parameters | Effect |
|------------|--------|
| `Username` | keep empty for multi-user mode, providing an username will switch to single-user. |
| `Wallpaper` | path to an image or a directory, `/etc/lightdm/lightdm-micro-greeter/` will be prepended. |
|`DPI`| dpi used. |
| `Entry.WidthChars` | entry width in chars. |
| `Label.Margin` | label margin in pixel. |
| `Label.Color` | label text color. |
| `Label.UsernameText` | label text when waiting for username. |
| `Label.PasswordText` | label text when waiting for password. |
| `Label.Color` | label text color. |
| `Box.OffsetTop` | box offset from top. |
| `Box.OffsetBottom` | box offset from bottom. |

### Tips & Tricks

- If `Wallpaper` is a directory, it must only contain images, as the greeter will randomly chose a file from this directory. 
- LightDM must have access to `Wallpaper` (using `/etc/lightdm` is pretty convenient).
- (`XLocationRatio`, `YLocationRatio`) define the location of the entry box center. `(0, 0)` is the top left corner, `(1, 1)` is the bottom right corner and `(0.5, 0.5)` the screen center.

## Dev notes

I had to rely on C macro `G_CALLBACK` to bind the LightDM server callbacks, which is unfortunately not accessible through the `import "C"` statement.  
To bind a go function to glib events, I've exported few go functions using `//export` statement and binded them with `G_CALLBACK` from C code `greeter_signal_connect.c`.  
With this architecture, it has been pretty difficult to avoid the use of global vars to carry information from the UI to these callbacks. If you have any idea how to make a cleaner code here, open a ticket, I'll be happy to discuss about it :)
