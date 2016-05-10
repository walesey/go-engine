package countdown

import "time"

type Countdown struct {
	times        []int64
	lastCount    int64
	lastEstimate int64
	count        int64
	total        int64
}

func (c *Countdown) Start(length int) {
	c.total = int64(length)
	c.times = []int64{}
	c.lastEstimate = 0
	c.count = int64(length)
	c.lastCount = time.Now().UnixNano()
}

func (c Countdown) SecondsRemaining() int64 {
	return c.lastEstimate
}

func (c Countdown) PercentageComplete() int {
	return int((float32(c.total-c.count) / float32(c.total)) * 100)
}

func (c *Countdown) Count() {
	c.times = append(c.times, time.Now().UnixNano()-c.lastCount)
	c.lastCount = time.Now().UnixNano()
	c.count--
	total := int64(0)
	for _, time := range c.times {
		total = total + time
	}
	averageTime := total / int64(len(c.times))
	timeRemaining := averageTime * c.count
	timeRemainingSeconds := timeRemaining / 1000000000
	if timeRemainingSeconds != c.lastEstimate {
		c.lastEstimate = timeRemainingSeconds
	}

}
