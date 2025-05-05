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
	ruleset, err := config.LoadRulesConfig("examples/helm_eval/rules.yaml")
	if err != nil {
		log.Fatalf("Failed to load rules config: %v", err)
	}

	input := interfaces.RuleInput{
		Title:       "Update Helm chart for API service",
		Author:      "team-b-dev",
		Files: []string{
			"charts/api-service/Chart.yaml",
			"charts/api-service/templates/deployment.yaml",
			"custom-changelog.md",
		},
		Diff:     "version: 1.2.4\n# updated deployment\n",
		Metadata: map[string]string{"team": "platform"},
	}

	results := core.EvaluateAll(input, ruleset)

	for _, r := range results {
		fmt.Printf("Rule: %-40s | Passed: %-5v | Message: %s\n",
			r.RuleID, r.Passed, r.Message)
	}
}