package main

import (
	"fmt"
	"time"
)

type Grinder struct {}

func Grind(grinders chan Grinder, orderID int, results chan<- int) {
	m := <- grinders
	grindingTime := time.Duration(1)
	time.Sleep(grindingTime * time.Second)
	fmt.Println(fmt.Sprintf("Grinded beans: %d", orderID))
	results <- orderID
	grinders <- m
}

func main() {
	// Creating my Grinder machines
	numberOfGrinders := 2
	grinders := make(chan Grinder, numberOfGrinders)
	for _ = range numberOfGrinders {
		grinders <- Grinder{}
	}

	numberOfOrders := 100
	results := make(chan int)
	for i := range numberOfOrders {
		go Grind(grinders, i, results)
	}
	for counter := 0; counter < 100; counter++ {
		<- results
	}
}
