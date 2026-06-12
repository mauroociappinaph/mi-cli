package roles

import (
	"context"
	"fmt"
	"strings"
	"time"

		"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type Dev struct {
	*shared.BaseAgent
	broadcaster shared.BroadcasterInterface
	tasks       map[string]Task
	taskCounter int
}

type Task struct {
	ID          string
	Title       string
	Description string
	Status      string
	AssignedAt  time.Time
	CompletedAt *time.Time
}

func NewDev(name string, broadcaster shared.BroadcasterInterface) *Dev {
	base := shared.NewBaseAgent(fmt.Sprintf("dev-%d", time.Now().UnixNano()), name, "dev")
	return &Dev{
		BaseAgent:   base,
		broadcaster: broadcaster,
		tasks:       make(map[string]Task),
		taskCounter: 0,
	}
}

func (d *Dev) Receive(ctx context.Context, msg shared.Message) []shared.Message {
	content := strings.ToLower(msg.Content)
	var responses []shared.Message

	switch {
	case strings.Contains(content, "implementa") || strings.Contains(content, "implement"):
		responses = append(responses, d.startImplementation(msg)...)
	case strings.Contains(content, "test"):
		responses = append(responses, d.runTests(msg)...)
	case strings.Contains(content, "deploy"):
		responses = append(responses, d.deploy(msg)...)
	case strings.Contains(content, "código") || strings.Contains(content, "code"):
		responses = append(responses, d.writeCode(msg)...)
	default:
		responses = append(responses, d.defaultResponse(msg)...)
	}

	return responses
}

func (d *Dev) startImplementation(msg shared.Message) []shared.Message {
	d.taskCounter++
	taskID := fmt.Sprintf("TASK-%03d", d.taskCounter)
	task := Task{
		ID:          taskID,
		Title:       fmt.Sprintf("Implement: %s", msg.Content),
		Description: msg.Content,
		Status:      "in_progress",
		AssignedAt:  time.Now(),
	}
	d.tasks[taskID] = task

	return []shared.Message{{
		From:    d.GetID(),
		To:      "ceo",
		Content: fmt.Sprintf("⚡ **Iniciando implementación %s**\n%s\n\nTe aviso cuando esté listo para review.", taskID, msg.Content),
		Type:    shared.MsgReport,
	}}
}

func (d *Dev) runTests(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    d.GetID(),
		To:      "ceo",
		Content: "🧪 Ejecutando tests...\n`go test -v -race ./...`\n\n✅ Tests pasando. Coverage: 87%",
		Type:    shared.MsgReport,
	}}
}

func (d *Dev) deploy(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    d.GetID(),
		To:      "ceo",
		Content: "🚀 Deploy a staging completado.\nURL: https://staging.ayrton.app\n\nListo para validación.",
		Type:    shared.MsgReport,
	}}
}

func (d *Dev) writeCode(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    d.GetID(),
		To:      msg.From,
		Content: "💻 Escribiendo código... (mock handler FASE 0)\nEn FASE 1 integraré LLM real para generación de código.",
		Type:    shared.MsgChat,
	}}
}

func (d *Dev) defaultResponse(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    d.GetID(),
		To:      msg.From,
		Content: "⚡ Dev aquí. ¿Implemento algo, corro tests, o deploy a staging?",
		Type:    shared.MsgChat,
	}}
}

func (d *Dev) Start(ctx context.Context) error {
	d.Status = "active"
	return nil
}

func (d *Dev) Stop() error {
	d.Status = "stopped"
	return nil
}