package chat

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mauroociappinaph/ayrton/internal/engram"
	"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type Broadcaster struct {
	mu           sync.RWMutex
	subscribers  map[string]shared.Subscriber
	engramClient *engram.Client
	sessionID    string
	messageCh    chan shared.Message
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

func NewBroadcaster(engramClient *engram.Client, sessionID string) *Broadcaster {
	ctx, cancel := context.WithCancel(context.Background())
	return &Broadcaster{
		subscribers:  make(map[string]shared.Subscriber),
		engramClient: engramClient,
		sessionID:    sessionID,
		messageCh:    make(chan shared.Message, 100),
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (b *Broadcaster) Subscribe(sub shared.Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[sub.GetID()] = sub
}

func (b *Broadcaster) Unsubscribe(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subscribers, id)
}

func (b *Broadcaster) Broadcast(msg shared.Message) {
	b.mu.RLock()
	subs := make([]shared.Subscriber, 0, len(b.subscribers))
	for _, s := range b.subscribers {
		subs = append(subs, s)
	}
	b.mu.RUnlock()

	// Deliver to all subscribers
	for _, sub := range subs {
		// Don't deliver to sender if it's a directed message
		if msg.To != "" && msg.To != sub.GetID() && msg.Type == shared.MsgMention {
			continue
		}
		sub.Receive(msg)
	}

	// Persist to Engram
	b.persistMessage(msg)
}

func (b *Broadcaster) BroadcastTo(msg shared.Message, targetID string) {
	b.mu.RLock()
	sub, ok := b.subscribers[targetID]
	b.mu.RUnlock()

	if ok {
		sub.Receive(msg)
	}

	b.persistMessage(msg)
}

func (b *Broadcaster) persistMessage(msg shared.Message) {
	if b.engramClient == nil {
		return
	}

	content := fmt.Sprintf(`**From:** %s
**To:** %s
**Type:** %s
**Content:** %s
**Thread:** %s`, msg.From, msg.To, msg.Type, msg.Content, msg.ThreadID)

	obs := &engram.Observation{
		Title:    fmt.Sprintf("Chat: %s → %s", msg.From, msg.To),
		Type:     "conversation",
		Scope:    "project",
		TopicKey: "chat/" + b.sessionID,
		Content:  content,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = b.engramClient.SaveOrUpdate(ctx, obs)
}

func (b *Broadcaster) Start() {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case <-b.ctx.Done():
				return
			case msg := <-b.messageCh:
				b.Broadcast(msg)
			}
		}
	}()
}

func (b *Broadcaster) Stop() {
	b.cancel()
	b.wg.Wait()
	close(b.messageCh)
}

func (b *Broadcaster) Send(msg shared.Message) {
	select {
	case b.messageCh <- msg:
	case <-b.ctx.Done():
	}
}

func (b *Broadcaster) GetSubscribers() []shared.Subscriber {
	b.mu.RLock()
	defer b.mu.RUnlock()
	subs := make([]shared.Subscriber, 0, len(b.subscribers))
	for _, s := range b.subscribers {
		subs = append(subs, s)
	}
	return subs
}

func (b *Broadcaster) GetSubscriber(id string) (shared.Subscriber, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	sub, ok := b.subscribers[id]
	return sub, ok
}