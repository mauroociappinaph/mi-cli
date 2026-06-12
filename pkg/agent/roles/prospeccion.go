package roles

import (
	"context"
	"fmt"
	"strings"
	"time"

		"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type Prospeccion struct {
	*shared.BaseAgent
	broadcaster shared.BroadcasterInterface
	leads       map[string]Lead
	leadCounter int
}

type Lead struct {
	ID          string
	Name        string
	Company     string
	Role        string
	Source      string
	Status      string
	Score       int
	CreatedAt   time.Time
}

func NewProspeccion(name string, broadcaster shared.BroadcasterInterface) *Prospeccion {
	base := shared.NewBaseAgent(fmt.Sprintf("prospeccion-%d", time.Now().UnixNano()), name, "prospeccion")
	return &Prospeccion{
		BaseAgent:   base,
		broadcaster: broadcaster,
		leads:       make(map[string]Lead),
		leadCounter: 0,
	}
}

func (p *Prospeccion) Receive(ctx context.Context, msg shared.Message) []shared.Message {
	content := strings.ToLower(msg.Content)
	var responses []shared.Message

	switch {
	case strings.Contains(content, "lead") || strings.Contains(content, "prospect"):
		responses = append(responses, p.findLeads(msg)...)
	case strings.Contains(content, "outreach") || strings.Contains(content, "contact"):
		responses = append(responses, p.outreach(msg)...)
	case strings.Contains(content, "research") || strings.Contains(content, "investig"):
		responses = append(responses, p.research(msg)...)
	case strings.Contains(content, "pipeline"):
		responses = append(responses, p.showPipeline()...)
	default:
		responses = append(responses, p.defaultResponse(msg)...)
	}

	return responses
}

func (p *Prospeccion) findLeads(msg shared.Message) []shared.Message {
	p.leadCounter++
	leadID := fmt.Sprintf("LEAD-%03d", p.leadCounter)

	// Mock leads based on query
	leads := []Lead{
		{ID: fmt.Sprintf("%s-1", leadID), Name: "Carlos Mendoza", Company: "FacturaYa", Role: "CTO", Source: "LinkedIn", Status: "new", Score: 85},
		{ID: fmt.Sprintf("%s-2", leadID), Name: "Ana Torres", Company: "HRTech Latam", Role: "Founder", Source: "Twitter", Status: "new", Score: 78},
		{ID: fmt.Sprintf("%s-3", leadID), Name: "Roberto Silva", Company: "LogiChain", Role: "VP Ops", Source: "Referral", Status: "contacted", Score: 92},
	}

	for _, l := range leads {
		p.leads[l.ID] = l
	}

	var lines []string
	for _, l := range leads {
		lines = append(lines, fmt.Sprintf("  • %s - %s (%s) | Score: %d | %s", l.Name, l.Company, l.Role, l.Score, l.Source))
	}

	return []shared.Message{{
		From:    p.GetID(),
		To:      "ceo",
		Content: fmt.Sprintf("🔍 **%d leads encontrados** para: %s\n%s", len(leads), msg.Content, strings.Join(lines, "\n")),
		Type:    shared.MsgReport,
	}}
}

func (p *Prospeccion) outreach(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    p.GetID(),
		To:      "ceo",
		Content: "📧 Outreach iniciado:\n• 20 emails personalizados enviados\n• 15 conexiones LinkedIn solicitadas\n• 5 calls agendadas para esta semana\n\nTracking: replies en pipeline.",
		Type:    shared.MsgReport,
	}}
}

func (p *Prospeccion) research(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    p.GetID(),
		To:      msg.From,
		Content: "🔬 Research completado:\n• Top 3 nichos SaaS B2B Latam: Facturación electrónica, HR Tech, Logística last-mile\n• TAM estimado: $2.4B combinado\n• Competencia: Fragmentada, sin líder claro\n• Dolor principal: Compliance local + integraciones\n\n¿Profundizo en alguno?",
		Type:    shared.MsgReport,
	}}
}

func (p *Prospeccion) showPipeline() []shared.Message {
	if len(p.leads) == 0 {
		return []shared.Message{{
			From:    p.GetID(),
			To:      "ceo",
			Content: "📊 Pipeline vacío. Ejecuta 'busca leads' para empezar.",
			Type:    shared.MsgReport,
		}}
	}

	counts := map[string]int{}
	for _, l := range p.leads {
		counts[l.Status]++
	}

	var lines []string
	for status, count := range counts {
		lines = append(lines, fmt.Sprintf("  • %s: %d", status, count))
	}

	return []shared.Message{{
		From:    p.GetID(),
		To:      "ceo",
		Content: fmt.Sprintf("📊 **Pipeline (%d leads):**\n%s", len(p.leads), strings.Join(lines, "\n")),
		Type:    shared.MsgReport,
	}}
}

func (p *Prospeccion) defaultResponse(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    p.GetID(),
		To:      msg.From,
		Content: "🔍 Prospección aquí. ¿Busco leads, hago outreach, research, o muestro pipeline?",
		Type:    shared.MsgChat,
	}}
}

func (p *Prospeccion) Start(ctx context.Context) error {
	p.Status = "active"
	return nil
}

func (p *Prospeccion) Stop() error {
	p.Status = "stopped"
	return nil
}