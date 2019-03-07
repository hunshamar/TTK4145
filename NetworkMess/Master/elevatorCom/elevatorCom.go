package elevatorCom

import(
	/*"../bcast"
	"../localip"
	"../peers"
	*/
	"../../dataTypes"
	"time"
)

func Transmit(TX chan<- dataTypes.LongMessage, LongMessage <-chan dataTypes.LongMessage){

		for {
			ToElevators := <-LongMessage
			TX <- ToElevators
			time.Sleep(100 * time.Millisecond)
		}
}

