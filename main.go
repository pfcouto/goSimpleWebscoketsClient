package main

import (
    "bufio"
	"encoding/json"
    "fmt"
    "log"
    "os"
    "github.com/gorilla/websocket"
	"net/http"
)

type Message struct {
    Topic string      `json:"topic"`
    Data  interface{} `json:"data"`
}

func main() {
    // WebSocket server URL
    serverURL := "ws://127.0.0.1:3000"

    // Additional headers, including the "cookie" header
    headers := http.Header{"Cookie": {"topics=robot,ai;device=goClient"}}

    // Establish WebSocket connection
    conn, _, err := websocket.DefaultDialer.Dial(serverURL, headers)
    if err != nil {
        log.Fatal("WebSocket connection error:", err)
    }
    defer conn.Close()

    // Read keyboard input in a loop
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Enter a command (H/L/R): ")
        scanner.Scan()
        input := scanner.Text()
        if err := scanner.Err(); err != nil {
            log.Fatal("Error reading input:", err)
        }

        // Handle different inputs
        switch input {
        case "H", "L", "R":
            // Send the command
            // Create a JSON message
            message := Message{
                Topic: "robot",
                Data:  map[string]string{"command": input},
            }

			jsonMessage, err := json.Marshal(message)

			if err != nil {
                log.Fatal("Error marshaling message to JSON:", err)
            }

            // Send the JSON message
            if err := conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
                log.Fatal("Error sending message:", err)
            }

            // Wait for ACK
            _, response, err := conn.ReadMessage()
            if err != nil {
                log.Fatal("Error reading response:", err)
            }

            // Unmarshal the ACK message
            var ackMessage string
            if err := json.Unmarshal(response, &ackMessage); err != nil {
                log.Fatal("Error unmarshaling ACK message:", err)
            }
			
			
			
            if input != "H" && ackMessage == "ACK" {
				fmt.Println("ACK received")
                aiMessage := Message{
                    Topic: "ai",
                    Data:  "D",
                }

                // Marshal the AI message to JSON
                jsonAI, err := json.Marshal(aiMessage)
                if err != nil {
                    log.Fatal("Error marshaling AI message to JSON:", err)
                }

                // Send the JSON AI message
                if err := conn.WriteMessage(websocket.TextMessage, jsonAI); err != nil {
                    log.Fatal("Error sending AI message:", err)
                }
            }
        default:
            fmt.Println("Invalid command. Please enter H, L, or R.")
        }
    }
}
