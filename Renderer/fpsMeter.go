package renderer

import (
	"fmt"
	"time"
)

type FPSMeter struct {
    start time.Time
    last time.Time
    frames int
    sampleTime float64
    FpsCap float64
}

func CreateFPSMeter( sampleTime float64 ) *FPSMeter {
	return &FPSMeter{ start: time.Now(), frames: 0, sampleTime: sampleTime, FpsCap: 120 }
}

func (fps *FPSMeter) UpdateFPSMeter() {
	fps.frames = fps.frames + 1
    elapsed := time.Since(fps.start)
    if elapsed.Seconds() >= fps.sampleTime {
    	fpsCount := (float64)(fps.frames) / fps.sampleTime
    	fmt.Printf("fps: %f\n",fpsCount)
    	fps.start = time.Now();
    	fps.frames = 0
    }

    frameTime := time.Since(fps.last)
    sleepTime := (time.Duration)((1000.0/fps.FpsCap) - (1000.0*frameTime.Seconds()))
    if sleepTime > 0 {
    	time.Sleep(sleepTime * time.Millisecond)
    }
    fps.last = time.Now() 
}

