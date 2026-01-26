package main

import (
	"fmt"
	
	"go-concurrent-cafeteria/machine"
	"go-concurrent-cafeteria/telemetry"
)

const numberOfOrders = 100

func main() {
	telemetryService := telemetry.NewTelemetryService()
	grinders, expressoMachines, steamers := machine.SetupMachines(telemetryService)

	orderedLatte := make(chan Latte)
	for i := range numberOfOrders {
		go makeALatte(i, grinders, expressoMachines, steamers, orderedLatte)
	}

	lattes := []Latte{}
	// Appends the lattes to a slice so we can deliver them!
	for _ = range numberOfOrders {
		lattes = append(lattes, <-orderedLatte)
		fmt.Println("Delivered Latte ", lattes[len(lattes) - 1].orderID)
	}
	fmt.Println(fmt.Sprintf("Delivered %d lattes: ", len(lattes)), lattes)
	telemetryService.PrintTelemetry()
}

type Latte struct {
	orderID int
	coffe 	int
	milk	int
}

func (l *Latte) IsDone() bool {
	// The Latte is done when it have both coffe and milk set
	return l.coffe != -1 && l.milk != -1
}

func makeALatte(orderID int, grinders chan machine.Grinder, expressoMachines chan machine.ExpressoMachine, steamers chan machine.Steamer, orderedReadyLatte chan Latte) {
	unfinishedLattes := make(chan Latte, numberOfOrders)
	unfinishedLattes<- Latte{
		orderID: orderID,
		coffe: -1,
		milk: -1,
	}
	go makeACoffe(orderID, grinders, expressoMachines, unfinishedLattes, orderedReadyLatte)

	steamMilk(orderID, steamers, unfinishedLattes, orderedReadyLatte)
}

func makeACoffe(orderID int, grinders chan machine.Grinder, expressoMachines chan machine.ExpressoMachine, unfinishedLattes chan Latte, readyLattes chan<- Latte) {
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
	latte := <-unfinishedLattes
	latte.coffe = expresso
	
	// Ready Latte for delivery or keep working on it!
	finishLatte(latte, unfinishedLattes, readyLattes)
}

func steamMilk(orderID int, steamers chan machine.Steamer, unfinishedLattes chan Latte, readyLattes chan<- Latte) {
	// Get a Steamer
	steamer := <-steamers
	// Steam Milk !
	steamer.SteamMilk(orderID)
	// Get steamed milk from Steamer
	steamedMilk := <-steamer.SteamedMilkQueue	
	// Return Steamer to the channel so others can use it
	steamers<- steamer

	// Put steamed milk on the Latte
	latte := <-unfinishedLattes
	latte.milk = steamedMilk
	
	// Ready Latte for delivery or keep working on it!
	finishLatte(latte, unfinishedLattes, readyLattes)
}

func finishLatte(latte Latte, unfinishedLattes chan Latte, readyLattes chan<- Latte) {
	// Check  if the Latte is ready, if so put it in the ready chan !
	if latte.IsDone() {
		readyLattes<- latte
	} else {
	//  If it is not finished, put it in the chan and let the other worker finish its part
		unfinishedLattes<- latte
	}
}
