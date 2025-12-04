package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	socketPath := os.Getenv("DOCKER_SOCKET")
	if socketPath == "" {
		socketPath = "/var/run/docker.sock"
	}

	log.Printf("starting herbst-docker-agent for node=%q, url=%q, socket=%q",
		nodeName, herbstURL, socketPath)

	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
		Timeout: 8 * time.Second,
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for {
		if ctx.Err() != nil {
			log.Println("shutdown requested, exiting agent loop")
			return
		}

		if err := runOnce(ctx, herbstURL, token, nodeName, httpClient); err != nil {
			log.Printf("agent cycle ended with error: %v", err)
		} else {
			log.Println("agent cycle ended without explicit error")
		}

		// kleiner Backoff, bevor wir neu verbinden
		select {
		case <-ctx.Done():
			log.Println("shutdown requested during backoff, exiting")
			return
		case <-time.After(5 * time.Second):
		}
	}
}

func runOnce(ctx context.Context, herbstURL, token, nodeName string, httpClient *http.Client) error {
	// eigene Connect-Deadline
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c, _, err := websocket.Dial(dialCtx, herbstURL, nil)
	if err != nil {
		return err
	}
	defer c.Close(websocket.StatusNormalClosure, "bye")

	log.Printf("Connected to %s as node %q", herbstURL, nodeName)

	hello := proto.HelloMessage{
		Type:     "hello",
		NodeName: nodeName,
		Token:    token,
		Kind:     "docker",
	}
	if err := sendJSON(ctx, c, hello); err != nil {
		return wrapErr("failed to send hello", err)
	}
	log.Println("Hello sent")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			containers, err := fetchDockerContainers(ctx, httpClient)
			if err != nil {
				log.Printf("failed to list containers: %v", err)
				// Kein Abbruch, einfach beim nächsten Tick nochmal probieren
				continue
			}

			msg := proto.ContainersMessage{
				Type:       "containers",
				NodeName:   nodeName,
				Containers: containers,
			}

			if err := sendJSON(ctx, c, msg); err != nil {
				// typischer Fall: broken pipe / server weg / unauthorized -> runOnce beendet sich,
				// main-Loop macht Reconnect
				return wrapErr("failed to send containers", err)
			}

			log.Printf("Sent %d containers for node %q", len(containers), nodeName)
		}
	}
}

func fetchDockerContainers(ctx context.Context, client *http.Client) ([]proto.Container, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://docker/containers/json?all=true", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, wrapErr("docker API returned error status", errStatus(resp.StatusCode))
	}

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
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// eigener Timeout fürs Schreiben
	writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.Write(writeCtx, websocket.MessageText, data)
}

// kleine Helfer für nicer Logs / Fehlermeldungen
type errStatus int

func (e errStatus) Error() string {
	return http.StatusText(int(e))
}

func wrapErr(msg string, err error) error {
	if err == nil {
		return nil
	}
	return &wrappedError{msg: msg, inner: err}
}

type wrappedError struct {
	msg   string
	inner error
}

func (w *wrappedError) Error() string {
	return w.msg + ": " + w.inner.Error()
}

func (w *wrappedError) Unwrap() error {
	return w.inner
}
