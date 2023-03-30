# lightdm-micro-greeter

A minimalist LightDM greeter written in Go, powered by [gotk3](https://github.com/gotk3/gotk3) and inspired by [lightdm-mini-greeter](https://github.com/prikhi/lightdm-mini-greeter).  

Many thanks to Matt Fischer for [his great blog post](http://www.mattfischer.com/blog/archives/5).

![screenshot](https://github.com/NiziL/lightdm-micro-greeter/blob/main/data/example.jpg)
*Example screenshot, unknown artist: please open an issue to credit him/her !*


## Features

- log in
- single-user or multi-user mode
- suited for HiDPI monitor
- customization
    - wallpaper 
        - autoscaling on your primary monitor
        - random selection from a directory
    - entry
        - width
        - colors (font, background and carret)
    - label
        - Choose text for username and password mode

## Installation

### Package manager

:rotating_light: Only manual installation is provided for now :rotating_light:  
Any help to package lightdm-micro-greeter for your favorite distribution is greatly appreciated ! 

### Manual 

#### Requirements 

You need the C shared libraries `lightdm-gobject-1`, `gobject-2.0` and `glib-2.0`. They might be shipped with LightDM, but I can't tell for sure. Depending on your distro, you might have to install some `-dev` or `-devel` packages.

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

If you have used `go install` and `GOBIN` is in your `PATH`, you could use [the desktop file](https://github.com/NiziL/lightdm-micro-greeter/blob/main/data/lightdm-micro-greeter.desktop) provided.

If you have used `go build`, you could just run the [install.sh](https://github.com/NiziL/lightdm-micro-greeter/blob/main/install.sh) script **as root**, which put the built binary inside `/usr/bin/`, setup the desktop file and a default configuration, and even modify your LightDM configuration file (while doing a backup).
```bash
sudo ./install.sh
```


## Configuration

All the configuration is handled within the `/etc/lightdm/lightdm-micro-greeter/config.json` file.
If the file does not exist, [this configuration](https://github.com/NiziL/lightdm-micro-greeter/blob/main/data/config.json) will be used.

| Parameters | Effect |
|------------|--------|
| `Username` | keep empty for multi-user mode, providing an username will switch to single-user. |
| `Wallpaper` | path to an image or a directory, `/etc/lightdm/lightdm-micro-greeter/` will be prepended. |
|`DPI`| dpi used. |
| `Entry.WidthChars` | entry width in chars. |
| `Entry.TextColor` | entry text color (hexcode or rgba). |
| `Entry.BackgroundColor` | entry background color (hexcode or rgba). |
| `Entry.CaretColor` | entry caret color (hexcode or rgba). |
| `Entry.TextAlignment` | entry text alignement, float between 0 (left) and 1 (right). |
| `Label.Margin` | label margin in pixel. |
| `Label.Color` | label text color (hexcode or rgba). |
| `Label.UsernameText` | label text when waiting for username. |
| `Label.PasswordText` | label text when waiting for password. |
| `Label.Color` | label text color (hexcode or rgba). |
| `Box.MarginTop` | box margin top, in pixel. |
| `Box.MarginBottom` | box margin bottom, in pixel. |
| `Box.MarginLeft` | box margin left, in pixel. |
| `Box.MarginRight` | box margin right, in pixel. |

If `Wallpaper` is a directory, it must only contain images, as the greeter will randomly chose a file from this directory. 


## Dev notes

I had to rely on C macro `G_CALLBACK` to bind the LightDM server callbacks, which is unfortunately not accessible through the `import "C"` statement.  
To bind a go function to glib events, I've exported few go functions using `//export` statement and binded them with `G_CALLBACK` from C code `greeter_signal_connect.c`.  
With this architecture, it has been pretty difficult to avoid the use of global vars to carry information from the UI to these callbacks. If you have any idea how to make a cleaner code here, open a ticket, I'll be more than happy to discuss about it :)

