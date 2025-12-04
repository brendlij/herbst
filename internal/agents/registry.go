// internal/agents/registry.go
package agents

import (
	"sync"
	"time"

	"herbst/internal/proto"
)

type NodeState struct {
	Name       string             `json:"name"`
	Kind       string             `json:"kind"`      // z.B. "docker"
	LastSeen   time.Time          `json:"lastSeen"`
	Containers []proto.Container  `json:"containers"`
}

type Registry struct {
	mu    sync.RWMutex
	nodes map[string]*NodeState
}

func NewRegistry() *Registry {
	return &Registry{
		nodes: make(map[string]*NodeState),
	}
}

func (r *Registry) UpdateContainers(nodeName string, kind string, containers []proto.Container) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ns, ok := r.nodes[nodeName]
	if !ok {
		ns = &NodeState{Name: nodeName}
		r.nodes[nodeName] = ns
	}
	ns.Kind = kind
	ns.Containers = containers
	ns.LastSeen = time.Now()
}

func (r *Registry) Snapshot() map[string]NodeState {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[string]NodeState, len(r.nodes))
	for k, v := range r.nodes {
		out[k] = *v
	}
	return out
}
