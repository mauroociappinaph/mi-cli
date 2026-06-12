package chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mauroociappinaph/ayrton/internal/engram"
	"github.com/mauroociappinaph/ayrton/pkg/agent"
	"github.com/mauroociappinaph/ayrton/pkg/agent/roles"
	"github.com/mauroociappinaph/ayrton/pkg/runtime"
	"github.com/mauroociappinaph/ayrton/pkg/shared"
)

type REPL struct {
	broadcaster   shared.BroadcasterInterface
	agentRegistry *agent.Registry
	runtime       *runtime.Pool
	engramClient  *engram.Client
	sessionID     string
	scanner       *bufio.Scanner
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewREPL(engramClient *engram.Client) *REPL {
	ctx, cancel := context.WithCancel(context.Background())
	sessionID := fmt.Sprintf("ayrton-chat-%d", time.Now().Unix())

	b := NewBroadcaster(engramClient, sessionID)
	reg := agent.NewRegistry(b)
	pool := runtime.NewPool(b)

	return &REPL{
		broadcaster:   b,
		agentRegistry: reg,
		runtime:       pool,
		engramClient:  engramClient,
		sessionID:     sessionID,
		scanner:       bufio.NewScanner(os.Stdin),
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (r *REPL) Start() error {
	// Auto-spawn meta-agents
	r.spawnMetaAgents()

	fmt.Println("🤖 Ayrton Chat iniciado")
	fmt.Println("   Escribe '@ayrton spawn <rol> \"nombre\"' para crear agentes")
	fmt.Println("   Usa '@rol mensaje' para hablar con un agente")
	fmt.Println("   Comandos: /help, /list, /kill, /history, /quit")
	fmt.Println()

	// Start runtime
	r.runtime.Start()

	// Main loop
	for {
		select {
		case <-r.ctx.Done():
			return nil
		default:
			fmt.Print("> ")
			if !r.scanner.Scan() {
				return nil // EOF
			}

			input := r.scanner.Text()
			if err := r.processInput(input); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}
		}
	}
}

func (r *REPL) processInput(input string) error {
	parsed := Parse(input, "ceo")
	if parsed == nil {
		return nil
	}

	switch parsed.Type {
	case MsgCommand:
		return r.handleCommand(parsed)
	case MsgMention:
		return r.handleMention(parsed)
	case MsgBroadcast:
		return r.handleBroadcast(parsed)
	}
	return nil
}

func (r *REPL) handleCommand(parsed *ParsedMessage) error {
	cmd, args := parsed.To, parsed.Content

	switch cmd {
	case "help":
		r.printHelp()
	case "list":
		r.listAgents()
	case "spawn":
		return r.handleSpawn(args)
	case "kill":
		return r.handleKill(args)
	case "history":
		r.showHistory()
	case "quit", "exit":
		r.cancel()
	case "status":
		r.showStatus()
	default:
		fmt.Printf("❓ Comando desconocido: /%s (usa /help)\n", cmd)
	}
	return nil
}

func (r *REPL) handleSpawn(args string) error {
	parts := strings.Fields(args)
	if len(parts) < 2 {
		fmt.Println("Uso: /spawn <rol> \"nombre\"")
		fmt.Println("Roles disponibles: pm, dev, marketing, ops, prospeccion")
		return nil
	}

	role := parts[0]
	name := strings.Join(parts[1:], " ")
	name = strings.Trim(name, "\"")

	_, err := r.agentRegistry.Spawn(role, name)
	if err != nil {
		return err
	}
	fmt.Printf("✅ %s (%s) se unió al chat\n", name, role)
	return nil
}

func (r *REPL) handleKill(args string) error {
	name := strings.TrimSpace(args)
	if name == "" {
		fmt.Println("Uso: /kill <nombre_agente>")
		return nil
	}

	if err := r.agentRegistry.Kill(name); err != nil {
		return err
	}
	fmt.Printf("🗑️ Agente %s eliminado\n", name)
	return nil
}

func (r *REPL) handleMention(parsed *ParsedMessage) error {
	msg := shared.Message{
		ID:        fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		From:      "ceo",
		To:        parsed.To,
		Content:   parsed.Content,
		Type:      shared.MsgMention,
		ThreadID:  fmt.Sprintf("thread-%d", time.Now().Unix()),
		Timestamp: time.Now(),
	}

	r.broadcaster.Send(msg)
	return nil
}

func (r *REPL) handleBroadcast(parsed *ParsedMessage) error {
	msg := shared.Message{
		ID:        fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		From:      "ceo",
		To:        "",
		Content:   parsed.Content,
		Type:      shared.MsgBroadcast,
		ThreadID:  fmt.Sprintf("thread-%d", time.Now().Unix()),
		Timestamp: time.Now(),
	}

	r.broadcaster.Send(msg)
	return nil
}

func (r *REPL) printHelp() {
	fmt.Println(`📖 Comandos disponibles:
  /help                    - Muestra esta ayuda
  /list                    - Lista agentes activos
  /spawn <rol> "nombre"    - Crea un agente (roles: pm, dev, marketing, ops, prospeccion)
  /kill <nombre>           - Elimina un agente
  /history                 - Muestra historial del chat
  /status                  - Estado del sistema
  /quit                    - Salir del chat

💬 Chat:
  @rol mensaje             - Envía mensaje a un agente específico
  mensaje                  - Broadcast a todos los agentes

👥 Roles disponibles:
  pm           - Product Manager (specs, backlog, priorización)
  dev          - Developer (código, tests, deploy)
  marketing    - Marketing (campañas, contenido, ads)
  ops          - Operations (presupuesto, recursos, approvals)
  prospeccion  - Prospección (research, leads, outreach)

🤖 Meta-agentes (auto-spawn):
  auditor      - Auditoría código/calidad/seguridad
  learning     - Aprendizaje patrones + Engram
  engram       - Memoria persistente`)
}

func (r *REPL) listAgents() {
	agents := r.agentRegistry.List()
	if len(agents) == 0 {
		fmt.Println("   (ningún agente activo)")
		return
	}
	for _, a := range agents {
		fmt.Printf("   • %s (%s) - %s\n", a.Name, a.Role, a.Status)
	}
}

func (r *REPL) showHistory() {
	fmt.Println("📜 Historial del chat (últimos 20):")
	// TODO: Implementar lectura desde Engram
	fmt.Println("   (pendiente: leer de Engram)")
}

func (r *REPL) showStatus() {
	agents := r.agentRegistry.List()
	fmt.Printf("📊 Estado del sistema:\n")
	fmt.Printf("   Sesión: %s\n", r.sessionID)
	fmt.Printf("   Agentes activos: %d\n", len(agents))
	fmt.Printf("   Runtime: %s\n", r.runtime.Status())
}

func (r *REPL) spawnMetaAgents() {
	// Spawn Auditor
	auditor := roles.NewAuditor(r.engramClient, r.broadcaster)
	r.agentRegistry.Register(auditor)
	r.runtime.Spawn(auditor)

	// Spawn Learning (usa tu Learning agent existente)
	learning := roles.NewLearning(r.engramClient, r.broadcaster)
	r.agentRegistry.Register(learning)
	r.runtime.Spawn(learning)

	fmt.Println("🤖 Meta-agentes iniciados: Auditor, Learning, Engram")
}