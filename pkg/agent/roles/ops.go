package roles

import (
	"context"
	"fmt"
	"strings"
	"time"

		"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type Ops struct {
	*shared.BaseAgent
	broadcaster shared.BroadcasterInterface
	budget      float64
	allocations map[string]float64
}

func NewOps(name string, broadcaster shared.BroadcasterInterface) *Ops {
	base := shared.NewBaseAgent(fmt.Sprintf("ops-%d", time.Now().UnixNano()), name, "ops")
	return &Ops{
		BaseAgent:   base,
		broadcaster: broadcaster,
		budget:      1000.0, // $1000 inicial
		allocations: make(map[string]float64),
	}
}

func (o *Ops) Receive(ctx context.Context, msg shared.Message) []shared.Message {
	content := strings.ToLower(msg.Content)
	var responses []shared.Message

	switch {
	case strings.Contains(content, "presupuesto") || strings.Contains(content, "budget"):
		responses = append(responses, o.showBudget()...)
	case strings.Contains(content, "apruebo") || strings.Contains(content, "aprobar"):
		responses = append(responses, o.approve(msg)...)
	case strings.Contains(content, "gasto") || strings.Contains(content, "spend"):
		responses = append(responses, o.trackSpend(msg)...)
	case strings.Contains(content, "recurso"):
		responses = append(responses, o.manageResources(msg)...)
	default:
		responses = append(responses, o.defaultResponse(msg)...)
	}

	return responses
}

func (o *Ops) showBudget() []shared.Message {
	var lines []string
	totalAllocated := 0.0
	for k, v := range o.allocations {
		lines = append(lines, fmt.Sprintf("  • %s: $%.2f", k, v))
		totalAllocated += v
	}

	return []shared.Message{{
		From:    o.GetID(),
		To:      "ceo",
		Content: fmt.Sprintf("💰 **Presupuesto:**\nTotal: $%.2f\nAsignado: $%.2f\nDisponible: $%.2f\n\n%s", o.budget, totalAllocated, o.budget-totalAllocated, strings.Join(lines, "\n")),
		Type:    shared.MsgReport,
	}}
}

func (o *Ops) approve(msg shared.Message) []shared.Message {
	// Simple parsing: "apruebo $50 marketing"
	parts := strings.Fields(msg.Content)
	if len(parts) >= 3 {
		amount := 0.0
		fmt.Sscanf(parts[1], "$%f", &amount)
		target := parts[2]

		if amount <= o.budget {
			o.budget -= amount
			o.allocations[target] += amount
			return []shared.Message{{
				From:    o.GetID(),
				To:      "ceo",
				Content: fmt.Sprintf("✅ **Aprobado**: $%.2f para %s\nPresupuesto restante: $%.2f", amount, target, o.budget),
				Type:    shared.MsgApprovalResponse,
			}}
		}
		return []shared.Message{{
			From:    o.GetID(),
			To:      "ceo",
			Content: fmt.Sprintf("❌ **Sin fondos**: Solo $%.2f disponibles", o.budget),
			Type:    shared.MsgError,
		}}
	}
	return []shared.Message{{
		From:    o.GetID(),
		To:      "ceo",
		Content: "Uso: 'apruebo $50 marketing'",
		Type:    shared.MsgChat,
	}}
}

func (o *Ops) trackSpend(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    o.GetID(),
		To:      "ceo",
		Content: "📊 Tracking activado. Todos los gastos se registran automáticamente.",
		Type:    shared.MsgReport,
	}}
}

func (o *Ops) manageResources(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    o.GetID(),
		To:      msg.From,
		Content: "⚙️ Recursos disponibles: 2 Devs, 1 Marketing, 1 Prospección\n¿Necesitás reasignar alguien?",
		Type:    shared.MsgChat,
	}}
}

func (o *Ops) defaultResponse(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    o.GetID(),
		To:      msg.From,
		Content: "⚙️ Ops aquí. ¿Presupuesto, aprobación de gasto, o recursos?",
		Type:    shared.MsgChat,
	}}
}

func (o *Ops) Start(ctx context.Context) error {
	o.Status = "active"
	return nil
}

func (o *Ops) Stop() error {
	o.Status = "stopped"
	return nil
}