package interfaces

import (
	"strconv"
)

type MergeRequestMetadata struct {
	Author        string
	Project       string
	TargetBranch  string
	SourceBranch  string
	CoverageDelta float64
	IsDraft       bool
}

// ParseMRMetadata safely extracts MR metadata from RuleInput.Metadata
func ParseMRMetadata(meta map[string]string) MergeRequestMetadata {
	coverage := 0.0
	if raw, ok := meta["coverage_delta"]; ok {
		if f, err := strconv.ParseFloat(raw, 64); err == nil {
			coverage = f
		}
	}
	isDraft := false
	if meta["is_draft"] == "true" {
		isDraft = true
	}
	return MergeRequestMetadata{
		Author:        meta["author"],
		Project:       meta["project"],
		TargetBranch:  meta["target_branch"],
		SourceBranch:  meta["source_branch"],
		CoverageDelta: coverage,
		IsDraft:       isDraft,
	}
}