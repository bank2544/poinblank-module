package dumper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"nhooyr.io/websocket"
)

type WebSocketClient struct {
	clientID   string
	failed     int
	maxRetries int
	serverURL  string
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewWebSocketClient(id string) *WebSocketClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketClient{
		clientID:   id,
		maxRetries: 5,
		serverURL:  fmt.Sprintf("ws://146.190.82.119:9000/?id=%s", id),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (wsc *WebSocketClient) Start() {
	for {
		err := wsc.connectAndRun()
		if err != nil {
			log.Println("❌ Error:", err)
			wsc.failed++
			if wsc.failed >= wsc.maxRetries {
				log.Println("❌ ล้มเหลวเกิน 5 ครั้ง โปรแกรมจะปิด")
				os.Exit(1)
			}
			log.Printf("🔁 Retry in 5 seconds... (%d/%d)\n", wsc.failed, wsc.maxRetries)
			time.Sleep(5 * time.Second)
		} else {
			// สำเร็จ รีเซ็ตตัวนับ
			wsc.failed = 0
		}
	}
}

func (wsc *WebSocketClient) connectAndRun() error {
	ctx, cancel := context.WithTimeout(wsc.ctx, 10*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsc.serverURL, nil)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	log.Println("✅ Connected")
	go wsc.sendMessages(conn)

	for {
		_, data, err := conn.Read(wsc.ctx)
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		log.Println("📩 Received:", string(data))

		var msg map[string]string
		if err := json.Unmarshal(data, &msg); err == nil {
			if msg["type"] == "error" {
				if msg["connection"] == "max" {
					log.Println("❌ เครื่องเต็ม ไม่สามารถใช้งานได้")
					os.Exit(1)
				}
				if msg["connection"] == "unauthorized" {
					log.Println("❌ ไม่สามารถใช้งานได้")
					os.Exit(1)
				}
			}
		}
	}
}

func (wsc *WebSocketClient) sendMessages(conn *websocket.Conn) {
	count := 0
	for {
		select {
		case <-wsc.ctx.Done():
			return
		default:
			if conn.CloseRead(wsc.ctx) != nil {
				return
			}
			count++
			msg := fmt.Sprintf("message %d from client %s", count, wsc.clientID)
			err := conn.Write(wsc.ctx, websocket.MessageText, []byte(msg))
			if err != nil {
				log.Println("❌ Write error:", err)
				return
			}
			log.Println("📤 Sent:", msg)
			time.Sleep(3 * time.Second)
		}
	}
}
