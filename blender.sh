#!/bin/sh

# Blenderをインストール
sudo snap install blender --classic

# レンダリング
sudo blender --background -noaudio randomCube.blend --threads 0 -E CYCLES --render-output img/anim -s 0 -e 0 -a