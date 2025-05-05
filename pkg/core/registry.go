package core

import (
	"github.com/rakeshavasarala/agents/shared/interfaces"
)

var ruleRegistry = make(map[string]interfaces.RuleEvaluator)

// RegisterRule adds a rule to the internal registry.
func RegisterRule(rule interfaces.RuleEvaluator) {
	ruleRegistry[rule.ID()] = rule
}

// GetRegisteredRules returns all active rules.
func GetRegisteredRules() []interfaces.RuleEvaluator {
	rules := []interfaces.RuleEvaluator{}
	for _, r := range ruleRegistry {
		rules = append(rules, r)
	}
	return rules
}