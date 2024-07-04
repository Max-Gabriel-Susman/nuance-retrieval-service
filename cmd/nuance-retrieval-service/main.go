package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"sync"

	"github.com/pinecone-io/go-pinecone/pinecone"
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
	ctx := context.Background()

	pineKey := os.Getenv("PINECONE_API_KEY")

	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: pineKey,
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	idxs, err := pc.ListIndexes(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, index := range idxs {
		fmt.Println(index)
	}

	idx, err := pc.Index(idxs[0].Host)
	defer idx.Close()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	res, err := idx.DescribeIndexStats(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(res)

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
		fmt.Println("Parsed message:", message, ":")
	} else {
		fmt.Println("No match found")
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
					respondToMessage(message)
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
