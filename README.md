# Concurrent Ticket Booking
A Ticket Booking Mechanism that allows concurrency while booking tickets. Implemented Locking and Concurrency in GoLang. 
WebApp persists the WebSocket connected to the Golang App which essentially triggers when seats get updated (without loading).
It's like the TATKAL Ticket Booking system where the number of users will fire booking requests concurrently and the Server can manage and allocate the seats to all users according to the seats.

### Architecture
![architectural_figure](https://github.com/omjogani/concurrent-ticket-booking/blob/master/Concurrent-ticket-booking.png?raw=true "Architectural Figure")

### Working Example
   
![concurrent-ticketbooking-demo](https://github.com/omjogani/concurrent-ticket-booking/assets/72139914/18c81613-b208-4d57-9562-81ac4b43b117)



### Technical Details

---
- Technology
    - Concurrency: Mutex & WaitGroups in Golang
    - WebSocket: Golang
    - User Interface: HTML, Tailwind CSS, JS 


>If you found this useful, make sure to give it a star 🌟
## Thank You!!
