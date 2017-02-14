package uiutil

import (
	"fmt"
	"github.com/robbiev/dilemma"
	"os"
)

func CreateDilemma(title, help string, options []string, ShownItems int) string {
	thisDilemma := dilemma.Config{
		Title:      title,
		Help:       help,
		Options:    options,
		ShownItems: ShownItems,
	}
	selected, exitKey, err := dilemma.Prompt(thisDilemma)
	if err != nil || exitKey == dilemma.CtrlC {
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	return selected
}
