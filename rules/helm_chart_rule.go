package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/rakeshavasarala/agents/shared/interfaces"
	"github.com/rakeshavasarala/agents/shared/utils"
)

type HelmChartVersionAndChangelogRule struct{}

func (r HelmChartVersionAndChangelogRule) ID() string {
	return "helm-chart-version-and-changelog"
}

func (r HelmChartVersionAndChangelogRule) Evaluate(input interfaces.RuleInput) interfaces.RuleResult {
	return r.EvaluateWithParams(input, map[string]interface{}{
		"changelog":    "CHANGELOG.md",
		"chart_paths":  []string{"charts/"},
	})
}

func (r HelmChartVersionAndChangelogRule) EvaluateWithParams(input interfaces.RuleInput, params map[string]interface{}) interfaces.RuleResult {
	changelog := utils.GetString(params, "changelog")
	if changelog == "" {
		changelog = "CHANGELOG.md"
	}
	chartPaths := utils.GetStringSlice(params, "chart_paths")
	warnAboutChartPath := len(chartPaths) == 0

	// Skip if draft MR
	if mr, ok := input.Raw.(interfaces.MergeRequestMetadata); ok && mr.IsDraft {
		return interfaces.RuleResult{
			RuleID:     r.ID(),
			Applicable: true,
			Passed:     true,
			Message:    "Skipped Helm chart checks for draft MR",
			Severity:   interfaces.SeverityInfo,
		}
	}

	chartModified, _ := chartWasModified(input.Files, chartPaths)
	if !chartModified {
		msg := "No Helm chart changes detected"
		if warnAboutChartPath {
			msg += " [warning: 'chart_paths' not configured, all files considered Helm-related]"
		}
		return interfaces.RuleResult{
			RuleID:     r.ID(),
			Applicable: false,
			Passed:     true,
			Message:    msg,
			Severity:   interfaces.SeverityInfo,
		}
	}

	// Check if changelog is touched
	changelogUpdated := changelogWasUpdated(input.Files, changelog)

	// Parse and compare versions
	oldVer, newVer, err := extractChartVersions(input.Diff)
	if err != nil {
		return interfaces.RuleResult{
			RuleID:     r.ID(),
			Applicable: true,
			Passed:     false,
			Message:    "Failed to parse chart versions: " + err.Error(),
			Severity:   interfaces.SeverityHigh,
		}
	}

	versionOk, versionMsg := checkVersionBump(oldVer, newVer)

	// Optional author label from Raw
	var author string
	if mr, ok := input.Raw.(interfaces.MergeRequestMetadata); ok {
		author = mr.Author
	}

	// Compose success or failure
	if versionOk && changelogUpdated {
		msg := "Helm chart version and changelog update detected"
		if warnAboutChartPath {
			msg += " [warning: 'chart_paths' not configured, all files considered Helm-related]"
		}
		if author != "" {
			msg += fmt.Sprintf(" (by %s)", author)
		}
		return interfaces.RuleResult{
			RuleID:     r.ID(),
			Applicable: true,
			Passed:     true,
			Message:    msg,
			Severity:   interfaces.SeverityInfo,
		}
	}

	failMsg := ""
	switch {
	case !versionOk && !changelogUpdated:
		failMsg = "Missing version bump and changelog update for Helm chart"
	case !versionOk:
		failMsg = versionMsg
	case !changelogUpdated:
		failMsg = "Helm chart changed but " + changelog + " was not updated"
	}
	if author != "" {
		failMsg += fmt.Sprintf(" (by %s)", author)
	}

	return interfaces.RuleResult{
		RuleID:     r.ID(),
		Applicable: true,
		Passed:     false,
		Message:    failMsg,
		Severity:   interfaces.SeverityHigh,
	}
}

func chartWasModified(files []string, paths []string) (bool, string) {
	for _, f := range files {
		if isHelmChartFile(f, paths) {
			return true, f
		}
	}
	return false, ""
}

func isHelmChartFile(path string, prefixes []string) bool {
	lower := strings.ToLower(path)

	if strings.HasSuffix(lower, "chart.yaml") || strings.HasSuffix(lower, "values.yaml") {
		if matchesPrefix(lower, prefixes) {
			return true
		}
	}

	if strings.Contains(lower, "/templates/") && strings.HasSuffix(lower, ".yaml") {
		if matchesPrefix(lower, prefixes) {
			return true
		}
	}

	return false
}

func matchesPrefix(path string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return true // default: allow any path
	}
	for _, p := range prefixes {
		if strings.HasPrefix(path, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

func changelogWasUpdated(files []string, changelog string) bool {
	for _, f := range files {
		if strings.EqualFold(f, changelog) || strings.Contains(strings.ToLower(f), strings.ToLower(changelog)) {
			return true
		}
	}
	return false
}

func extractChartVersions(diff string) (*semver.Version, *semver.Version, error) {
	versionRegex := regexp.MustCompile(`(?m)^[-+]\s*.*version:\s*([0-9]+\.[0-9]+\.[0-9]+)`)
	matches := versionRegex.FindAllStringSubmatch(diff, -1)

	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("no version lines found in diff")
	}

	var oldVerStr, newVerStr string
	if len(matches) >= 2 {
		oldVerStr = matches[0][1]
		newVerStr = matches[1][1]
	} else {
		newVerStr = matches[0][1]
	}

	newVer, errNew := semver.NewVersion(newVerStr)
	if errNew != nil {
		return nil, nil, errNew
	}

	var oldVer *semver.Version
	var errOld error
	if oldVerStr != "" {
		oldVer, errOld = semver.NewVersion(oldVerStr)
		if errOld != nil {
			oldVer = nil // still allow version bump w/o old version
		}
	}

	return oldVer, newVer, nil
}

func checkVersionBump(old, new *semver.Version) (bool, string) {
	if new == nil {
		return false, "Invalid or missing new chart version"
	}
	if old != nil && !new.GreaterThan(old) {
		return false, fmt.Sprintf("Version downgrade or no bump: %s → %s", old, new)
	}
	return true, "Valid version bump"
}