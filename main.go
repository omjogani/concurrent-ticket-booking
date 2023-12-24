package main

import (
	"fmt"
	"sync"
)

var seats = make([][]int, 6)
var mutexes [6][20]sync.Mutex

func processSeats() {
	for i := range seats {
		seats[i] = make([]int, 20)
	}
}

func displaySeats() {
	divider := 0
	for _, row := range seats {
		if divider == 3 {
			fmt.Println()
		}
		for _, value := range row {
			fmt.Print(value, "  ")
		}
		divider++
		fmt.Println()
	}
}

func getMeASeat(seatNumber int) {
	// if seatNumber == 0 {
	// start from starting
	status := false
	for row := 0; row < len(seats); row++ {
		for col := 0; col < len(seats[0]); col++ {
			// mutexes[row][col].Lock()
			if seats[row][col] == 1 {
				continue
			}
			seats[row][col] = seatNumber
			// mutexes[row][col].Unlock()
			status = true
			break
		}
		if status {
			break
		}
	}
	// } else {
	// get me specific seat
	// }
}

func allocateSeat(totalCustomer int) {
	for idx := 1; idx <= totalCustomer; idx++ {
		go func() {
			getMeASeat(idx)
		}()
	}
}

// func allocateSeatRandomly(totalCustomer int) {
// 	for idx := range totalCustomer {
// 		go func () {

// 		}
// 	}
// }

func main() {
	processSeats()
	defer displaySeats()
	allocateSeat(120)
}
