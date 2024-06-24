package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	channel chan string
	name    string
}

var (
	clients    = make(map[Client]bool)
	broadcast  = make(chan string)
	register   = make(chan Client)
	unregister = make(chan Client)
	mutex      = &sync.Mutex{}
)

func main() {
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer server.Close()

	go handleConnections()

	fmt.Println("Server listening on port 8080")
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(conn)
	}
}

func handleConnections() {
	for {
		select {
		case message := <-broadcast:
			mutex.Lock()
			for client := range clients {
				select {
				case client.channel <- message:
				default:
					close(client.channel)
					delete(clients, client)
				}
			}
			mutex.Unlock()
		case client := <-register:
			mutex.Lock()
			clients[client] = true
			mutex.Unlock()
		case client := <-unregister:
			mutex.Lock()
			if _, ok := clients[client]; ok {
				delete(clients, client)
				close(client.channel)
			}
			mutex.Unlock()
		}
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	channel := make(chan string)
	client := Client{channel: channel, name: conn.RemoteAddr().String()}

	register <- client

	go func() {
		for message := range channel {
			fmt.Fprintln(conn, message)
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		broadcast <- fmt.Sprintf("%s: %s", client.name, message)
	}

	unregister <- client
	fmt.Printf("Client %s disconnected\n", client.name)
}
