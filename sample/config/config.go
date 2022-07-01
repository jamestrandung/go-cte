package config

import (
	"fmt"

	"github.com/jamestrandung/go-die/die"
)

var Engine = &CostEngine{
	Engine: die.NewEngine(),
}

var printDebugLog = true

func Print(values ...any) {
	if printDebugLog {
		fmt.Println(values...)
	}
}

func Printf(format string, values ...any) {
	if printDebugLog {
		fmt.Printf(format, values...)
	}
}

type CostEngine struct {
	die.Engine
	// Add common utilities like logger, statsD, UCM client, etc.
	// for all component codes to share.
}
