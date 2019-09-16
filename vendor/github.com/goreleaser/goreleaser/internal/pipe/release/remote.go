package release

import (
	"fmt"
	"strings"

	"github.com/goreleaser/goreleaser/internal/git"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/pkg/errors"
)

// remoteRepo gets the repo name from the Git config.
func remoteRepo() (result config.Repo, err error) {
	if !git.IsRepo() {
		return result, errors.New("current folder is not a git repository")
	}
	out, err := git.Run("config", "--get", "remote.origin.url")
	if err != nil {
		return result, fmt.Errorf("repository doesn't have an `origin` remote")
	}
	return extractRepoFromURL(out), nil
}

func extractRepoFromURL(s string) config.Repo {
	// removes the .git suffix and any new lines
	s = strings.NewReplacer(
		".git", "",
		"\n", "",
	).Replace(s)
	// if the URL contains a :, indicating a SSH config,
	// remove all chars until it, including itself
	// on HTTP and HTTPS URLs it will remove the http(s): prefix,
	// which is ok. On SSH URLs the whole user@server will be removed,
	// which is required.
	s = s[strings.LastIndex(s, ":")+1:]
	// split by /, the last to parts should be the owner and name
	ss := strings.Split(s, "/")
	return config.Repo{
		Owner: ss[len(ss)-2],
		Name:  ss[len(ss)-1],
	}
}
