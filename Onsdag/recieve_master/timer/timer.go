
package timer

import "time"


const _pollRate = 20 * time.Millisecond

type Timer_s struct{
    start_time_ms int64 
    run_time_ms int64 
}

func WatchdogTimedOut(w* Timer_s) bool{


    if ((w.start_time_ms + w.run_time_ms) < (time.Now().UnixNano() / int64(time.Millisecond))){
        return true
    }
    return false
}


func WatchdogInit(time_ms int64, w* Timer_s){
	w.start_time_ms = time.Now().UnixNano() / int64(time.Millisecond)
	w.run_time_ms = time_ms
}

func WatchdogPoll(receiver chan<- bool, w* Timer_s){
	for {
		time.Sleep(_pollRate)
		v := WatchdogTimedOut(w)
		if v == true {
			receiver <- true
		}
	}
}

func WatchdogReset(w* Timer_s){
	w.start_time_ms = time.Now().UnixNano() / int64(time.Millisecond)	
}