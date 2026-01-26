package machine

import (
	"fmt"
	"time"
	"math/rand/v2"

	"go-concurrent-cafeteria/telemetry"
)

const (
	numberOfGrinders         = 2
	numberOfExpressoMachines = 2
	numberOfSteamers = 1
)

type Machine struct {
	telemetryService telemetry.TelemetryService
}

type Grinder struct {
	Machine
	GroundsQueue chan int
}

func (g *Grinder) GrindBeans(orderID int) {
	workOnIt()
	fmt.Println(fmt.Sprintf("Grinded beans: %d", orderID))
	g.telemetryService.AddSpan(orderID, "grinded beans")
	g.GroundsQueue<- orderID
}

type ExpressoMachine struct {
	Machine
	CoffeQueue chan int
}

func (em *ExpressoMachine) MakeExpresso(grindedBeans int) {
	workOnIt()
	fmt.Println(fmt.Sprintf("Made the expresso: %d", grindedBeans))
	em.telemetryService.AddSpan(grindedBeans, "made the expresso")
	em.CoffeQueue<- grindedBeans
}

type Steamer struct {
	Machine
	SteamedMilkQueue chan int
}

func (s *Steamer) SteamMilk(orderID int) {
	workOnIt()
	fmt.Println(fmt.Sprintf("Steamed the milk: %d", orderID))
	s.telemetryService.AddSpan(orderID, "steamed the milk")
	s.SteamedMilkQueue<- orderID
}

func workOnIt() {
	workTime := time.Duration(rand.IntN(5) + 1) // Returns a randon Duration between [1, 6) (6 is not included)
	time.Sleep(workTime * time.Second) // or time.Millisecond, for quicker simulation
}

func SetupMachines(telemetryService telemetry.TelemetryService) (chan Grinder, chan ExpressoMachine, chan Steamer) {
	m := Machine{
		telemetryService: telemetryService,
	}
	
	beans := make(chan int, 100)
	// Creating my Grinder machines channel (so I can limit their use)
	grinders := make(chan Grinder, numberOfGrinders)
	// Setting up the Worker Pool (of Grinders)
	for _ = range numberOfGrinders {
		grinders<- Grinder{
			Machine: m,
			GroundsQueue: beans,
		}
	}

	coffeCups := make(chan int, 100)
	// Crating my Expresso machines channel
	expressoMachines := make(chan ExpressoMachine, numberOfExpressoMachines)
	// Setting up the Worker Pool (of ExpressoMachines)
	for _ = range numberOfExpressoMachines {
		expressoMachines<- ExpressoMachine{
			Machine: m,
			CoffeQueue: coffeCups,
		}
	}

	milkCups := make(chan int, 100)
	// Creating my Steamer machines channel
	steamers := make(chan Steamer, numberOfSteamers)
	// Setting up the Steamer Worker Pool
	for _ = range numberOfSteamers {
		steamers<- Steamer{
			Machine: m,
			SteamedMilkQueue: milkCups,
		}
	}

	return grinders, expressoMachines, steamers
}