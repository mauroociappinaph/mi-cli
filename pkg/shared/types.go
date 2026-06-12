package shared

import (
	"time"
)

// MessageType defines the type of message
type MessageType int

const (
	MsgChat MessageType = iota
	MsgCommand
	MsgMention
	MsgBroadcast
	MsgSystem
	MsgApprovalRequest
	MsgApprovalResponse
	MsgReport
	MsgError
	MsgInteractionCompleted
	MsgPatternDetected
	MsgStrategyReview
	MsgMemoryQuery
)

// Message represents a chat message
type Message struct {
	ID        string
	From      string
	To        string
	Content   string
	Type      MessageType
	ThreadID  string
	Timestamp time.Time
	Payload   map[string]interface{}
}

// Subscriber defines the interface for message subscribers
type Subscriber interface {
	Receive(msg Message)
	GetID() string
	GetRole() string
}

// BroadcasterInterface defines the interface for message broadcasting
type BroadcasterInterface interface {
	Subscribe(sub Subscriber)
	Unsubscribe(id string)
	Broadcast(msg Message)
	BroadcastTo(msg Message, targetID string)
	Send(msg Message)
	Start()
	Stop()
	GetSubscribers() []Subscriber
	GetSubscriber(id string) (Subscriber, bool)
}

// AgentStatus represents agent status
type AgentStatus string

const (
	AgentStatusActive   AgentStatus = "active"
	AgentStatusIdle     AgentStatus = "idle"
	AgentStatusStopped  AgentStatus = "stopped"
	AgentStatusError    AgentStatus = "error"
)

// AgentInfo basic agent information
type AgentInfo struct {
	ID     string
	Name   string
	Role   string
	Status AgentStatus
}

// BaseAgent provides common agent functionality
type BaseAgent struct {
	ID     string
	Name   string
	Role   string
	Status string
}

func (b *BaseAgent) GetID() string   { return b.ID }
func (b *BaseAgent) GetName() string { return b.Name }
func (b *BaseAgent) GetRole() string { return b.Role }
func (b *BaseAgent) GetStatus() string { return b.Status }

func NewBaseAgent(id, name, role string) *BaseAgent {
	return &BaseAgent{
		ID:     id,
		Name:   name,
		Role:   role,
		Status: "active",
	}
}