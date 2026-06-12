package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var autonomousCmd = &cobra.Command{
	Use:   "autonomous [issue-number]",
	Short: "Ejecuta loop autónomo SDD para issue con label 'autonomous'",
	Long:  `Procesa un issue de GitHub mediante pipeline SDD completo: propose → spec → design → tasks → apply → verify → archive`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		issueNumber := args[0]
		fmt.Fprintf(cmd.OutOrStdout(), "🤖 Loop autónomo iniciado para issue #%s\n", issueNumber)
		fmt.Println("📋 Pipeline SDD: propose → spec → design → tasks → apply → verify → archive")
		
		// TODO: Implementar cada fase SDD
		fmt.Println("⚠️  Implementación pendiente - cada fase invocará skills SDD correspondientes")
	},
}

func init() {
	rootCmd.AddCommand(autonomousCmd)
	
	autonomousCmd.Flags().Bool("dry-run", false, "Simula sin crear branch/PR")
	autonomousCmd.Flags().String("label", "autonomous", "Label que dispara el loop")
}
