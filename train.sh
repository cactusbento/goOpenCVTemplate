#!/bin/sh
/opt/opencv3/bin/opencv_traincascade -data cascade/ -vec vec/pos.vec -bg annotations/neg.txt -w 24 -h 24 \
	-numPos 25 -numNeg 125 -numStages 9 \
	-maxFalseAlarmRate 0.175 \
	-minHitRate 0.98 \
	-acceptanceRatioBreakValue 0.000001
