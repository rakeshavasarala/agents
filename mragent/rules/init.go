package rules

import "github.com/rakeshavasarala/agents/mragent/pkg/core"

func init() {
	core.RegisterRule(SensitiveFileRule{})
		core.RegisterRule(HelmChartVersionAndChangelogRule{})
}