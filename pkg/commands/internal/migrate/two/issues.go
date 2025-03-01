package two

type Issues struct {
	MaxIssuesPerLinter *int  `yaml:"max-issues-per-linter,omitempty" toml:"max-issues-per-linter,omitempty"`
	MaxSameIssues      *int  `yaml:"max-same-issues,omitempty" toml:"max-same-issues,omitempty"`
	UniqByLine         *bool `yaml:"uniq-by-line,omitempty" toml:"uniq-by-line,omitempty"`

	DiffFromRevision  *string `yaml:"new-from-rev,omitempty" toml:"new-from-rev,omitempty"`
	DiffFromMergeBase *string `yaml:"new-from-merge-base,omitempty" toml:"new-from-merge-base,omitempty"`
	DiffPatchFilePath *string `yaml:"new-from-patch,omitempty" toml:"new-from-patch,omitempty"`
	WholeFiles        *bool   `yaml:"whole-files,omitempty" toml:"whole-files,omitempty"`
	Diff              *bool   `yaml:"new,omitempty" toml:"new,omitempty"`

	NeedFix *bool `yaml:"fix,omitempty" toml:"fix,omitempty"`
}
