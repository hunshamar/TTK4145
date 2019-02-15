package main

import "./elevio"
import "./timer"
import "./orders"
import "fmt"



func main(){

    numFloors := 4


    elevio.Init("localhost:15657", numFloors)
    
    var d elevio.MotorDirection = elevio.MD_Up
    //elevio.SetMotorDirection(d)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)    
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)
    

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

    
    elevio.SetMotorDirection(elevio.MD_Down)
    
    for(elevio.GetFloor() == -1){
        fmt.Printf("Hello %s", timer.Hello())
    }
    elevio.SetMotorDirection(elevio.MD_Stop)



    /*
    for(-1 == <-drv_floors){
        elevio.SetMotorDirection(d)
        fmt.Printf("while");
    }
    */

    /* ---------------*/


    for {   
        select {
        case a := <- drv_buttons:
            fmt.Printf("%+v\n", a)
            elevio.SetButtonLamp(a.Button, a.Floor, true)
            orders.Add(int(a.Button), int(a.Floor))
            
        case a := <- drv_floors:
            fmt.Printf("%+v  ", a)

            orders.PrintOrders()

            //t := timer.Door_timer
        
            

            timer.Start(3000)
            elevio.SetMotorDirection(elevio.MD_Stop)
            for (!timer.TimedOut()){

            }

            


            elevio.SetMotorDirection(d)


            if a == numFloors-1 {
                d = elevio.MD_Down
            } else if a == 0 {
                d = elevio.MD_Up
            }
            elevio.SetMotorDirection(d)
            
            
        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }
            
        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        }
    }    
}
