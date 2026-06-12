package main

import (
	"os"

	"github.com/tuusuario/mi-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}