package elevatorCom

import(
	"../../network/bcast"
	"../../timer"
	"../../dataTypes"
	"../../config"
)







func Receive(RX chan<- dataTypes.ElevatorInfo, elevNumber int, disconnected chan<- int){
	watchdogTimer := timer.Timer_s{}
	timer.WatchdogInit(int64(1000),&watchdogTimer)
	RxElev 				:= make(chan dataTypes.ElevatorInfo)
	watchdogTimedOut 	:= make(chan bool)
	go bcast.Receiver(config.ElevatorTXPort+elevNumber, RxElev)
	go timer.WatchdogPoll(watchdogTimedOut,&watchdogTimer)
	for {
		select {
			// Greit at heiser ikke kan ta imot hall bestillinger når hardwarefunctioning = false? Tror lett å fikse
		case FromElevator := <-RxElev:
			FromElevator.Number = elevNumber

			if FromElevator.HardwareFunctioning == false{
				FromElevator.State = dataTypes.S_Disconnected
				disconnected <- elevNumber
			}

			RX <- FromElevator
			

			timer.WatchdogReset(&watchdogTimer)

		case <-watchdogTimedOut:
			disconnected <- elevNumber
		}
		
	}



}



