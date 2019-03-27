package main


import(
    "./FSM"
    "os"
    "fmt"
    "strconv"
)


func main(){

    elevNum := os.Args[1:][0]
    portNum := os.Args[1:][1] 

    elevatorNumber,err := strconv.Atoi(elevNum)
    if err != nil || elevatorNumber < 1{
        fmt.Println("Error. Not valid elevator number")
        os.Exit(1)
    }

    fmt.Println("elev:",elevatorNumber, "on port", portNum)
    FSM.StateMachine(elevatorNumber, portNum)    
}
