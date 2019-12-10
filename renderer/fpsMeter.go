package renderer

import (
	"time"
)

type FPSMeter struct {
	start      time.Time
	last       time.Time
	frames     int
	sampleTime float64
	FpsCap     float64
	FrameTime  float64
	value      float64
}

func CreateFPSMeter(sampleTime float64) *FPSMeter {
	return &FPSMeter{start: time.Now(), last: time.Now(), frames: 0, sampleTime: sampleTime, FpsCap: 120}
}

func (fps *FPSMeter) UpdateFPSMeter() float64 {
	fps.frames = fps.frames + 1
	elapsed := time.Since(fps.start)
	if elapsed.Seconds() >= fps.sampleTime {
		fps.value = (float64)(fps.frames) / fps.sampleTime
		fps.start = time.Now()
		fps.frames = 0
	}

	fps.FrameTime = time.Since(fps.last).Seconds()
	sleepTime := (time.Duration)((1000.0 / fps.FpsCap) - (1000.0 * fps.FrameTime))
	if sleepTime > 0 {
		time.Sleep(sleepTime * time.Millisecond)
	}
	frameTime := time.Since(fps.last)
	fps.last = time.Now()
	return frameTime.Seconds()
}

func (fps *FPSMeter) Value() float64 {
	return fps.value
}
