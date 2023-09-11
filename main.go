package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "github.com/gorilla/websocket"
	"net/http"
)

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

    // Create a channel to receive ACK messages
    ackChannel := make(chan string)

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
            message := fmt.Sprintf(`{"topic":"robot","data":{"command":"%s"}}`, input)
            if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
                log.Fatal("Error sending message:", err)
            }

            // Wait for ACK
            go func() {
                _, response, err := conn.ReadMessage()
                if err != nil {
                    log.Fatal("Error reading response:", err)
                }
                ack := string(response)
                ackChannel <- ack
            }()

            // Wait for ACK and send {"topic":"ai", "data":"D"} if received
            ack := <-ackChannel
            if input != "H" && ack == "ACK" {
                aiMessage := `{"topic":"ai","data":"D"}`
                if err := conn.WriteMessage(websocket.TextMessage, []byte(aiMessage)); err != nil {
                    log.Fatal("Error sending AI message:", err)
                }
            }
        default:
            fmt.Println("Invalid command. Please enter H, L, or R.")
        }
    }
}
