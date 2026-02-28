package main

import (
	"github.com/shell-sage/cmd"
	_ "github.com/shell-sage/internal/ollama" // registers the ollama provider via init()
)

func main() {
	cmd.Execute()
}
