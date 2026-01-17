package main

import (
	"fmt"
	"go-concurrent-cafeteria/machine"
)

func main() {
	grinders, expressoMachines, steamers := machine.SetupMachines()

	numberOfOrders := 100
	orderedLatte := make(chan Latte, numberOfOrders)
	for i := range numberOfOrders {
		go makeALatte(i, grinders, expressoMachines, steamers, orderedLatte)
	}

	lattes := []Latte{}
	// Appends the lattes to a slice so we can deliver them!
	for _ = range numberOfOrders {
		lattes = append(lattes, <-orderedLatte)
		fmt.Println(lattes, len(lattes))
	}
}

type Latte struct {
	orderID int
	coffe 	int
	milk	int
}

func makeALatte(orderID int, grinders chan machine.Grinder, expressoMachines chan machine.ExpressoMachine, steamers chan machine.Steamer, lattes chan Latte) {
	lattes<- Latte{
		orderID: orderID,
		coffe: -1,
		milk: -1,
	}
	go makeACoffe(orderID, grinders, expressoMachines, lattes)

	steamMilk(orderID, steamers, lattes)
}

func makeACoffe(orderID int, grinders chan machine.Grinder, expressoMachines chan machine.ExpressoMachine, lattes chan Latte) {
	// Get a Grinder
	grinder := <-grinders
	// Grind Beans !
	grinder.GrindBeans(orderID)
	// Get Grounds from Grinder
	grounds := <-grinder.GroundsQueue
	// Return the Grinder to the channel so others can use it
	grinders<- grinder

	// Get a ExpressoMachine
	expressoMachine := <-expressoMachines
	// Make Expresso !
	expressoMachine.MakeExpresso(grounds)
	// Get coffe from ExpressoMachine
	expresso := <-expressoMachine.CoffeQueue
	// Return the ExpressoMachine to the channel so others can use it
	expressoMachines<- expressoMachine

	// Put coffe on the ordered Latte
	latte := <-lattes
	latte.coffe = expresso
	lattes<- latte
}

func steamMilk(orderID int, steamers chan machine.Steamer, lattes chan Latte) {
	// Get a Steamer
	steamer := <-steamers
	// Steam Milk !
	steamer.SteamMilk(orderID)
	// Get steamed milk from Steamer
	steamedMilk := <-steamer.SteamedMilkQueue	
	// Return Steamer to the channel so others can use it
	steamers<- steamer

	// Put steamed milk on the Latte
	latte := <-lattes
	latte.milk = steamedMilk
	lattes<- latte
}
