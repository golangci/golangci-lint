package result

type Result struct {
	Issues           []Issue
	MaxIssuesPerFile int // Needed for gofmt and goimports where it is 1
}
