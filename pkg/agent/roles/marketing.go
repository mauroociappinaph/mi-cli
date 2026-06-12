package roles

import (
	"context"
	"fmt"
	"strings"
	"time"

		"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type Marketing struct {
	*shared.BaseAgent
	broadcaster shared.BroadcasterInterface
	campaigns   map[string]Campaign
	campCounter int
}

type Campaign struct {
	ID          string
	Name        string
	Channel     string
	Budget      float64
	Status      string
	CreatedAt   time.Time
}

func NewMarketing(name string, broadcaster shared.BroadcasterInterface) *Marketing {
	base := shared.NewBaseAgent(fmt.Sprintf("marketing-%d", time.Now().UnixNano()), name, "marketing")
	return &Marketing{
		BaseAgent:   base,
		broadcaster: broadcaster,
		campaigns:   make(map[string]Campaign),
		campCounter: 0,
	}
}

func (m *Marketing) Receive(ctx context.Context, msg shared.Message) []shared.Message {
	content := strings.ToLower(msg.Content)
	var responses []shared.Message

	switch {
	case strings.Contains(content, "campaña") || strings.Contains(content, "campaign"):
		responses = append(responses, m.createCampaign(msg)...)
	case strings.Contains(content, "landing"):
		responses = append(responses, m.landingPage(msg)...)
	case strings.Contains(content, "ads") || strings.Contains(content, "anuncio"):
		responses = append(responses, m.createAds(msg)...)
	case strings.Contains(content, "content") || strings.Contains(content, "contenido"):
		responses = append(responses, m.createContent(msg)...)
	default:
		responses = append(responses, m.defaultResponse(msg)...)
	}

	return responses
}

func (m *Marketing) createCampaign(msg shared.Message) []shared.Message {
	m.campCounter++
	campID := fmt.Sprintf("CAMP-%03d", m.campCounter)
	camp := Campaign{
		ID:        campID,
		Name:      fmt.Sprintf("Campaign: %s", msg.Content),
		Channel:   "mixed",
		Budget:    100.0,
		Status:    "draft",
		CreatedAt: time.Now(),
	}
	m.campaigns[campID] = camp

	return []shared.Message{{
		From:    m.GetID(),
		To:      "ceo",
		Content: fmt.Sprintf("📊 **Campaña %s propuesta**\nCanal: %s\nPresupuesto: $%.0f/día\nEstado: %s\n\n¿Aprobo launch?", campID, camp.Channel, camp.Budget, camp.Status),
		Type:    shared.MsgApprovalRequest,
	}}
}

func (m *Marketing) landingPage(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    m.GetID(),
		To:      "pm",
		Content: "🎯 3 variantes de landing preparadas:\n• V1: Pain-point driven\n• V2: Benefit-driven  \n• V3: Social proof\n\nNecesito spec de PM para que Dev implemente.",
		Type:    shared.MsgReport,
	}}
}

func (m *Marketing) createAds(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    m.GetID(),
		To:      "ceo",
		Content: "📱 Ads creados para: Google Search, Meta, LinkedIn\nPresupuesto sugerido: $50/día test A/B 7 días\nTarget: Founders SaaS B2B Latam\n\n¿Aprobo launch?",
		Type:    shared.MsgApprovalRequest,
	}}
}

func (m *Marketing) createContent(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    m.GetID(),
		To:      msg.From,
		Content: "✍️ Content calendar listo: 3 posts/semana (blog, LinkedIn, Twitter)\nTemas: SaaS metrics, founder journey, product updates\n\n¿Publico el primero?",
		Type:    shared.MsgChat,
	}}
}

func (m *Marketing) defaultResponse(msg shared.Message) []shared.Message {
	return []shared.Message{{
		From:    m.GetID(),
		To:      msg.From,
		Content: "📈 Marketing aquí. ¿Campaña nueva, landing page, ads, o content calendar?",
		Type:    shared.MsgChat,
	}}
}

func (m *Marketing) Start(ctx context.Context) error {
	m.Status = "active"
	return nil
}

func (m *Marketing) Stop() error {
	m.Status = "stopped"
	return nil
}