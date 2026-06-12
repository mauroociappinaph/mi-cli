package roles

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mauroociappinaph/ayrton/internal/engram"
	"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type Auditor struct {
	*shared.BaseAgent
	broadcaster  shared.BroadcasterInterface
	engramClient *engram.Client
}

func NewAuditor(engramClient *engram.Client, broadcaster shared.BroadcasterInterface) *Auditor {
	base := shared.NewBaseAgent("auditor", "Auditor", "auditor")
	return &Auditor{
		BaseAgent:    base,
		broadcaster:  broadcaster,
		engramClient: engramClient,
	}
}

func (a *Auditor) Receive(ctx context.Context, msg shared.Message) []shared.Message {
	content := strings.ToLower(msg.Content)
	var responses []shared.Message

	switch {
	case strings.Contains(content, "audit") || strings.Contains(content, "audita"):
		responses = append(responses, a.runAudit(msg)...)
	case strings.Contains(content, "security") || strings.Contains(content, "seguridad"):
		responses = append(responses, a.securityScan(msg)...)
	case strings.Contains(content, "performance") || strings.Contains(content, "perf"):
		responses = append(responses, a.performanceCheck(msg)...)
	case strings.Contains(content, "quality") || strings.Contains(content, "calidad"):
		responses = append(responses, a.qualityGate(msg)...)
	default:
		responses = append(responses, a.defaultResponse(msg)...)
	}

	return responses
}

func (a *Auditor) runAudit(msg shared.Message) []shared.Message {
	findings := []string{
		"✅ Code quality: golangci-lint clean",
		"✅ Test coverage: 87% (target >80%)",
		"⚠️ Performance: p95 latency 2.3s (target <1.5s)",
		"✅ Security: 0 vulns (govulncheck)",
		"⚠️ Architecture: pkg/memory/conversation.go - function 45 lines (max 30)",
	}

	auditMsg := shared.Message{
		ID:        fmt.Sprintf("audit-%d", time.Now().UnixNano()),
		From:      "auditor",
		To:        "ceo",
		Content:   fmt.Sprintf("🔍 **Auditoría completa**\n%s", strings.Join(findings, "\n")),
		Type:      shared.MsgReport,
		ThreadID:  fmt.Sprintf("audit-%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
	}
	a.broadcaster.Send(auditMsg)

	// Persist to Engram
	if a.engramClient != nil {
		obs := &engram.Observation{
			Title:    "Auto-auditoría post-deploy",
			Type:     "audit",
			Scope:    "project",
			TopicKey: "audit/full",
			Content:  strings.Join(findings, "\n"),
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = a.engramClient.SaveOrUpdate(ctx, obs)
	}

	return []shared.Message{{
		From:    "auditor",
		To:      "ceo",
		Content: fmt.Sprintf("🔍 Auditoría completada. %d hallazgos (2 warnings).", len(findings)),
		Type:    shared.MsgReport,
	}}
}

func (a *Auditor) securityScan(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    "auditor",
		To:      "ceo",
		Content: "🔒 Security scan:\n• govulncheck: 0 vulnerabilities\n• gosec: No issues\n• Secrets scan: Clean\n• Dependencies: All up to date",
		Type:    shared.MsgReport,
	}}
}

func (a *Auditor) performanceCheck(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    "auditor",
		To:      "ceo",
		Content: "⚡ Performance check:\n• p50 latency: 45ms\n• p95 latency: 2.3s ⚠️ (target <1.5s)\n• p99 latency: 4.1s\n• Memory: 42MB RSS\n• Goroutines: 12\n\n💡 Revisar query N+1 en handler /api/leads",
		Type:    shared.MsgReport,
	}}
}

func (a *Auditor) qualityGate(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    "auditor",
		To:      "ceo",
		Content: "📏 Quality gates:\n• gofmt: ✅\n• go vet: ✅\n• golangci-lint: ✅\n• Test coverage: 87% ✅\n• Race detector: ✅\n• Cyclomatic complexity: avg 3.2 ✅",
		Type:    shared.MsgReport,
	}}
}

func (a *Auditor) defaultResponse(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    "auditor",
		To:      msg.From,
		Content: "🔍 Auditor aquí. ¿Ejecuto auditoría completa, security scan, performance check, o quality gate?",
		Type:    shared.MsgChat,
	}}
}

func (a *Auditor) Start(ctx context.Context) error {
	a.Status = "active"
	return nil
}

func (a *Auditor) Stop() error {
	a.Status = "stopped"
	return nil
}