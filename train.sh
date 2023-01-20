#!/bin/sh
/opt/opencv3/bin/opencv_traincascade -data cascade/ -vec vec/pos.vec -bg annotations/neg.txt -w 24 -h 24 \
	-numPos 16 -numNeg 160 -numStages 9 \
	-maxFalseAlarmRate 0.15 \
	-minHitRate 0.99 \
	-acceptanceRatioBreakValue 0.00001
