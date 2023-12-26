package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

func (s *Server) broadCastSeatsUpdate(payload []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(payload); err != nil {
				fmt.Println("Write Error:", err)
			}
		}(ws)
	}
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		_, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				// Connection from other side is closed!
				break
			}
			fmt.Println("Read Error: ", err)
			continue
		}
		msg, err := json.Marshal(seats)
		if err != nil {
			fmt.Println(err)
		}
		s.broadCastSeatsUpdate(msg)
	}
}

func (s *Server) getSeatStatusUpdate(ws *websocket.Conn) {
	fmt.Println("Seats Status Request from : ", ws.RemoteAddr())
	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *Server) bookTotalTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requested")
	numberOfTicket := r.URL.Query().Get("tickets")

	if numberOfTicket == "" {
		// not specified number of tickets to book
		return
	} else {
		numberOfTicketInt, err := strconv.Atoi(numberOfTicket)
		if err != nil {
			fmt.Println("Error in converting numberOfTicket to int")
		} else {
			allocateSeat(numberOfTicketInt)
			displaySeats()
			payload, err := json.Marshal(seats)

			s.broadCastSeatsUpdate(payload)
			if err != nil {
				fmt.Println("Error in Marshal Operation")
			}
		}
	}
}

func main() {
	processSeats()
	server := NewServer()
	PORT := ":3550"
	http.Handle("/", websocket.Handler(server.getSeatStatusUpdate))
	http.HandleFunc("/book-ticket", server.bookTotalTicket)

	fmt.Println("Server is listening at", PORT)
	http.ListenAndServe(PORT, nil)
}
