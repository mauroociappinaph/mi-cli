package agent

// MockHandlers proporciona handlers falsos para testing FASE 0 sin LLM real
// En FASE 1 se reemplazarán con integración real de Anthropic/OpenAI SDK

import (
	"context"
	"fmt"
	"strings"

	"github.com/mauroociappinaph/ayrton/pkg/shared"
)

// MockHandler es una función que simula respuesta de agente
type MockHandler func(ctx context.Context, msg Message) []Message

// PMHandlers handlers mock para Product Manager
var PMHandlers = map[string]MockHandler{
	"create_spec": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "pm",
			To:    "ceo",
			Content: "📝 Spec creada: Landing page SaaS B2B - Hero, Pain points, Demo, Pricing, CTA. Estimado: 4h dev.",
			Type:  shared.MsgReport,
		}}
	},
	"prioritize": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "pm",
			To:    "ceo",
			Content: "📊 Priorización: 1) Landing SaaS (ROI alto, esfuerzo bajo) 2) Dashboard analytics 3) Mobile app",
			Type:  shared.MsgReport,
		}}
	},
	"breakdown": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "pm",
			To:    "dev",
			Content: "🔧 Tareas: 1) Setup repo 2) Hero section 3) Form capture 4) Stripe integration 5) Tests 6) Deploy",
			Type:  shared.MsgMention,
		}}
	},
}

// DevHandlers handlers mock para Developer
var DevHandlers = map[string]MockHandler{
	"implement": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "dev",
			To:    "ceo",
			Content: "⚡ Implementando... Te aviso en 2h. Setup completado, empezando Hero section.",
			Type:  shared.MsgReport,
		}}
	},
	"test": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "dev",
			To:    "ceo",
			Content: "🧪 Tests pasando: `go test -v -race ./...`\nCoverage: 87% | Race: clean",
			Type:  shared.MsgReport,
		}}
	},
	"deploy": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "dev",
			To:    "ceo",
			Content: "🚀 Deploy staging OK: https://staging.ayrton.app\nListo para validación.",
			Type:  shared.MsgReport,
		}}
	},
	"write_code": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "dev",
			To:    msg.From,
			Content: "💻 Generando código (mock FASE 0)... En FASE 1 usaré LLM real.",
			Type:  shared.MsgChat,
		}}
	},
}

// MarketingHandlers handlers mock para Marketing
var MarketingHandlers = map[string]MockHandler{
	"campaign": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "marketing",
			To:    "ceo",
			Content: "📊 Campaña: Landing test A/B, $50/día, 7 días. Target: Founders SaaS Latam. ¿Aprobo?",
			Type:  shared.shared.MsgApprovalRequest,
		}}
	},
	"landing": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "marketing",
			To:    "pm",
			Content: "🎯 3 variantes: V1 (pain-point), V2 (benefit), V3 (social proof). Necesito spec de PM.",
			Type:  shared.MsgReport,
		}}
	},
	"ads": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "marketing",
			To:    "ceo",
			Content: "📱 Ads: Google Search + Meta + LinkedIn. $50/día test A/B 7 días. Target: Founders SaaS B2B Latam. ¿Aprobo?",
			Type:  shared.shared.MsgApprovalRequest,
		}}
	},
	"content": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "marketing",
			To:    msg.From,
			Content: "✍️ Content calendar: 3 posts/sem (blog, LinkedIn, Twitter). Temas: SaaS metrics, founder journey, product updates.",
			Type:  shared.MsgChat,
		}}
	},
}

// OpsHandlers handlers mock para Operations
var OpsHandlers = map[string]MockHandler{
	"budget": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "ops",
			To:    "ceo",
			Content: "💰 Presupuesto: $1000 total | Asignado: $150 | Disponible: $850\n• Marketing: $100\n• Dev tools: $50",
			Type:  shared.MsgReport,
		}}
	},
	"approve": func(ctx context.Context, msg Message) []Message {
		// Simple parsing
		content := msg.Content
		if strings.Contains(content, "apruebo") {
			return []Message{{
				From:  "ops",
				To:    "ceo",
				Content: "✅ Aprobado: $50 para marketing. Restante: $800",
				Type:  shared.MsgApprovalResponse,
			}}
		}
		return []Message{{
			From:  "ops",
			To:    "ceo",
			Content: "Uso: 'apruebo $50 marketing'",
			Type:  shared.MsgChat,
		}}
	},
	"resources": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "ops",
			To:    msg.From,
			Content: "⚙️ Recursos: 2 Devs, 1 Marketing, 1 Prospección. ¿Reasignar?",
			Type:  shared.MsgChat,
		}}
	},
}

