package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type Registry struct {
	mu          sync.RWMutex
	agents      map[string]Agent
	broadcaster shared.BroadcasterInterface
}

func NewRegistry(broadcaster shared.BroadcasterInterface) *Registry {
	return &Registry{
		agents:      make(map[string]Agent),
		broadcaster: broadcaster,
	}
}

func (r *Registry) Register(a Agent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[a.GetID()] = a
	if r.broadcaster != nil {
		r.broadcaster.Subscribe(&agentSubscriber{agent: a})
	}
}

func (r *Registry) Unregister(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.agents, id)
}

func (r *Registry) Get(id string) (Agent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.agents[id]
	return a, ok
}

func (r *Registry) List() []shared.AgentInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]shared.AgentInfo, 0, len(r.agents))
	for _, a := range r.agents {
		list = append(list, shared.AgentInfo{
			ID:     a.GetID(),
			Name:   a.GetName(),
			Role:   a.GetRole(),
			Status: shared.AgentStatus(a.GetStatus()),
		})
	}
	return list
}

func (r *Registry) Spawn(role, name string) (Agent, error) {
	a, ok := CreateAgent(role, name, r.broadcaster)
	if !ok {
		return nil, fmt.Errorf("rol desconocido: %s", role)
	}

	r.Register(a)
	ctx := context.Background()
	if err := a.Start(ctx); err != nil {
		r.Unregister(a.GetID())
		return nil, err
	}
	return a, nil
}

func (r *Registry) Kill(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, a := range r.agents {
		if a.GetName() == name {
			a.Stop()
			delete(r.agents, id)
			return nil
		}
	}
	return fmt.Errorf("agente no encontrado: %s", name)
}

func (r *Registry) Broadcast(msg shared.Message) {
	if r.broadcaster != nil {
		r.broadcaster.Send(msg)
	}
}

type agentSubscriber struct {
	agent Agent
}

func (s *agentSubscriber) Receive(msg shared.Message) {
	ctx := context.Background()
	_ = s.agent.Receive(ctx, Message{
		ID:        msg.ID,
		From:      msg.From,
		To:        msg.To,
		Content:   msg.Content,
		Type:      MessageType(msg.Type),
		ThreadID:  msg.ThreadID,
		Timestamp: msg.Timestamp,
		Payload:   msg.Payload,
	})
}

func (s *agentSubscriber) GetID() string { return s.agent.GetID() }
func (s *agentSubscriber) GetRole() string { return s.agent.GetRole() }