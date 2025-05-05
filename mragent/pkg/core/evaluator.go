package core

import (
	"github.com/rakeshavasarala/agents/mragent/pkg/config"
	"github.com/rakeshavasarala/agents/shared/interfaces"
)

func getRuleConfig(id string, ruleset *config.Ruleset) *config.RuleConfig {
	for _, rc := range ruleset.Rules {
		if rc.ID == id {
			return &rc
		}
	}
	return nil
}

func EvaluateAll(input interfaces.RuleInput, ruleset *config.Ruleset) []interfaces.RuleResult {
	results := []interfaces.RuleResult{}
	for _, rule := range GetRegisteredRules() {
		cfg := getRuleConfig(rule.ID(), ruleset)
		if cfg != nil && cfg.Enabled {
			if ruleWithParams, ok := rule.(interfaces.ParametrizedRuleEvaluator); ok {
				results = append(results, ruleWithParams.EvaluateWithParams(input, cfg.Params))
			} else {
				results = append(results, rule.Evaluate(input))
			}
		}
	}
	return results
}