#!/bin/bash

DISPLAY=:0.0 XAUTHORITY=/home/pi/.Xauthority /usr/bin/feh --cache-size 256 --auto-zoom --quiet --preload --randomize --full-screen --reload 60 -Y --slideshow-delay 15.0 /home/pi/share/japan/new/ --recursive --caption-path /captions --fontpath /home/pi/share --font lobster.ttf/48 --draw-tinted
