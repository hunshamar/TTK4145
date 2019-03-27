
package config

const DoorOpenTimeMs  int 		= 3000
const ElevatorTravelTimeMs int 	= 3500
const WatchdogTimeMs int 		= 1000
const NumFloors int 			= 4
const NumElevators int 			= 3
const NumOrderButtons int 		= 3
const MasterTXPort int 			= 16569
const ElevatorTXPort int 		= 16500
const BackupRXPort int	    	= 15600
const InfCost int 				= 4*(ElevatorTravelTimeMs + DoorOpenTimeMs )*NumFloors 