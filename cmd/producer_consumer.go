package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

// single pizza make log
type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func makePizza(pizzaNum int) *PizzaOrder {
	pizzaNum++
	if pizzaNum <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Sprintf("Received order #%d!\n", pizzaNum)

		fmt.Sprintf("Making pizza #%d, it will take %d seconds ...\n", pizzaNum, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd <= 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		if rnd < 2 {
			msg = fmt.Sprintf("** We run out of ingredients for pizza #%d!", pizzaNum)
		} else if rnd <= 5 {
			msg = fmt.Sprintf("** The cook is quit while making pizza #%d!", pizzaNum)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d!", pizzaNum)
		}

		return &PizzaOrder{
			pizzaNumber: pizzasMade,
			message:     msg,
			success:     success,
		}

	}

	return &PizzaOrder{
		pizzaNumber: pizzasMade,
	}
}

func pizzaria(pizzaMaker *Producer) {

	// var i = 0

	for {
		// currentPizza := makePizza(i)

	}

}

func main1() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	color.Cyan("The pizza store is open for business!")
	color.Cyan("=====================================")

	pizzaMaker := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}
	go pizzaria(pizzaMaker)

}
