package roles

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mauroociappinaph/ayrton/internal/engram"
	"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type Learning struct {
	*shared.BaseAgent
	broadcaster  shared.BroadcasterInterface
	engramClient *engram.Client
	patterns     map[string]*Pattern
}

type Pattern = engram.Pattern

func NewLearning(engramClient *engram.Client, broadcaster shared.BroadcasterInterface) *Learning {
	base := shared.NewBaseAgent("learning", "Learning", "learning")
	l := &Learning{
		BaseAgent:    base,
		broadcaster:  broadcaster,
		engramClient: engramClient,
		patterns:     make(map[string]*Pattern),
	}
	l.loadPatterns()
	return l
}

func (l *Learning) loadPatterns() {
	if l.engramClient == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results, err := l.engramClient.SearchWithOptions(ctx, "", engram.SearchOptions{
		Type:    "pattern",
		Project: "project",
		Scope:   "project",
		Limit:   100,
	})
	if err != nil {
		return
	}
	for _, r := range results {
		if p := l.parsePattern(r.Content); p != nil {
			l.patterns[p.Trigger] = p
		}
	}
}

func (l *Learning) Receive(ctx context.Context, msg shared.Message) []shared.Message {
	content := strings.ToLower(msg.Content)
	var responses []shared.Message

	switch {
	case strings.Contains(content, "learn") || strings.Contains(content, "aprende"):
		responses = append(responses, l.extractLesson(msg)...)
	case strings.Contains(content, "pattern") || strings.Contains(content, "patrón"):
		responses = append(responses, l.showPatterns(msg)...)
	case strings.Contains(content, "recall") || strings.Contains(content, "recuerda"):
		responses = append(responses, l.recall(msg)...)
	default:
		responses = append(responses, l.maybeLearnFromChat(msg)...)
	}

	return responses
}

func (l *Learning) extractLesson(msg shared.Message) []shared.Message {
	payload := msg.Payload
	if payload == nil {
		payload = make(map[string]interface{})
	}

	trigger := getString(payload, "trigger")
	action := getString(payload, "action")
	outcome := getString(payload, "outcome")
	agents := getStringSlice(payload, "agents")
	ctx := getMap(payload, "context")

	if trigger == "" || action == "" || outcome == "" {
		return []shared.Message{{
			From:    "learning",
			To:      msg.From,
			Content: "Uso: trigger=\"...\" action=\"...\" outcome=\"...\" agents=[...] context={...}",
			Type:    shared.MsgChat,
		}}
	}

	id, err := l.engramClient.SaveLesson(trigger, action, outcome, ctx, agents, 0.5)
	if err != nil {
		return []shared.Message{{
			From:    "learning",
			To:      msg.From,
			Content: fmt.Sprintf("❌ Error guardando lesson: %v", err),
			Type:    shared.MsgError,
		}}
	}

	p := l.updatePattern(trigger, action, outcome, agents, ctx)

	if p.Confidence > 0.8 {
		return l.proposeStrategyChange(p)
	}

	return []shared.Message{{
		From:    "learning",
		To:      msg.From,
		Content: fmt.Sprintf("🧠 Lesson guardada (ID: %d): %s → %s [confidence: %.2f]", id, trigger, outcome, p.Confidence),
		Type:    shared.MsgReport,
	}}
}

func (l *Learning) updatePattern(trigger, action, outcome string, agents []string, context map[string]interface{}) *Pattern {
	p, exists := l.patterns[trigger]
	if !exists {
		p = &Pattern{
			Trigger:    trigger,
			Action:     action,
			Outcome:    outcome,
			Frequency:  0,
			Confidence: 0,
			Agents:     agents,
		}
		l.patterns[trigger] = p
	}
	p.Frequency++
	p.LastSeen = time.Now()
	p.Confidence = l.calculateConfidence(p)

	_ = l.engramClient.SavePattern(trigger, p)
	return p
}

func (l *Learning) calculateConfidence(p *Pattern) float64 {
	base := float64(p.Frequency) / 10.0
	if base > 1.0 {
		base = 1.0
	}
	return base * 0.9
}

func (l *Learning) proposeStrategyChange(p *Pattern) []shared.Message {
	return []shared.Message{{
		From:    "learning",
		To:      "ceo",
		Content: fmt.Sprintf(`💡 **Patrón detectado (confidence: %.2f)**
**Trigger:** %s
**Acción recurrente:** %s
**Outcome:** %s
**Frecuencia:** %d veces

**Propuesta:** Cambiar estrategia para optimizar este patrón.`, p.Confidence, p.Trigger, p.Action, p.Outcome, p.Frequency),
		Type: shared.MsgReport,
	}}
}

