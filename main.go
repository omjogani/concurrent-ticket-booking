package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

var seats = make([][]string, 6)
var mutexes [6][20]sync.Mutex

func processSeats() {
	for i := range seats {
		seats[i] = make([]string, 20)
	}
	for i := 0; i < len(seats); i++ {
		for j := 0; j < len(seats[0]); j++ {
			seats[i][j] = "⋅"
		}
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
	for {
		for row := 0; row < len(seats[0]); row++ {
			for col := 0; col < len(seats); col++ {
				mutexes[col][row].Lock()
				if seats[col][row] == "⋅" {
					// seat is unoccupied & free to book
					seats[col][row] = "x"
					mutexes[col][row].Unlock()
					return
				}
				mutexes[col][row].Unlock()
			}
		}
	}
}

func allocateSeat(totalCustomer int) {
	var wg sync.WaitGroup
	wg.Add(totalCustomer)

	for idx := 1; idx <= totalCustomer; idx++ {
		go func(seatNumber int) {
			defer wg.Done()
			getMeASeat(seatNumber)
		}(idx)
	}

	// wait till every seats get allocated
	wg.Wait()
}

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) bookTickets(ws *websocket.Conn) {
	// ask the system how many tickets you want to book
	fmt.Println("New Incoming Connection from User: ", ws.RemoteAddr())
	numberOfTicket := strings.Split(ws.Request().URL.RawQuery, "=")[1]

	if numberOfTicket == "" {
		payload := "Number of Ticket not mentioned!"
		ws.Write([]byte(payload))
	} else {
		numberOfTicketInt, err := strconv.Atoi(numberOfTicket)
		if err != nil {
			fmt.Println("Error in converting numberOfTicket to int")
		} else {
			processSeats()
			allocateSeat(numberOfTicketInt)
			displaySeats()
			payload := "Tickets booked!"
			ws.Write([]byte(payload))
		}
	}
}

func main() {
	server := NewServer()
	http.Handle("/book-ticket", websocket.Handler(server.bookTickets))
	defer displaySeats()
	http.ListenAndServe(":3550", nil)
}
