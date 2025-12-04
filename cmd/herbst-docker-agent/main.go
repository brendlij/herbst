package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"herbst/internal/proto"

	"nhooyr.io/websocket"
)

func main() {
	herbstURL := os.Getenv("HERBST_URL")
	token := os.Getenv("HERBST_TOKEN")
	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		nodeName = "docker-node"
	}

	if herbstURL == "" {
		log.Fatal("HERBST_URL is required")
	}
	if token == "" {
		log.Fatal("HERBST_TOKEN is required")
	}

	// Docker socket path (default für Linux)
	socketPath := os.Getenv("DOCKER_SOCKET")
	if socketPath == "" {
		socketPath = "/var/run/docker.sock"
	}

	ctx := context.Background()

	// HTTP-Client, der über den Unix-Socket spricht
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
		Timeout: 5 * time.Second,
	}

	// WebSocket zu herbst
	c, _, err := websocket.Dial(ctx, herbstURL, nil)
	if err != nil {
		log.Fatalf("failed to connect to herbst: %v", err)
	}
	defer c.Close(websocket.StatusNormalClosure, "bye")

	log.Printf("Connected to %s as node %q", herbstURL, nodeName)

	// HELLO
	hello := proto.HelloMessage{
		Type:     "hello",
		NodeName: nodeName,
		Token:    token,
		Kind:     "docker",
	}
	if err := sendJSON(ctx, c, hello); err != nil {
		log.Fatalf("failed to send hello: %v", err)
	}
	log.Println("Hello sent")

	// Polling-Loop
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		containers, err := fetchDockerContainers(ctx, httpClient)
		if err != nil {
			log.Printf("failed to list containers: %v", err)
			continue
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

		log.Printf("Sent %d containers for node %q", len(containers), nodeName)
	}
}

func fetchDockerContainers(ctx context.Context, client *http.Client) ([]proto.Container, error) {
	// Docker HTTP API – identisch zu deinem Server-Code
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://docker/containers/json?all=true", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw []struct {
		ID      string   `json:"Id"`
		Names   []string `json:"Names"`
		Image   string   `json:"Image"`
		State   string   `json:"State"`
		Status  string   `json:"Status"`
		Created int64    `json:"Created"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	out := make([]proto.Container, 0, len(raw))
	for _, c := range raw {
		name := c.ID
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		} else if len(c.ID) >= 12 {
			name = c.ID[:12]
		}

		out = append(out, proto.Container{
			ID:      c.ID,
			Name:    name,
			Image:   c.Image,
			State:   c.State,
			Status:  c.Status,
			Created: c.Created,
		})
	}

	return out, nil
}

func sendJSON(ctx context.Context, c *websocket.Conn, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.Write(ctx, websocket.MessageText, b)
}
