package message

import (
	"fmt"
	"regexp"
)

type Message struct {
	Text string
}

func NewMessage(text string) Message {
	return Message{Text: text}
}

func (m Message) RespondToMessage() {
	// fmt.Println("message is: ", message) // delete l8r
	fmt.Println("pre parsed message: ", m.Text) // delete l8r

	// Define the regular expression pattern to match the message
	re := regexp.MustCompile(`\[::1\]:\d+: (.+)`)

	// Find the match
	match := re.FindStringSubmatch(m.Text)

	// Check if a match is found
	if len(match) > 1 {
		m.Text = match[1]
		fmt.Println("Parsed message:", m.Text, ":") // delete l8r
	} else {
		fmt.Println("No match found")
	}
}
