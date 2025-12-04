package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"herbst/internal/proto"

	"nhooyr.io/websocket"
)

func main() {
	herbstURL := os.Getenv("HERBST_URL")
	token := os.Getenv("HERBST_TOKEN")
	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		nodeName = "fake-node"
	}

	if herbstURL == "" {
		log.Fatal("HERBST_URL is required")
	}
	if token == "" {
		log.Fatal("HERBST_TOKEN is required")
	}

	ctx := context.Background()

	// WebSocket verbinden
	c, _, err := websocket.Dial(ctx, herbstURL, nil)
	if err != nil {
		log.Fatalf("failed to connect to herbst: %v", err)
	}
	defer c.Close(websocket.StatusNormalClosure, "bye")

	log.Printf("Connected to %s as node %q", herbstURL, nodeName)

	// ---- HELLO senden ----
	hello := proto.HelloMessage{
		Type:     "hello",
		NodeName: nodeName,
		Token:    token,
		Kind:     "docker", // fürs Routing im Server
	}
	if err := sendJSON(ctx, c, hello); err != nil {
		log.Fatalf("failed to send hello: %v", err)
	}
	log.Println("Hello sent")

	// ---- Fake-Container-Loop ----
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		containers := []proto.Container{
			{
				ID:      "fake1234567890",
				Name:    "fake-nginx",
				Image:   "nginx:alpine",
				State:   "running",
				Status:  "Up 5 minutes",
				Created: time.Now().Add(-5 * time.Minute).Unix(),
			},
			{
				ID:      "fakeabcdef1234",
				Name:    "fake-redis",
				Image:   "redis:alpine",
				State:   "running",
				Status:  "Up 2 minutes",
				Created: time.Now().Add(-2 * time.Minute).Unix(),
			},
		}

		msg := proto.ContainersMessage{
			Type:       "containers",
			NodeName:   nodeName,
			Containers: containers,
		}

		if err := sendJSON(ctx, c, msg); err != nil {
			log.Printf("failed to send containers: %v", err)
			return
		}

		log.Printf("Sent %d fake containers for node %q", len(containers), nodeName)
	}
}

func sendJSON(ctx context.Context, c *websocket.Conn, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.Write(ctx, websocket.MessageText, b)
}
