package server

import (
	"log"
	"math"
	"sync/atomic"
	"time"
)

const channelBufferSize = 100

// Monitor struct holds collection of monitored variables
type Monitor struct {
	printCh chan bool

	sendTimeCh      chan time.Duration
	sendTimes       []time.Duration
	physicsTimeCh   chan time.Duration
	physicsTimes    []time.Duration
	droppedMessages uint64
}

func newMonitor() *Monitor {
	return &Monitor{
		printCh:         make(chan bool),
		sendTimeCh:      make(chan time.Duration, channelBufferSize),
		sendTimes:       newDurationSlice(),
		physicsTimeCh:   make(chan time.Duration, channelBufferSize),
		physicsTimes:    newDurationSlice(),
		droppedMessages: 0,
	}
}

func (m *Monitor) run() {
	m.runPrintTicker()
	go m.loop()
}

func (m *Monitor) addSendTime(duration time.Duration) {
	select {
	case m.sendTimeCh <- duration:
	}
}

func (m *Monitor) addPhysicsTime(duration time.Duration) {
	select {
	case m.physicsTimeCh <- duration:
	}
}

func (m *Monitor) addDroppedMessage() {
	atomic.AddUint64(&m.droppedMessages, 1)
}

func (m *Monitor) loop() {
	for {
		select {
		case <-m.printCh:
			m.print()
		case duration := <-m.sendTimeCh:
			m.sendTimes = append(m.sendTimes, duration)
		case duration := <-m.physicsTimeCh:
			m.physicsTimes = append(m.physicsTimes, duration)
		}
	}
}

func newDurationSlice() []time.Duration {
	return make([]time.Duration, 0, 100)
}

func (m *Monitor) print() {
	m.printStats(m.sendTimes, "sendTime")
	m.printStats(m.physicsTimes, "physicsTime")

	droppedMessages := atomic.SwapUint64(&m.droppedMessages, 0)
	log.Printf("dropped messages: %d", droppedMessages)

	m.sendTimes = newDurationSlice()
	m.physicsTimes = newDurationSlice()
}

func (m *Monitor) printStats(durations []time.Duration, name string) {
	if len(durations) == 0 {
		log.Printf("%s: no samples", name)
		return
	}

	count := len(durations)
	min, max, avg := minMaxAvg(durations)
	stdDev := stdDev(durations, avg)

	log.Printf("%s: avg: %s, min: %s, max: %s, std_dev: %s, count: %d",
		name, avg, min, max, stdDev, count)
}

func minMaxAvg(durations []time.Duration) (time.Duration, time.Duration,
	time.Duration) {
	var max, sum time.Duration
	min := durations[0]

	count := int64(len(durations))

	for d := range durations {
		duration := durations[d]
		sum += duration

		if duration < min {
			min = duration
		}

		if duration > max {
			max = duration
		}
	}

	avg := time.Duration(sum.Nanoseconds() / count)
	return min, max, avg
}

func stdDev(durations []time.Duration, avg time.Duration) time.Duration {
	var sqSum time.Duration
	for d := range durations {
		duration := durations[d]
		sqSum += (duration - avg) * (duration - avg)
	}

	count := int64(len(durations))

	return time.Duration(math.Sqrt(float64(sqSum.Nanoseconds() / count)))
}

func (m *Monitor) runPrintTicker() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for _ = range ticker.C {
			m.printCh <- true
		}
	}()
}
