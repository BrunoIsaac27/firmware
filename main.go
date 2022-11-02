// Perseus is an efficient firmware written in Go
// with the aim of having high performance and
// productivity at the same time. Currently this
// version should not be considered stable for production,
// and does not have remote upgrade possibilities.

package main

import (
	"fmt"
	"log"
	"machine"
	"perseus/hardware"
	"perseus/system"
	"runtime"
	"sync"
	"time"
)

var active bool

func emergency() {
	log.Fatal(`[System Perseus] - Emergency stop triggered:
	- ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ -  
	 Sorry! but something did not
	go as expected. There was some instability in the OS, maybe you can solve this just
	by restarting ;)

	If the problem persists, please go to the critical call center: https://critical.perseus.com
	- ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - ⚠️ - 
	`)

}

func main() {
	var wg sync.WaitGroup

	led := machine.GPIO2
	button := machine.GPIO25
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	button.Configure(machine.PinConfig{machine.PinInput})
	alarm := system.New(led)
	alarm.Tone(6000, 1000)

	channel, channel2 := make(chan int), make(chan int)
	go func(enable, disable <-chan int) {
		var status uint8
		defer emergency()
		timer := time.NewTicker(time.Second)
		defer timer.Stop()
		for true {

			select {
			case msg := <-enable:

				fmt.Println("enable", msg, status)
				status = 0
			case msg := <-disable:
				status = 3
				fmt.Println("disable", msg)
			case <-timer.C:
				if status != 0 {
					continue
				}
				led.Set(!led.Get())
				fmt.Println("nothing chose")
			}

		}
	}(channel, channel2)

	i := hardware.NewInterface(button, func(state hardware.EventButton) {
		if state == hardware.PressedIn {
			channel <- 0
		}

	})
	mananger := system.CreateIO(func() {
		log.Println("Hello!")
		channel2 <- 0
	})

	mananger.Configure(system.Config{
		CallMode: system.TimePressedCall,
		During:   time.Second * 3,
	})

	i2 := hardware.NewInterface(button, mananger.Manager)

	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			i2.Process()
			runtime.Gosched()
		}

	}()
	go func() {
		defer wg.Done()
		for {
			i.Process()
			runtime.Gosched()
		}
	}()
	wg.Wait()

}
