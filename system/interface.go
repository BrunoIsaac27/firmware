package system

import (
	"log"
	"perseus/hardware"
	"time"
)

const (
	//Call when pressing the button for n time
	TimePressedCall = iota
	//Called when pressing button
	ImmediateCall
	//Called when release button
	ImmediateCallDrop
)

type systemIO struct {
	callback func()
	Config
	timeIn    time.Time
	timeSince time.Duration
	SinceOut  time.Duration
}

type Config struct {
	CallMode uint8
	During   time.Duration
}

func CreateIO(callback func()) systemIO {
	return systemIO{
		callback: callback,
	}
}

func (system *systemIO) Configure(config Config) {
	system.Config = config
}

func (system *systemIO) Manager(state hardware.EventButton) {

	if state == hardware.PressedIn {
		system.timeIn = time.Now()

		if system.Config.CallMode == ImmediateCall {
			//call here
			system.callback()
		}

	} else if state == hardware.PressedOut {
		system.SinceOut = time.Since(system.timeIn)
		log.Printf("[System Perseus] - Button duration: %vs", system.SinceOut.Seconds())
		if system.Config.CallMode == TimePressedCall {
			if system.SinceOut.Seconds() >= system.Config.During.Seconds() {
				//call here
				system.callback()
			}
		}

		if system.Config.CallMode == ImmediateCallDrop {
			//call here
			system.callback()
		}
	}

}
