#!/usr/bin/bash

CONFIG_DIR=/etc/lightdm/lightdm-micro-greeter

# install binary
cp lightdm-micro-greeter /usr/bin

# prepare config folder
if [[ ! -d $CONFIG_DIR ]]; then
    mkdir $CONFIG_DIR
    cp data/config.json $CONFIG_DIR
fi

# prepare desktop file
cp data/lightdm-micro-greeter.desktop /usr/share/xgreeters/

# modify greeter in lightdm.conf
cp /etc/lightdm/lightdm.conf /etc/lightdm/.lightdm_conf.backup
sed -i "s/^greeter-session=.*$/greeter-session=lightdm-micro-greeter/g" /etc/lightdm/lightdm.conf

