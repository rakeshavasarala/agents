package interfaces

// RuleEvaluator is implemented by any rule that wants to inspect an input and return a result.
// Useful for security checks, compliance policies, code quality, etc.
type RuleEvaluator interface {
	ID() string
	Evaluate(input RuleInput) RuleResult
}

// Agent is a high-level interface that runs a set of rules or logic on a given input.
type Agent interface {
	ID() string
	Evaluate(input RuleInput) []RuleResult
}

type ParametrizedRuleEvaluator interface {
	RuleEvaluator
	EvaluateWithParams(input RuleInput, params map[string]interface{}) RuleResult
}