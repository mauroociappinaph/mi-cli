package agent

import (
	"context"

	"github.com/mauroociappinaph/ayrton/pkg/shared"
)

// Message extends shared.Message with helper methods
type Message = shared.Message

// Agent interface
type Agent interface {
	GetID() string
	GetName() string
	GetRole() string
	GetStatus() string
	Receive(ctx context.Context, msg Message) []Message
	Start(ctx context.Context) error
	Stop() error
}

// BaseAgent embeds shared.BaseAgent
type BaseAgent = shared.BaseAgent

// NewBaseAgent creates a new base agent
var NewBaseAgent = shared.NewBaseAgent