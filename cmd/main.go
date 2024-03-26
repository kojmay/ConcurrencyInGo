package main

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

// barber type
type Barber struct {
	name       string
	number     int
	isAsleep   bool
	haircutNum int
}

// customer stype
type Customer struct {
	number int
}

// the waitingChan serves a waiting room
var waitingChan chan *Customer

// the barber chan
var barberChan chan *Barber

// barber number
const barberNum = 2
const haircutTime = 2
const waitingRoomCapacity = 5

// time to close the Shop
var closeChan chan interface{}

// simulate creating createCustomers
func createCustomers(ctx context.Context, wg *sync.WaitGroup) {
	color.Red("Start creating customers\n")
	i := 0
	for {
		select {
		case <-ctx.Done():
			// recieved shop close signal
			color.Cyan("Close signal recieved, stop creating customers...\n")
			defer wg.Done()
			return
		default:
			// simulate creating customers irregularly
			randTime := rand.Intn(6)
			time.Sleep(time.Duration(randTime) * time.Second)
			color.Cyan("A new customer comes...\n")
			i++
			if len(waitingChan) == waitingRoomCapacity {
				color.Cyan("\t The room is full, so the #%d customer leaves...\n", i)
				continue
			} else {
				customer := &Customer{
					number: i,
				}
				waitingChan <- customer
				color.Cyan("\t #%d customer sits down, there are %d sits left.\n", i, waitingRoomCapacity-len(waitingChan))
				// color.Cyan("\t waiting chan len: %d", len(waitingChan))
			}
		}
	}
}

// barbers start to check the room, and choose to sleep or work
func barbersAppear(ctx context.Context, wg *sync.WaitGroup) {
	color.Red("Barbers start to work.\n")
	for {
		// color.Black("waiting chan len: %d", len(waitingChan))
		if len(waitingChan) > 0 {
			for customer := range waitingChan {
				// color.Red("Test: customer %d, baber: %d", customer.number, )
				for barber := range barberChan {
					wg.Add(1)
					go working(barber, customer, wg)
					break
				}
			}
		}
		select {
		case <-ctx.Done():
			color.Red("Close signal recieved")
			defer wg.Done()
			return
		default:
			continue
		}
	}
}

// barber start cutting hair for customer
func working(barber *Barber, customer *Customer, wg *sync.WaitGroup) {
	defer wg.Done()
	if barber.isAsleep {
		color.Blue("\t#%d barber wakes up, because of #%d customer\n", barber.number, customer.number)
	}
	barber.haircutNum++
	barber.isAsleep = false
	color.Blue("\t#%d customer starts taking a hair cut, his/her barber is #%d barber\n", customer.number, barber.number)
	time.Sleep(haircutTime * time.Second)
	color.Blue("\t#%d customer is leaving, it's #%d barber's #%d customer.\n", customer.number, barber.number, barber.haircutNum)
	if len(waitingChan) == 0 {
		// no customer
		color.Blue("\t#%d Barber is going to sleep as there is no customer in room \n", barber.number)
		barber.isAsleep = true
	}
	barberChan <- barber
}

func clockTimer(cancel context.CancelFunc) {
	time.Sleep(20 * time.Second)
	cancel()
	color.Red("Close signal prodcasted!!!")
}
func main() {
	// init barberChan
	barberChan = make(chan *Barber, barberNum)
	for i := 0; i < barberNum; i++ {
		barberChan <- &Barber{
			number:     i + 1,
			isAsleep:   true,
			haircutNum: 0,
		}
	}
	waitingChan = make(chan *Customer, waitingRoomCapacity)
	// closeChan = make(chan interface{})

	ctx, cancel := context.WithCancel(context.Background())
	go clockTimer(cancel)
	var wg sync.WaitGroup
	// start create cuctomers
	wg.Add(2)
	go createCustomers(ctx, &wg)
	// barbers start to check the room, and choose to sleep or work according to the situation
	go barbersAppear(ctx, &wg)
	wg.Wait()
}
