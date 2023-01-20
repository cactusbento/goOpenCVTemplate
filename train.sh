#!/bin/sh
/opt/opencv3/bin/opencv_traincascade -data cascade/ -vec vec/pos.vec -bg annotations/neg.txt -w 32 -h 32 -numPos 24 -numNeg 96 -numStages 9 -maxFalseAlarmRate 0.25 -minHitRate 0.999
