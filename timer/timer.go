
package timer

import "time"



type Door_timer struct{
    start_time_ms int64
    run_time_ms int64
}

func Start(t *Door_timer, run_time_ms int64) {
        
    t.start_time_ms = time.Now().UnixNano() / int64(time.Millisecond)

    t.run_time_ms = run_time_ms

}

func TimedOut(t Door_timer) bool{


    if ((t.start_time_ms + t.run_time_ms) < (time.Now().UnixNano() / int64(time.Millisecond))){
        return true
    }
    return false
}