func (l *Learning) showPatterns(msg shared.Message) []shared.Message {
	if len(l.patterns) == 0 {
		return []shared.Message{{
			From:    "learning",
			To:      msg.From,
			Content: "📝 No hay patrones aprendidos aún.",
			Type:    shared.MsgReport,
		}}
	}

	var lines []string
	for _, p := range l.patterns {
		lines = append(lines, fmt.Sprintf("  • %s → %s (freq: %d, conf: %.2f)", p.Trigger, p.Outcome, p.Frequency, p.Confidence))
	}

	return []shared.Message{{
		From:    "learning",
		To:      msg.From,
		Content: fmt.Sprintf("📚 **Patrones aprendidos (%d):**\n%s", len(l.patterns), strings.Join(lines, "\n")),
		Type:    shared.MsgReport,
	}}
}

func (l *Learning) recall(msg shared.Message) []shared.Message {
	query := strings.TrimSpace(strings.TrimPrefix(msg.Content, "recall"))
	query = strings.TrimSpace(strings.TrimPrefix(query, "recuerda"))

	if l.engramClient == nil {
		return []shared.Message{{
			From:    "learning",
			To:      msg.From,
			Content: "Engram no disponible",
			Type:    shared.MsgError,
		}}
	}

	results, err := l.engramClient.GetLessons(query, 10)
	if err != nil {
		return []shared.Message{{
			From:    "learning",
			To:      msg.From,
			Content: fmt.Sprintf("❌ Error: %v", err),
			Type:    shared.MsgError,
		}}
	}

	if len(results) == 0 {
		return []shared.Message{{
			From:    "learning",
			To:      msg.From,
			Content: "No se encontraron lessons para: " + query,
			Type:    shared.MsgReport,
		}}
	}

	var lines []string
	for _, r := range results {
		lines = append(lines, fmt.Sprintf("  • %s (rank: %.2f)", r.Title, r.Rank))
	}

	return []shared.Message{{
		From:    "learning",
		To:      msg.From,
		Content: fmt.Sprintf("🔍 **Recall para '%s' (%d resultados):**\n%s", query, len(results), strings.Join(lines, "\n")),
		Type:    shared.MsgReport,
	}}
}

func (l *Learning) maybeLearnFromChat(msg shared.Message) []shared.Message {
	if msg.Type == shared.MsgReport && strings.Contains(strings.ToLower(msg.Content), "✅") {
		// Could extract lesson here
	}
	return nil
}

func (l *Learning) parsePattern(content string) *Pattern {
	lines := strings.Split(content, "\n")
	p := &Pattern{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "**Trigger:** ") {
			p.Trigger = strings.TrimPrefix(line, "**Trigger:** ")
		} else if strings.HasPrefix(line, "**Action:** ") {
			p.Action = strings.TrimPrefix(line, "**Action:** ")
		} else if strings.HasPrefix(line, "**Outcome:** ") {
			p.Outcome = strings.TrimPrefix(line, "**Outcome:** ")
		} else if strings.HasPrefix(line, "**Frequency:** ") {
			fmt.Sscanf(strings.TrimPrefix(line, "**Frequency:** "), "%d", &p.Frequency)
		} else if strings.HasPrefix(line, "**Confidence:** ") {
			fmt.Sscanf(strings.TrimPrefix(line, "**Confidence:** "), "%f", &p.Confidence)
		} else if strings.HasPrefix(line, "**Last Seen:** ") {
			if t, err := time.Parse(time.RFC3339, strings.TrimPrefix(line, "**Last Seen:** ")); err == nil {
				p.LastSeen = t
			}
		} else if strings.HasPrefix(line, "**Agents:** ") {
			// Parse agents slice
		}
	}

	if p.Trigger == "" {
		return nil
	}
	return p
}

func (l *Learning) Start(ctx context.Context) error {
	l.Status = "active"
	return nil
}

func (l *Learning) Stop() error {
	l.Status = "stopped"
	return nil
}

// Helpers for extracting from Payload
func getString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getStringSlice(m map[string]interface{}, key string) []string {
	if m == nil {
		return nil
	}
	if v, ok := m[key].([]string); ok {
		return v
	}
	return nil
}

func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	if v, ok := m[key].(map[string]interface{}); ok {
		return v
	}
	return nil
}