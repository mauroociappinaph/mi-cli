package roles

import (
	"context"
	"fmt"
	"strings"
	"time"

		"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type PM struct {
	*shared.BaseAgent
	broadcaster shared.BroadcasterInterface
	specs       map[string]Spec
	specCounter int
}

type Spec struct {
	ID          string
	Title       string
	Description string
	Status      string
	AssignedTo  string
	CreatedAt   time.Time
}

func NewPM(name string, broadcaster shared.BroadcasterInterface) *PM {
	base := shared.NewBaseAgent(fmt.Sprintf("pm-%d", time.Now().UnixNano()), name, "pm")
	return &PM{
		BaseAgent:   base,
		broadcaster: broadcaster,
		specs:       make(map[string]Spec),
		specCounter: 0,
	}
}

func (p *PM) Receive(ctx context.Context, msg shared.Message) []shared.Message {
	content := strings.ToLower(msg.Content)
	var responses []shared.Message

	switch {
	case strings.Contains(content, "spec") || strings.Contains(content, "espec"):
		responses = append(responses, p.createSpec(msg)...)
	case strings.Contains(content, "backlog"):
		responses = append(responses, p.showBacklog()...)
	case strings.Contains(content, "prioriz"):
		responses = append(responses, p.prioritize(msg)...)
	case strings.Contains(content, "tarea") || strings.Contains(content, "task"):
		responses = append(responses, p.createTask(msg)...)
	default:
		responses = append(responses, p.defaultResponse(msg)...)
	}

	return responses
}

func (p *PM) createSpec(msg shared.Message) []shared.Message {
	p.specCounter++
	specID := fmt.Sprintf("SPEC-%03d", p.specCounter)
	spec := Spec{
		ID:          specID,
		Title:       fmt.Sprintf("Spec desde: %s", msg.Content),
		Description: msg.Content,
		Status:      "draft",
		CreatedAt:   time.Now(),
	}
	p.specs[specID] = spec

	return []shared.Message{{
		From:    p.GetID(),
		To:      "ceo",
		Content: fmt.Sprintf("📝 **Spec %s creada**\nTítulo: %s\nEstado: %s\n\n¿Asigno a Dev para implementar?", specID, spec.Title, spec.Status),
		Type:    shared.MsgReport,
	}}
}

func (p *PM) showBacklog() []shared.Message {
	if len(p.specs) == 0 {
		return []shared.Message{{
			From:    p.GetID(),
			To:      "ceo",
			Content: "📋 Backlog vacío. No hay specs pendientes.",
			Type:    shared.MsgReport,
		}}
	}

	var lines []string
	for _, s := range p.specs {
		lines = append(lines, fmt.Sprintf("  • %s: %s [%s]", s.ID, s.Title, s.Status))
	}

	return []shared.Message{{
		From:    p.GetID(),
		To:      "ceo",
		Content: fmt.Sprintf("📋 **Backlog (%d specs):**\n%s", len(p.specs), strings.Join(lines, "\n")),
		Type:    shared.MsgReport,
	}}
}

func (p *PM) prioritize(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    p.GetID(),
		To:      "ceo",
		Content: "📊 Priorización aplicada: specs con mayor ROI y menor esfuerzo primero. ¿Confirmas orden?",
		Type:    shared.MsgReport,
	}}
}

func (p *PM) createTask(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    p.GetID(),
		To:      "dev",
		Content: fmt.Sprintf("🔧 Nueva tarea para Dev: %s", msg.Content),
		Type:    shared.MsgMention,
	}}
}

func (p *PM) defaultResponse(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    p.GetID(),
		To:      msg.From,
		Content: "📋 PM aquí. ¿Necesitás una spec, priorizar backlog, o crear tarea para Dev?",
		Type:    shared.MsgChat,
	}}
}

func (p *PM) Start(ctx context.Context) error {
	p.Status = "active"
	return nil
}

func (p *PM) Stop() error {
	p.Status = "stopped"
	return nil
}