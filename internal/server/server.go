package server

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"sync"
)

type Client struct {
	channel chan string
	name    string
}

type Server struct {
	Listener   net.Listener
	Clients    map[Client]bool
	Broadcast  chan string
	Register   chan Client
	Unregister chan Client
	Mutex      *sync.Mutex
}

func NewServer(listener net.Listener) Server {
	return Server{
		Listener:   listener,
		Clients:    make(map[Client]bool),
		Broadcast:  make(chan string),
		Register:   make(chan Client),
		Unregister: make(chan Client),
		Mutex:      &sync.Mutex{},
	}

}

func (s Server) HandleConnections() {
	for {
		select {
		case message := <-s.Broadcast:
			s.Mutex.Lock()
			for client := range s.Clients {
				select {
				case client.channel <- message:
					respondToMessage(message)
				default:
					close(client.channel)
					delete(s.Clients, client)
				}
			}
			s.Mutex.Unlock()
		case client := <-s.Register:
			s.Mutex.Lock()
			s.Clients[client] = true
			s.Mutex.Unlock()
		case client := <-s.Unregister:
			s.Mutex.Lock()
			if _, ok := s.Clients[client]; ok {
				delete(s.Clients, client)
				close(client.channel)
			}
			s.Mutex.Unlock()
		}
	}
}

func (s Server) HandleClient(conn net.Conn) {
	defer conn.Close()

	channel := make(chan string)
	client := Client{channel: channel, name: conn.RemoteAddr().String()}

	s.Register <- client

	go func() {
		for message := range channel {
			fmt.Fprintln(conn, message)
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		s.Broadcast <- fmt.Sprintf("%s: %s", client.name, message)
	}

	s.Unregister <- client
	fmt.Printf("Client %s disconnected\n", client.name)
}

func respondToMessage(message string) {
	// fmt.Println("message is: ", message) // delete l8r
	fmt.Println("pre parsed message: ", message) // delete l8r
	// Example input string
	input := message

	// Define the regular expression pattern to match the message
	re := regexp.MustCompile(`\[::1\]:\d+: (.+)`)

	// Find the match
	match := re.FindStringSubmatch(input)

	// Check if a match is found
	if len(match) > 1 {
		message := match[1]
		fmt.Println("Parsed message:", message, ":") // delete l8r
	} else {
		fmt.Println("No match found")
	}
}
