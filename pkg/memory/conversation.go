package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/mauroociappinaph/ayrton/internal/engram"
)

// ConversationMemory gestiona el historial de chat compartido
type ConversationMemory struct {
	engramClient *engram.Client
	sessionID    string
}

func NewConversationMemory(engramClient *engram.Client, sessionID string) *ConversationMemory {
	return &ConversationMemory{
		engramClient: engramClient,
		sessionID:    sessionID,
	}
}

// SaveMessage guarda un mensaje del chat
func (c *ConversationMemory) SaveMessage(from, to, content, msgType string) error {
	if c.engramClient == nil {
		return nil
	}

	obs := &engram.Observation{
		Title:    fmt.Sprintf("%s → %s", from, to),
		Type:     "conversation",
		Scope:    "project",
		TopicKey: "chat/" + c.sessionID,
		Content: fmt.Sprintf(`**From:** %s
**To:** %s
**Type:** %s
**Content:** %s
**Timestamp:** %s`, from, to, msgType, content, time.Now().Format(time.RFC3339)),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.engramClient.SaveOrUpdate(ctx, obs)
	return err
}

// GetHistory recupera historial del chat
func (c *ConversationMemory) GetHistory(limit int) ([]engram.SearchResult, error) {
	if c.engramClient == nil {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.engramClient.Search(ctx, "", engram.SearchOptions{
		Type:    "conversation",
		Project: "project",
		Scope:   "project",
		Limit:   limit,
	})
}

// GetContextForAgent recupera contexto relevante para un agente
func (c *ConversationMemory) GetContextForAgent(agentRole, task string, limit int) (string, error) {
	if c.engramClient == nil {
		return "", nil
	}

	query := fmt.Sprintf("%s %s", agentRole, task)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results, err := c.engramClient.Search(ctx, query, engram.SearchOptions{
		Project: "project",
		Scope:   "project",
		Limit:   limit,
	})
	if err != nil {
		return "", err
	}

	return c.formatResults(results), nil
}

func (c *ConversationMemory) formatResults(results []engram.SearchResult) string {
	if len(results) == 0 {
		return ""
	}

	var lines []string
	for _, r := range results {
		lines = append(lines, fmt.Sprintf("• %s: %s", r.Title, truncate(r.Content, 200)))
	}
	return "Contexto relevante:\n" + fmt.Sprintf("%s", lines)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}