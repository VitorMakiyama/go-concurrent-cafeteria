package machine

import (
	"fmt"
	"time"
	"math/rand"
)

const (
	numberOfGrinders         = 2
	numberOfExpressoMachines = 2
	numberOfSteamers = 1
)

type Grinder struct {
	GroundsQueue chan int
}

func workOnIt() {
	workTime := time.Duration(rand.Float32())
	time.Sleep(workTime * time.Second) // or time.Millisecond, for quicker simulation
}

func (g *Grinder) GrindBeans(orderID int) {
	workOnIt()
	fmt.Println(fmt.Sprintf("Grinded beans: %d", orderID))
	g.GroundsQueue<- orderID
}

type ExpressoMachine struct {
	CoffeQueue chan int
}

func (em *ExpressoMachine) MakeExpresso(grindedBeans int) {
	workOnIt()
	fmt.Println(fmt.Sprintf("Made the expresso: %d", grindedBeans))
	em.CoffeQueue<- grindedBeans
}

type Steamer struct {
	SteamedMilkQueue chan int
}

func (s *Steamer) SteamMilk(orderID int) {
	workOnIt()
	fmt.Println(fmt.Sprintf("Steamed the milk: %d", orderID))
	s.SteamedMilkQueue<- orderID
}

func SetupMachines() (chan Grinder, chan ExpressoMachine, chan Steamer) {
	beans := make(chan int, 100)
	// Creating my Grinder machines channel (so I can limit their use)
	grinders := make(chan Grinder, numberOfGrinders)
	// Setting up the Worker Pool (of Grinders)
	for _ = range numberOfGrinders {
		grinders<- Grinder{ GroundsQueue: beans }
	}

	coffeCups := make(chan int, 100)
	// Crating my Expresso machines channel
	expressoMachines := make(chan ExpressoMachine, numberOfExpressoMachines)
	// Setting up the Worker Pool (of ExpressoMachines)
	for _ = range numberOfExpressoMachines {
		expressoMachines<- ExpressoMachine{ CoffeQueue: coffeCups }
	}

	milkCups := make(chan int, 100)
	// Creating my Steamer machines channel
	steamers := make(chan Steamer, numberOfSteamers)
	// Setting up the Steamer Worker Pool
	for _ = range numberOfSteamers {
		steamers<- Steamer{ SteamedMilkQueue: milkCups }
	}

	return grinders, expressoMachines, steamers
}