package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func getAvailableShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/zsh"
	}
	return shell
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket connection")

	// Send service info to the client
	serviceInfo := "Service Information" // Replace with actual service info
	serviceMsg := Message{Type: "message", Data: serviceInfo}
	if err := conn.WriteJSON(serviceMsg); err != nil {
		log.Println("Failed to send service info:", err)
		return
	}

	shell := getAvailableShell()

	cmd := exec.Command(shell)
	tty, err := pty.Start(cmd)
	if err != nil {
		log.Println("Failed to start shell:", err)
		return
	}
	defer func() {
		cmd.Process.Kill()
		tty.Close()
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := tty.Read(buf)
			if err != nil {
				return
			}
			outputMsg := Message{Type: "output", Data: string(buf[:n])}
			if err := conn.WriteJSON(outputMsg); err != nil {
				log.Println("Failed to send output:", err)
				return
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		if msg.Type == "input" {
			if _, err := tty.Write([]byte(msg.Data)); err != nil {
				log.Println("Write to TTY error:", err)
				break
			}
		}
	}

	log.Println("User disconnected")
}