// ProspeccionHandlers handlers mock para Prospección
var ProspeccionHandlers = map[string]MockHandler{
	"leads": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "prospeccion",
			To:    "ceo",
			Content: "🔍 3 leads SaaS B2B Latam:\n1. Carlos Mendoza - FacturaYa (CTO) - Score 85\n2. Ana Torres - HRTech Latam (Founder) - Score 78\n3. Roberto Silva - LogiChain (VP Ops) - Score 92",
			Type:  shared.MsgReport,
		}}
	},
	"outreach": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "prospeccion",
			To:    "ceo",
			Content: "📧 Outreach: 20 emails, 15 LinkedIn, 5 calls agendadas. Tracking replies en pipeline.",
			Type:  shared.MsgReport,
		}}
	},
	"research": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "prospeccion",
			To:    msg.From,
			Content: "🔬 Research: Top 3 nichos - Facturación electrónica AR, HR Tech PyMEs, Logística last-mile. TAM $2.4B. Dolor: compliance local + integraciones.",
			Type:  shared.MsgReport,
		}}
	},
	"pipeline": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "prospeccion",
			To:    "ceo",
			Content: "📊 Pipeline: 3 new, 2 contacted, 1 qualified, 0 closed",
			Type:  shared.MsgReport,
		}}
	},
}

// AuditorHandlers handlers mock para Auditor
var AuditorHandlers = map[string]MockHandler{
	"audit": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "auditor",
			To:    "ceo",
			Content: "🔍 Auditoría:\n✅ Code quality: golangci-lint clean\n✅ Test coverage: 87%\n⚠️ Perf: p95 2.3s (target <1.5s)\n✅ Security: 0 vulns\n⚠️ Arch: pkg/memory/conversation.go - 45 lines (max 30)",
			Type:  shared.MsgReport,
		}}
	},
	"security": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "auditor",
			To:    "ceo",
			Content: "🔒 Security: govulncheck 0 vulns, gosec clean, secrets scan clean, deps up to date",
			Type:  shared.MsgReport,
		}}
	},
	"performance": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "auditor",
			To:    "ceo",
			Content: "⚡ Perf: p50 45ms, p95 2.3s ⚠️, p99 4.1s, Mem 42MB, Goroutines 12\n💡 Revisar query N+1 en /api/leads",
			Type:  shared.MsgReport,
		}}
	},
}

// LearningHandlers handlers mock para Learning
var LearningHandlers = map[string]MockHandler{
	"learn": func(ctx context.Context, msg Message) []Message {
		trigger := msg.GetString("trigger")
		action := msg.GetString("action")
		outcome := msg.GetString("outcome")
		if trigger == "" || action == "" || outcome == "" {
			return []Message{{
				From:  "learning",
				To:    msg.From,
				Content: "Uso: trigger=\"...\" action=\"...\" outcome=\"...\"",
				Type:  shared.MsgChat,
			}}
		}
		return []Message{{
			From:  "learning",
			To:    msg.From,
			Content: fmt.Sprintf("🧠 Lesson: %s → %s [confidence: 0.50]", trigger, outcome),
			Type:  shared.MsgReport,
		}}
	},
	"pattern": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "learning",
			To:    msg.From,
			Content: "📚 Patrones: deploy staging → success (freq: 12, conf: 0.73), spec review → 0 bugs (freq: 8, conf: 0.68)",
			Type:  shared.MsgReport,
		}}
	},
	"recall": func(ctx context.Context, msg Message) []Message {
		return []Message{{
			From:  "learning",
			To:    msg.From,
			Content: "🔍 Recall: 'deploy staging' → 3 lessons, 'spec review' → 2 lessons, 'error handling' → 1 lesson",
			Type:  shared.MsgReport,
		}}
	},
}

// GetHandler busca handler por acción
func GetHandler(role, action string) MockHandler {
	switch role {
	case "pm":
		return PMHandlers[action]
	case "dev":
		return DevHandlers[action]
	case "marketing":
		return MarketingHandlers[action]
	case "ops":
		return OpsHandlers[action]
	case "prospeccion":
		return ProspeccionHandlers[action]
	case "auditor":
		return AuditorHandlers[action]
	case "learning":
		return LearningHandlers[action]
	}
	return nil
}

// ExecuteHandler ejecuta handler si existe
func ExecuteHandler(ctx context.Context, role, action string, msg Message) []Message {
	if h := GetHandler(role, action); h != nil {
		return h(ctx, msg)
	}
	// Default response
	return []Message{{
		From:  role,
		To:    msg.From,
		Content: fmt.Sprintf("%s: ¿En qué te ayudo? (acciones: %s)", role, getAvailableActions(role)),
		Type:  shared.MsgChat,
	}}
}

func getAvailableActions(role string) string {
	switch role {
	case "pm":
		return "create_spec, prioritize, breakdown"
	case "dev":
		return "implement, test, deploy, write_code"
	case "marketing":
		return "campaign, landing, ads, content"
	case "ops":
		return "budget, approve, resources"
	case "prospeccion":
		return "leads, outreach, research, pipeline"
	case "auditor":
		return "audit, security, performance"
	case "learning":
		return "learn, pattern, recall"
	}
	return ""
}