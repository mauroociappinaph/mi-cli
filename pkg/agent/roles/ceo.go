package roles

import (
	"context"
	"fmt"
	"time"

	"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type CEO struct {
	*shared.BaseAgent
	broadcaster shared.BroadcasterInterface
}

func NewCEO(broadcaster shared.BroadcasterInterface) *CEO {
	base := shared.NewBaseAgent("ceo", "CEO", "ceo")
	return &CEO{
		BaseAgent:   base,
		broadcaster: broadcaster,
	}
}

func (c *CEO) Receive(ctx context.Context, msg shared.Message) []shared.Message {
	return nil
}

func (c *CEO) Start(ctx context.Context) error {
	c.Status = "active"
	return nil
}

func (c *CEO) Stop() error {
	c.Status = "stopped"
	return nil
}

func (c *CEO) Approve(actionID, comment string) {
	msg := shared.Message{
		ID:        fmt.Sprintf("approval-%d", time.Now().UnixNano()),
		From:      "ceo",
		To:        "system",
		Content:   fmt.Sprintf("✅ APROBADO: %s\n%s", actionID, comment),
		Type:      shared.MsgApprovalResponse,
		ThreadID:  actionID,
		Timestamp: time.Now(),
	}
	c.broadcaster.Send(msg)
}

func (c *CEO) Reject(actionID, reason string) {
	msg := shared.Message{
		ID:        fmt.Sprintf("rejection-%d", time.Now().UnixNano()),
		From:      "ceo",
		To:        "system",
		Content:   fmt.Sprintf("❌ RECHAZADO: %s\nRazón: %s", actionID, reason),
		Type:      shared.MsgApprovalResponse,
		ThreadID:  actionID,
		Timestamp: time.Now(),
	}
	c.broadcaster.Send(msg)
}

func (c *CEO) Redirect(fromAgent, toAgent, instruction string) {
	msg := shared.Message{
		ID:        fmt.Sprintf("redirect-%d", time.Now().UnixNano()),
		From:      "ceo",
		To:        toAgent,
		Content:   fmt.Sprintf("🔄 REDIRECCIÓN DEL CEO: %s", instruction),
		Type:      shared.MsgMention,
		ThreadID:  fmt.Sprintf("redirect-%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
	}
	c.broadcaster.Send(msg)
}