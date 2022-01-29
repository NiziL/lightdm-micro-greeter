# lightdm-micro-greeter
A minimalist LightDM greeter written in Go, powered by [gotk3](github.com/gotk3/gotk3)

## Features
- multi-user login
- no user list 
- only starts the default session
- must have a file /etc/lightdm/wallpaper.png 
- shitty UX

## Installation
And the shitty UX experience starts right now ! 
Only manual installation is provided, and I doubt it will easily work on your machine.

Here's my main process on ArchLinux.
You'll need the C shared libraries `lightdm-gobject-1`, `glib-2.0` and `gobject` installed. It should be the case after installing LightDM, but I can't tell for sure.
```bash
git clone https://github.com/nizil/lightdm-micro-greeter
cd lightdm-miro-greeter
go build
sudo cp lightdm-micro-greeter /usr/bin
sudo cp lightdm-micro-greeter.desktop /usr/share/xgreeters
```
Then, ensure `/etc/lightdm/lightdm.conf` contains `greeter-session=lightdm-micro-greeter` and restart LightDM.
You also need a background image at `/etc/lightdm/wallpaper.png`.

If it doesn't work, you might have to change the `cgo` flags in the preambule of `main.go`.
Feel free to contact me through an issue if you want to try this greeter and need some help :)

If it works... Yay ! Don't forget to keep another greeter installed on your machine. You know, just in case...

## Backlog 
- shutdown, reboot and so
- better UI/UX
    - HiDPI handling
    - error message feedback
    - background auto scaling
- config file
    - username autofill
    - entry size and location
    - background file
- user list 
- sessions list
