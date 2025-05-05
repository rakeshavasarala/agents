module github.com/rakeshavasarala/agents/mragent

replace github.com/rakeshavasarala/agents/shared => ../shared

go 1.24.2

require (
	github.com/rakeshavasarala/agents/shared v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.1

)

require (
	github.com/Masterminds/semver v1.5.0
	github.com/Masterminds/semver/v3 v3.3.1 // indirect
)
