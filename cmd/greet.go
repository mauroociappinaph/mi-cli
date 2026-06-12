package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var greetCmd = &cobra.Command{
	Use:   "greet [name]",
	Short: "Saluda a alguien",
	Long:  `Saluda a la persona indicada con un mensaje amistoso.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := "mundo"
		if len(args) > 0 {
			name = args[0]
		}

		greeting := fmt.Sprintf("¡Hola, %s! 👋", name)

		if viper.GetString("output") == "json" {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "{\"greeting\": \"%s\"}\n", greeting)
		} else {
			cmd.Println(greeting)
		}
	},
}

func init() {
	rootCmd.AddCommand(greetCmd)

	greetCmd.Flags().StringP("lang", "l", "es", "idioma del saludo (es|en)")
	_ = viper.BindPFlag("lang", greetCmd.Flags().Lookup("lang"))
}