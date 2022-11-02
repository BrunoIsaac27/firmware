package threads

import (
	"machine"
)

type Notify struct {
	LED    machine.Pin
	Signal bool
}

type Thread struct {
	Channel   <-chan Notify
	Goroutine int
}

func NewThread(Goroutine int, Channel chan Notify) Thread {
	thread := Thread{
		Channel,
		Goroutine,
	}
	return thread
}

func (t Thread) AllocThreads() {
	for i := 0; i < t.Goroutine; i++ {
		go t.worker()
	}
}

func (t Thread) worker() {
	for signal := range t.Channel {
		signal.LED.Set(signal.Signal)

	}
}
