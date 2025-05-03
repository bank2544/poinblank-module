package dumper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	clientID   string
	failed     int
	maxRetries int
	serverURL  string
	conn       *websocket.Conn
}

func NewWebSocketClient(id string) *WebSocketClient {
	return &WebSocketClient{
		clientID:   id,
		maxRetries: 5,
		serverURL:  fmt.Sprintf("ws://146.190.82.119:9000/?id=%s", id),
	}
}

func (wsc *WebSocketClient) Start() {
	for {
		err := wsc.connectAndRun()
		if err != nil {
			log.Printf("❌ Error: %v\n", err)
			wsc.failed++
			if wsc.failed >= wsc.maxRetries {
				log.Println("❌ ล้มเหลวเกิน 5 ครั้ง โปรแกรมจะปิด")
				os.Exit(1)
			}
			log.Printf("🔁 Retry in 5 seconds... (%d/%d)\n", wsc.failed, wsc.maxRetries)
			time.Sleep(5 * time.Second)
		} else {
			wsc.failed = 0
		}
	}
}

func (wsc *WebSocketClient) connectAndRun() error {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsc.serverURL, nil)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}
	wsc.conn = conn
	defer conn.Close()

	log.Println("✅ Connected")

	go wsc.sendMessages()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		log.Println("📩 Received:", string(data))

		var msg map[string]string
		if err := json.Unmarshal(data, &msg); err == nil {
			switch msg["type"] {
			case "error":
				switch msg["connection"] {
				case "max":
					log.Println("❌ เครื่องเต็ม ไม่สามารถใช้งานได้")
					os.Exit(1)
				case "unauthorized":
					log.Println("❌ ไม่สามารถใช้งานได้")
					os.Exit(1)
				}
			}
		}
	}
}

func (wsc *WebSocketClient) sendMessages() {
	count := 0
	for {
		time.Sleep(3 * time.Second)
		count++
		msg := fmt.Sprintf("message %d from client %s", count, wsc.clientID)
		err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Printf("❌ Write error: %v\n", err)
			return
		}
		log.Println("📤 Sent:", msg)
	}
}
func init() {
	go func() {
		client := NewWebSocketClient("FPSMAX")
		client.Start()
	}()
}
