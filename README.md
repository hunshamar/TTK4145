Elevator Project
================


Summary
-------
Software for controlling `n` elevators working in parallel across `m` floors written in Go

How to Run 
-----------------
For each elevator there are three files that need to run:

Use `chmod +x <filename>`in order to give permission to run the files.


- A simulator or a network connection to a hardware elevator, both located in the repository. To run the simulator write for example `./SimElevatorServer --port 15657 --numfloors 4` If no floor or port is specified the port will be 15657 and numfloors will be 4. 

- An elevator executable named `SingleElevator`. To run the elevator software, write `./SingleElevator`followed by an elevator number and the port for connecting to the TPC elevator. For example for running three elevators on three different simulator ports write `./SingleElevator 1 15657`, `./SingleElevator 2 15658` and `./SingleElevator 3 15659`. 

- A master/backup file named `Master`. This also needs a unique number. Use 1,2,3, the same as the elevators when running. For example `./Master 1`, `./Master 2` etc. Master 1 will start as master, while the others will become backups, if the master crashes the backup with highest priority (number 2) will become master. If multiple masters are connected to the network simultaneously, the master with the lowest number (highest priority) will become master. 

All of the ports and constants, like number of elevators, number of floors, door open time, UDP-ports etc. are specified in the `config` file. If they are changed the elevator and master files have to be rebuildt by running the `BuildAll`file. 

Short Explanation of the Software Design
----------------------------------------

Master-Slave system. 

`n` elevators communicating through UDP with a master. The elevators sends their orders, floor, state and direction to the master at 10 Hz . The master will send back messages to the elevators that they are to either handle the order or just turn heir hall light on to indicate to the passenger waiting for the elevator that an elevator is coming. 

The master has `n` watchdog timers. The timers are reset every time the master receives a message from a specific elevator. If a second passes without receiving a message, the elevator is marked as disconnected by the master until it starts receiving messages again. 

In each elevator there is a watchdogtimer that starts every time the elevator starts moving, if more than 5 seconds (2x times the travel time between floors) passes without reaching a floor, a hardware alert is send to the master, and the elevator is marked as disconnected. 

A disconnected elevator will not receive or be able to send hall orders to the master, but it will still be able to handle cab orders, as long as the hardware of the elevator is functioning. 
When an elevator disconnects from the master, the master will still keep track of the elevators cab orders, so they will not be lost if the elevator software crashes completely. 

When an elevator disconnects from the master, the master will redistribute the orders to a connected elevator. If there are no connected elevators, the master will give the orders to the first elevator to connect to the network again.

Each time the master receives the info from an elevator it will pass it through an algorithm to determine which elevator is to handle the orders received, and then broadcast out to the elevators which orders they are to handle. The master transmit frequency is a dependent of the elevator transmit frequency. If it receives 10 messages a second, it will transmit 10 messages back to the elevators a second. 

Both the messages to and from the master double as an "I am alive" signal and provides information about the orders and state of the elevators. This will ensure that the system is robust for high packet loss. 

To ensure that for example orders updated in the elevator are not immediately overwritten by the messages from the master, the orders have five different states. These different order states have logic such that for example an `order received`can not be overwritten by an `enmpty order` before it is either turned into a `turn light on`or `handle order`, ensuring no hall orders are lost.

- `0: Empty order` 
- `1: Order received`, when an elevator receives an order it will send this status to the master and then the master will either send back `2` or `3` (only one elevator will receive `3`, the rest will receive `2`)
- `-1: Executed order`, when an elevator has executed the order it will send back this value to the master, so it can remove it from the order list. 
- `2: Turn on light`, turn the hall lights on 
- `3: Handle order`, go to the floor and execute the order

 
If the master crashes the backups have watchdog timers with a time scaled with their priority. Backup 2 will take over 2 seconds after the master disconnects. This will give a slight delay to the hall orders if any are received during this interval, but the orders will not be lost. The backup will always keep track of the orders transmitted by the elevators, but it will not send anything out before its state is switched to master. For the master and backups to be able communicate and separate themselves from each other the master will broadcast its unique ID/priority, and the backups will broadcast a backup ID set to -1. 

The cost function is based on which elevator becomes idle first after receiving the order, simulating the path taken by all the elevators, using the `elevatorLogic` module in `SingleElevator`. 


Libraries/Code Not Written by the Group
---------------------------------------
The `Network`-module and `elevio`-module are taken from [github.com/TTK4145](https://github.com/TTK4145). The `costFunction` module is based on the `Time until completion/idle` cost function at [github.com/TTK4145/Project-resources](https://github.com/TTK4145/Project-resources). 
