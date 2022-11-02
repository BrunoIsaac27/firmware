package hardware

import (
	"machine"
)

const (
	PressedIn = iota
	PressedOut
)

type EventButton int

type Interface struct {
	callback  func(EventButton)
	button    machine.Pin
	toogle    bool
	toogleOld bool
	state     int
}

func NewInterface(button machine.Pin, callback func(EventButton)) Interface {
	return Interface{
		callback: callback,
		button:   button,
	}
}

/*
Checks the status of the pin and calls the set
callback function if the pin has any status changes.

	Note: this function must be called in a loop inside goroutine,
	and at the end of this call runtime.Goshed() must be invoked.

Example:

	go func() {
		defer wg.Done()
		for {
			i2.Process()
			runtime.Gosched()
		}
	}()
*/
func (i *Interface) Process() {
	//fmt.Println("starting...")

	i.toogle = i.button.Get()
	if i.toogle != i.toogleOld {

		if i.toogle == true {
			//i.state = !i.state
			//i.state = PressedOut
			i.callback(PressedOut)

		} else {
			i.callback(PressedIn)
		}

		i.toogleOld = i.toogle

	}
}
