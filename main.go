package main

import (
	"fmt"
	"go-concurrent-cafeteria/machine"
)

func main() {
	grinders, expressoMachines := machine.SetupMachines()

	numberOfOrders := 100
	orderedLatte := make(chan int, numberOfOrders)
	for i := range numberOfOrders {
		go makeACoffe(i, grinders, expressoMachines, orderedLatte)
	}

	lattes := []int{}
	for _ = range numberOfOrders {
		lattes = append(lattes, <-orderedLatte)
		fmt.Println(lattes, len(lattes))
	}
}

func makeACoffe(orderID int, grinders chan machine.Grinder, expressoMachines chan machine.ExpressoMachine, latte chan<- int) {
	// Get a Grinder
	grinder := <-grinders
	grinder.GrindBeans(orderID)
	// Return the Grinder to the channel so others can use it
	grinders<- grinder
	// Get beans from Grinder
	beans := <-grinder.Beans

	// Get a ExpressoMachine
	expressoMachine := <-expressoMachines
	expressoMachine.MakeExpresso(beans)
	// Return the ExpressoMachine to the channel so others can use it
	expressoMachines<- expressoMachine
	// Get coffe from ExpressoMachine
	latte<- <-expressoMachine.Coffe
}
