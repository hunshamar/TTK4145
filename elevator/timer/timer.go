
package timer

import "time"


const _pollRate = 20 * time.Millisecond

type Door_timer struct{
    start_time_ms int64 
    run_time_ms int64 
}

var t = Door_timer{0,0}

func Start(run_time_ms int64) {
        
    t.start_time_ms = time.Now().UnixNano() / int64(time.Millisecond)

    t.run_time_ms = run_time_ms

}

func TimedOut() bool{


    if ((t.start_time_ms + t.run_time_ms) < (time.Now().UnixNano() / int64(time.Millisecond))){
        return true
    }
    return false
}

func PollTimer(receiver chan<- bool){
	prev := false
	for {
		time.Sleep(_pollRate)
		v := TimedOut()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}

func Hello() string{
    return "world"
}

