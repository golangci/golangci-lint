// Package changelog provides the release changelog to goreleaser.
package changelog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/git"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// ErrInvalidSortDirection happens when the sort order is invalid
var ErrInvalidSortDirection = errors.New("invalid sort direction")

// Pipe for checksums
type Pipe struct{}

func (Pipe) String() string {
	return "generating changelog"
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	// TODO: should probably have a different field for the filename and its
	// contents.
	if ctx.ReleaseNotes != "" {
		notes, err := loadFromFile(ctx.ReleaseNotes)
		if err != nil {
			return err
		}
		ctx.ReleaseNotes = notes
	}
	if ctx.Config.Changelog.Skip {
		return pipe.Skip("changelog should not be built")
	}
	if ctx.Snapshot {
		return pipe.Skip("not available for snapshots")
	}
	if ctx.ReleaseNotes != "" {
		return nil
	}
	if err := checkSortDirection(ctx.Config.Changelog.Sort); err != nil {
		return err
	}
	entries, err := buildChangelog(ctx)
	if err != nil {
		return err
	}

	changelogStringJoiner := "\n"
	if ctx.TokenType == context.TokenTypeGitLab || ctx.TokenType == context.TokenTypeGitea {
		// We need two or more whitespace to let markdown interpret
		// it as newline. See https://docs.gitlab.com/ee/user/markdown.html#newlines for details
		log.Debug("is gitlab or gitea changelog")
		changelogStringJoiner = "   \n"
	}
	ctx.ReleaseNotes = fmt.Sprintf("## Changelog\n\n%v\n", strings.Join(entries, changelogStringJoiner))
	var path = filepath.Join(ctx.Config.Dist, "CHANGELOG.md")
	log.WithField("changelog", path).Info("writing")
	return ioutil.WriteFile(path, []byte(ctx.ReleaseNotes), 0644)
}

func loadFromFile(file string) (string, error) {
	bts, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	log.WithField("file", file).Info("loaded custom release notes")
	log.WithField("file", file).Debugf("custom release notes: \n%s", string(bts))
	return string(bts), nil
}

func checkSortDirection(mode string) error {
	switch mode {
	case "":
		fallthrough
	case "asc":
		fallthrough
	case "desc":
		return nil
	}
	return ErrInvalidSortDirection
}

func buildChangelog(ctx *context.Context) ([]string, error) {
	log, err := getChangelog(ctx.Git.CurrentTag)
	if err != nil {
		return nil, err
	}
	var entries = strings.Split(log, "\n")
	entries = entries[0 : len(entries)-1]
	entries, err = filterEntries(ctx, entries)
	if err != nil {
		return entries, err
	}
	return sortEntries(ctx, entries), nil
}

func filterEntries(ctx *context.Context, entries []string) ([]string, error) {
	for _, filter := range ctx.Config.Changelog.Filters.Exclude {
		r, err := regexp.Compile(filter)
		if err != nil {
			return entries, err
		}
		entries = remove(r, entries)
	}
	return entries, nil
}

func sortEntries(ctx *context.Context, entries []string) []string {
	var direction = ctx.Config.Changelog.Sort
	if direction == "" {
		return entries
	}
	var result = make([]string, len(entries))
	copy(result, entries)
	sort.Slice(result, func(i, j int) bool {
		_, imsg := extractCommitInfo(result[i])
		_, jmsg := extractCommitInfo(result[j])
		if direction == "asc" {
			return strings.Compare(imsg, jmsg) < 0
		}
		return strings.Compare(imsg, jmsg) > 0
	})
	return result
}

func remove(filter *regexp.Regexp, entries []string) (result []string) {
	for _, entry := range entries {
		_, msg := extractCommitInfo(entry)
		if !filter.MatchString(msg) {
			result = append(result, entry)
		}
	}
	return result
}

func extractCommitInfo(line string) (hash, msg string) {
	ss := strings.Split(line, " ")
	return ss[0], strings.Join(ss[1:], " ")
}

func getChangelog(tag string) (string, error) {
	prev, err := previous(tag)
	if err != nil {
		return "", err
	}
	if isSHA1(prev) {
		return gitLog(prev, tag)
	}
	return gitLog(fmt.Sprintf("tags/%s..tags/%s", prev, tag))
}

func gitLog(refs ...string) (string, error) {
	var args = []string{"log", "--pretty=oneline", "--abbrev-commit", "--no-decorate", "--no-color"}
	args = append(args, refs...)
	return git.Run(args...)
}

func previous(tag string) (result string, err error) {
	result, err = git.Clean(git.Run("describe", "--tags", "--abbrev=0", fmt.Sprintf("tags/%s^", tag)))
	if err != nil {
		result, err = git.Clean(git.Run("rev-list", "--max-parents=0", "HEAD"))
	}
	return
}

// nolint: gochecknoglobals
var validSHA1 = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)

// isSHA1 te lets us know if the ref is a SHA1 or not
func isSHA1(ref string) bool {
	return validSHA1.MatchString(ref)
}
