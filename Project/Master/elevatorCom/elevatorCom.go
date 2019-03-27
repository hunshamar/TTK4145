package elevatorCom

import (
	"../../config"
	"../../dataTypes"
	"../../network/bcast"
	"../../timer"
)

func ReceiveAndCheckStatus(RX chan<- dataTypes.ElevatorInfo, elevNumber int, disconnected chan<- int) {
	watchdogTimer := timer.Timer_s{}
	watchdogTimer.Start(config.WatchdogTimeMs)
	RxElev := make(chan dataTypes.ElevatorInfo)
	watchdogTimedOut := make(chan bool)
	go bcast.Receiver(config.ElevatorTXPort+elevNumber, RxElev)
	go watchdogTimer.PollTimer(watchdogTimedOut)

	for {
		select {
		case FromElevator := <-RxElev:
			FromElevator.Number = elevNumber

			if FromElevator.HardwareFunctioning == false {
				FromElevator.State = dataTypes.S_Disconnected
				disconnected <- elevNumber
			}

			RX <- FromElevator

			watchdogTimer.Reset()

		case <-watchdogTimedOut:
			disconnected <- elevNumber
			watchdogTimer.Reset() 
		}
	}
}
