#!/bin/sh
/opt/opencv3/bin/opencv_createsamples -info annotations/pos.txt -w 24 -h 24 -num $(ls img/pos | wc -l) -vec vec/pos.vec
