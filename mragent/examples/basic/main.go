package main

import (
	"fmt"

	"github.com/rakeshavasarala/agents/mragent/pkg/core"
	_ "github.com/rakeshavasarala/agents/mragent/rules"
	"github.com/rakeshavasarala/agents/shared/interfaces"
)

func main() {
	// Typed MR metadata (no coverage)
	mr := interfaces.MergeRequestMetadata{
		Author:       "rakesh",
		Project:      "platform-service",
		TargetBranch: "main",
		SourceBranch: "feature/sensitive-update",
		IsDraft:      false,
	}

	input := interfaces.RuleInput{
		Source:   "gitlab",
		Subject:  "Fix CI pipeline",
		Content:  "Adjusts timeouts and updates chart.",
		Files:    []string{".gitlab-ci.yml", "charts/api/Chart.yaml"},
		Diff: `
		- version: 1.2.3
		+ version: 1.2.4
		`,
		Tags:     []string{"ci", "helm"},
		//Metadata: map[string]string{"author": "rakesh"}, // optional
		Raw:      mr,
	}

	rules := core.GetRegisteredRules()
	for _, rule := range rules {
		result := rule.Evaluate(input)
		icon := "✅"
		if !result.Passed {
			icon = "❌"
		}
		fmt.Printf("%s Rule: %-40s | Passed: %-5v | Message: %s\n",
			icon, result.RuleID, result.Passed, result.Message)
			}
}