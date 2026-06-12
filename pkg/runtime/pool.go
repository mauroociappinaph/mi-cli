package runtime

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mauroociappinaph/ayrton/pkg/agent"
	"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type Pool struct {
	mu          sync.RWMutex
	workers     map[string]*worker
	broadcaster shared.BroadcasterInterface
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

type worker struct {
	agent   agent.Agent
	ctx     context.Context
	cancel  context.CancelFunc
	inbox   chan shared.Message
	started time.Time
}

func NewPool(broadcaster shared.BroadcasterInterface) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		workers:      make(map[string]*worker),
		broadcaster:  broadcaster,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (p *Pool) Spawn(a agent.Agent) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	id := a.GetID()
	if _, exists := p.workers[id]; exists {
		return fmt.Errorf("worker %s already exists", id)
	}

	ctx, cancel := context.WithCancel(p.ctx)
	inbox := make(chan shared.Message, 50)

	w := &worker{
		agent:   a,
		ctx:     ctx,
		cancel:  cancel,
		inbox:   inbox,
		started: time.Now(),
	}

	p.workers[id] = w

	p.wg.Add(1)
	go p.runWorker(w)

	return nil
}

func (p *Pool) runWorker(w *worker) {
	defer p.wg.Done()

	for {
		select {
		case <-w.ctx.Done():
			return
		case msg := <-w.inbox:
			responses := w.agent.Receive(w.ctx, agent.Message{
				ID:        msg.ID,
				From:      msg.From,
				To:        msg.To,
				Content:   msg.Content,
				Type:      agent.MessageType(msg.Type),
				ThreadID:  msg.ThreadID,
				Timestamp: msg.Timestamp,
				Payload:   msg.Payload,
			})

			for _, resp := range responses {
				outMsg := shared.Message{
					ID:        fmt.Sprintf("msg-%d", time.Now().UnixNano()),
					From:      w.agent.GetID(),
					To:        resp.To,
					Content:   resp.Content,
					Type:      shared.MessageType(resp.Type),
					ThreadID:  resp.ThreadID,
					Timestamp: time.Now(),
				}
				if p.broadcaster != nil {
					p.broadcaster.Send(outMsg)
				}
			}
		}
	}
}

func (p *Pool) Send(msg shared.Message) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if msg.To == "" {
		// Broadcast to all
		for _, w := range p.workers {
			select {
			case w.inbox <- msg:
			default:
				// Inbox full, skip
			}
		}
	} else {
		// Direct message
		if w, ok := p.workers[msg.To]; ok {
			select {
			case w.inbox <- msg:
			default:
				// Inbox full
			}
		}
	}
}

func (p *Pool) Kill(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if w, ok := p.workers[id]; ok {
		w.cancel()
		delete(p.workers, id)
		return nil
	}
	return fmt.Errorf("worker not found: %s", id)
}

func (p *Pool) List() []WorkerInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()

	list := make([]WorkerInfo, 0, len(p.workers))
	for _, w := range p.workers {
		list = append(list, WorkerInfo{
			ID:        w.agent.GetID(),
			Name:      w.agent.GetName(),
			Role:      w.agent.GetRole(),
			Status:    w.agent.GetStatus(),
			Uptime:    time.Since(w.started).Round(time.Second).String(),
		})
	}
	return list
}

func (p *Pool) Status() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return fmt.Sprintf("%d workers active", len(p.workers))
}

func (p *Pool) Start() {
	// Pool starts automatically when workers are spawned
}

func (p *Pool) Stop() {
	p.cancel()
	p.wg.Wait()

	p.mu.Lock()
	defer p.mu.Unlock()
	for _, w := range p.workers {
		w.cancel()
	}
	p.workers = make(map[string]*worker)
}

type WorkerInfo struct {
	ID     string
	Name   string
	Role   string
	Status string
	Uptime string
}