package main

import (
	"fmt"
	"log"

	"github.com/rakeshavasarala/agents/mragent/pkg/config"
	"github.com/rakeshavasarala/agents/mragent/pkg/core"
	_ "github.com/rakeshavasarala/agents/mragent/rules"
	"github.com/rakeshavasarala/agents/shared/interfaces"
)

func main() {
	// Load rule config with parameters
	ruleset, err := config.LoadRulesConfig("examples/with_config/rules.yaml")
	if err != nil {
		log.Fatalf("Failed to load rules config: %v", err)
	}

	// Simulated MR input
	input := interfaces.RuleInput{
		Title:       "Update pipeline",
		Files:       []string{".gitlab-ci.yml", "README.md"},
		Description: "Changing CI timeout settings",
	}

	// Evaluate based on rules.yaml
	results := core.EvaluateAll(input, ruleset)

	for _, r := range results {
		fmt.Printf("Rule: %-25s | Passed: %-5v | Severity: %-8s | Message: %s\n",
			r.RuleID, r.Passed, r.Severity, r.Message)
	}
}