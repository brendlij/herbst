package agents

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"herbst/internal/config"
	"herbst/internal/proto"

	"nhooyr.io/websocket"
)

// serverSecret is used to generate deterministic tokens for agents
// Generated once at startup, persists for the lifetime of the process
var (
	serverSecret     []byte
	serverSecretOnce sync.Once
)

func getServerSecret() []byte {
	serverSecretOnce.Do(func() {
		serverSecret = make([]byte, 32)
		if _, err := rand.Read(serverSecret); err != nil {
			log.Printf("Warning: failed to generate secure secret: %v", err)
			// Fallback to a less secure but functional secret
			serverSecret = []byte("herbst-fallback-secret-change-me")
		}
	})
	return serverSecret
}

// GenerateToken creates a deterministic token for an agent name
// The token is derived from the agent name + server secret using HMAC-SHA256
func GenerateToken(agentName string) string {
	h := hmac.New(sha256.New, getServerSecret())
	h.Write([]byte(agentName))
	return hex.EncodeToString(h.Sum(nil))
}

type Server struct {
	reg     *Registry
	allowed map[string]string // nodeName -> token
}

func NewServer(cfg *config.Config, reg *Registry) *Server {
	allowed := make(map[string]string)
	for _, a := range cfg.Docker.Agents {
		if a.Token != "" {
			// Use configured token if provided
			allowed[a.Name] = a.Token
		} else {
			// Generate token from agent name + server secret
			allowed[a.Name] = GenerateToken(a.Name)
		}
	}

	return &Server{
		reg:     reg,
		allowed: allowed,
	}
}

func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("ws accept:", err)
		return
	}
	defer c.Close(websocket.StatusNormalClosure, "bye")

	// eigener Context, nicht r.Context()
	ctx := context.Background()

	// ---- HELLO lesen ----
	_, data, err := c.Read(ctx)
	if err != nil {
		log.Println("ws read hello:", err)
		return
	}

	var hello proto.HelloMessage
	if err := json.Unmarshal(data, &hello); err != nil {
		log.Println("invalid hello:", err)
		return
	}

	// Token-Check
	expectedToken, ok := s.allowed[hello.NodeName]
	if !ok || expectedToken != hello.Token {
		log.Printf("unauthorized agent: %s\n", hello.NodeName)
		c.Close(websocket.StatusPolicyViolation, "unauthorized")
		return
	}

	log.Printf("Agent connected: %s (kind=%s)\n", hello.NodeName, hello.Kind)

	// ---- Message-Loop ----
	for {
		_, msg, err := c.Read(ctx)
		if err != nil {
			log.Printf("ws read for %s: %v\n", hello.NodeName, err)
			return
		}

		var base struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(msg, &base); err != nil {
			log.Println("invalid message:", err)
			continue
		}

		switch base.Type {
		case "containers":
			var cm proto.ContainersMessage
			if err := json.Unmarshal(msg, &cm); err != nil {
				log.Println("invalid containers msg:", err)
				continue
			}
			s.reg.UpdateContainers(hello.NodeName, hello.Kind, cm.Containers)
		default:
			// später: metrics, logs, etc.
		}
	}
}
