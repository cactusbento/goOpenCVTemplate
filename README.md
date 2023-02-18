# goOpenCVTemplate 

A starter repository for training a cascade classifier on game objects using [gocv](https://gocv.io/) on X11.

Mainly created because I wanted to experiment on visual game object detection for game automation.

---------


Follow the [OpenCV Cascade Classifier training article](https://docs.opencv.org/4.x/dc/d88/tutorial_traincascade.html) when using this repository, most things are automated here.

Note:
 * Training parameters should be adjusted to match the desired output
 * Adjust OpenCV application paths to match your system's OpenCV installation
 * `uninput`, `xdotool`, and `X11` are the current requirements, so no Windows, Mac, or Wayland

Instructions:
 * Press `v` to take positive screenshots.
 * Press `b` to take negative screenshots.
 * `./annotate.sh` to launch OpenCV's annotate program and start annotating positives.
 * ''

