#!/bin/bash

# aptアップデート,snapインストール,Blenderインストール
sudo apt -y update
sudo apt list --upgradable
sudo apt install snapd
sudo snap install blender --classic

# Blenderをレンダリング
# https://docs.blender.org/manual/en/latest/advanced/command_line/arguments.html
sudo blender --background -noaudio blend/Miraikomachi.blend --threads 0 -E CYCLES --render-output /img //anim --render-frame 2000
#                                                                                                   |          