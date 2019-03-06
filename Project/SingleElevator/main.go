package main


import "./FSM"




func main(){

    FSM.Loop()


    //go timer.Start(door_timer, run_time_ms)
    //go timer.TimedOut(door_timer) 

    /*
    
    //buttons
    buttons_poll() // sjekker alle knapper, skrur på lys når trykket

    //Orders
    etasjepanel_opp[4]  = {0,0,0,x}
    etasjepanel_ned[4] = {x,0,0,0}
    heispanel[4] = {0,0,0,0}

    //elevator_control_logic


    //fsm


    //oversikt over tilstander heisen kan være i
    struct states{
        currious_floor // 1 2 3 4 
        idle
        moving_up
        moving_down
        door_open
    }states_t
    

    //Timer

    timer_Set(ms)
    bool checktimer()

    */


    FSM.Init()




     
}
