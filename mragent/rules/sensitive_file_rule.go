package rules

import (
	"strings"

	"github.com/rakeshavasarala/agents/shared/interfaces"
	"github.com/rakeshavasarala/agents/shared/utils"
)

type SensitiveFileRule struct{}

func (r SensitiveFileRule) ID() string {
	return "sensitive-file-check"
}

// Evaluate falls back to EvaluateWithParams with default patterns
func (r SensitiveFileRule) Evaluate(input interfaces.RuleInput) interfaces.RuleResult {
	return r.EvaluateWithParams(input, map[string]interface{}{
		"patterns": []string{".gitlab-ci.yml", "secrets.yaml"},
	})
}

func (r SensitiveFileRule) EvaluateWithParams(input interfaces.RuleInput, params map[string]interface{}) interfaces.RuleResult {
	patterns := utils.GetStringSlice(params, "patterns")
	if len(patterns) == 0 {
		return interfaces.RuleResult{
			RuleID:     r.ID(),
			Applicable: false,
			Passed:     true,
			Message:    "No sensitive file patterns provided",
			Severity:   interfaces.SeverityInfo,
		}
	}

	for _, file := range input.Files {
		for _, pattern := range patterns {
			if strings.Contains(file, pattern) {
				// Optionally use typed data from Raw
				if mr, ok := input.Raw.(interfaces.MergeRequestMetadata); ok {
					return interfaces.RuleResult{
						RuleID:     r.ID(),
						Applicable: true,
						Passed:     false,
						Severity:   interfaces.SeverityHigh,
						Message:    "Sensitive file changed: " + file + " (author: " + mr.Author + ")",
					}
				}

				return interfaces.RuleResult{
					RuleID:     r.ID(),
					Applicable: true,
					Passed:     false,
					Severity:   interfaces.SeverityHigh,
					Message:    "Sensitive file changed: " + file,
				}
			}
		}
	}

	return interfaces.RuleResult{
		RuleID:     r.ID(),
		Applicable: true,
		Passed:     true,
		Message:    "No sensitive files changed",
		Severity:   interfaces.SeverityInfo,
	}
}