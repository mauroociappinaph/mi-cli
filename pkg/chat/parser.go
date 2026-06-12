package chat

import (
	"regexp"
	"strings"
)

type ParsedMessage struct {
	From      string // "ceo" (user) or agent role
	To        string // "" = broadcast, "@rol" = directed, "/comando" = command
	Content   string
	Type      MessageType
	Raw       string
}

type MessageType int

const (
	MsgChat MessageType = iota
	MsgCommand
	MsgMention
	MsgBroadcast
)

var (
	mentionRegex = regexp.MustCompile(`^@(\w+)\s+(.+)$`)
	commandRegex = regexp.MustCompile(`^/(\w+)(?:\s+(.+))?$`)
)

func Parse(input string, from string) *ParsedMessage {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	// Check for @mention
	if m := mentionRegex.FindStringSubmatch(input); m != nil {
		return &ParsedMessage{
			From:    from,
			To:      m[1],
			Content: m[2],
			Type:    MsgMention,
			Raw:     input,
		}
	}

	// Check for /command
	if m := commandRegex.FindStringSubmatch(input); m != nil {
		return &ParsedMessage{
			From:    from,
			To:      "/" + m[1],
			Content: m[2],
			Type:    MsgCommand,
			Raw:     input,
		}
	}

	// Regular chat message (broadcast to all)
	return &ParsedMessage{
		From:    from,
		To:      "",
		Content: input,
		Type:    MsgBroadcast,
		Raw:     input,
	}
}

func IsSystemCommand(input string) bool {
	return commandRegex.MatchString(strings.TrimSpace(input))
}

func GetCommand(input string) (string, string) {
	m := commandRegex.FindStringSubmatch(strings.TrimSpace(input))
	if m == nil {
		return "", ""
	}
	return m[1], m[2]
}