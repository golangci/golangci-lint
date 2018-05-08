package result

type Issue struct {
	FromLinter string
	Text       string
	File       string
	LineNumber int
	HunkPos    int
}

func (i Issue) FilePath() string {
	return i.File
}

func (i Issue) Line() int {
	return i.LineNumber
}

func NewIssue(fromLinter, text, file string, lineNumber, hunkPos int) Issue {
	return Issue{
		FromLinter: fromLinter,
		Text:       text,
		File:       file,
		LineNumber: lineNumber,
		HunkPos:    hunkPos,
	}
}
