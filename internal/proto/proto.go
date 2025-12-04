// internal/proto/proto.go
package proto

type HelloMessage struct {
	Type     string `json:"type"` // "hello"
	NodeName string `json:"nodeName"`
	Token    string `json:"token"`
	Kind     string `json:"kind"` // z.B. "docker"
}

type Container struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	State   string `json:"state"`
	Status  string `json:"status"`
	Created int64  `json:"created"`
}

type ContainersMessage struct {
	Type       string      `json:"type"` // "containers"
	NodeName   string      `json:"nodeName"`
	Containers []Container `json:"containers"`
}
