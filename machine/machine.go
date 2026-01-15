package machine

import (
	"fmt"
	"time"
	"math/rand"
)

const (
	numberOfGrinders         = 2
	numberOfExpressoMachines = 2
)

type Grinder struct {
	Beans chan int
}

func (g *Grinder) GrindBeans(orderID int) {
	grindingTime := time.Duration(rand.Float32())
	time.Sleep(grindingTime * time.Second) // or time.Millisecond, for quicker simulation
	fmt.Println(fmt.Sprintf("Grinded beans: %d", orderID))
	g.Beans<- orderID
}

type ExpressoMachine struct {
	Coffe chan int
}

func (em *ExpressoMachine) MakeExpresso(grindedBeans int) {
	expressoTime := time.Duration(rand.Float32())
	time.Sleep(expressoTime * time.Second) // or time.Millisecond, for quicker simulation
	fmt.Println(fmt.Sprintf("Made the expresso: %d", grindedBeans))
	em.Coffe<- grindedBeans
}

func SetupMachines() (chan Grinder, chan ExpressoMachine) {
	beans := make(chan int, 100)
	// Creating my Grinder machines channel (so I can limit their use)
	grinders := make(chan Grinder, numberOfGrinders)
	// Setting up the Worker Pool (of Grinders)
	for _ = range numberOfGrinders {
		grinders <- Grinder{ Beans: beans }
	}

	coffe := make(chan int, 100)
	// Crating my Expresso machines channel
	expressoMachines := make(chan ExpressoMachine, numberOfExpressoMachines)
	// Setting up the Worker Pool (of ExpressoMachines)
	for _ = range numberOfExpressoMachines {
		expressoMachines <- ExpressoMachine{ Coffe: coffe }
	}

	return grinders, expressoMachines
}