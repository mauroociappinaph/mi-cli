package cmd

import (
	"fmt"
	"os"

	"github.com/mauroociappinaph/ayrton/internal/engram"
	"github.com/mauroociappinaph/ayrton/pkg/chat"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Inicia el chat multi-agente de Ayrton",
	Long: `Inicia una sesión de chat colaborativo con agentes autónomos.

Comandos disponibles en el chat:
  /help              - Muestra ayuda
  /list              - Lista agentes activos
  /spawn <rol> "nombre" - Crea agente (pm, dev, marketing, ops, prospeccion)
  /kill <nombre>     - Elimina agente
  /history           - Historial del chat
  /status            - Estado del sistema
  /quit              - Salir

Chat:
  @rol mensaje       - Mensaje dirigido a un agente
  mensaje            - Broadcast a todos

Roles disponibles:
  pm           - Product Manager (specs, backlog, priorización)
  dev          - Developer (código, tests, deploy)
  marketing    - Marketing (campañas, contenido, ads)
  ops          - Operations (presupuesto, recursos, approvals)
  prospeccion  - Prospección (research, leads, outreach)

Meta-agentes (auto-spawn):
  auditor      - Auditoría código/calidad/seguridad
  learning     - Aprendizaje patrones + Engram
  engram       - Memoria persistente`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize Engram client
		engramClient, err := engram.NewClient()
		if err != nil {
			return fmt.Errorf("inicializar Engram: %w", err)
		}
		defer engramClient.Close()

		// Create and start REPL
		repl := chat.NewREPL(engramClient)
		return repl.Start()
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}