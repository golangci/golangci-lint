package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
)

func main() {
	const (
		tmplPath = "README.md.tmpl"
		outPath  = "README.md"
	)

	if err := genReadme(tmplPath, outPath); err != nil {
		log.Fatalf("failed: %s", err)
	}
	log.Printf("Successfully generated %s", outPath)
}

func genReadme(tmplPath, outPath string) error {
	ctx, err := buildTemplateContext()
	if err != nil {
		return err
	}

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	tmpl := template.Must(template.ParseFiles(tmplPath))
	return tmpl.Execute(out, ctx)
}

func buildTemplateContext() (map[string]interface{}, error) {
	golangciYaml, err := ioutil.ReadFile(".golangci.yml")
	if err != nil {
		return nil, fmt.Errorf("can't read .golangci.yml: %s", err)
	}

	if err = exec.Command("go", "install", "./cmd/...").Run(); err != nil {
		return nil, fmt.Errorf("can't run go install: %s", err)
	}

	lintersOut, err := exec.Command("golangci-lint", "linters").Output()
	if err != nil {
		return nil, fmt.Errorf("can't run linters cmd: %s", err)
	}

	lintersOutParts := bytes.Split(lintersOut, []byte("\n\n"))

	return map[string]interface{}{
		"GolangciYaml":                     string(golangciYaml),
		"LintersCommandOutputEnabledOnly":  string(lintersOutParts[0]),
		"LintersCommandOutputDisabledOnly": string(lintersOutParts[1]),
		"EnabledByDefaultLinters":          getLintersListMarkdown(true),
		"DisabledByDefaultLinters":         getLintersListMarkdown(false),
		"ThanksList":                       getThanksList(),
	}, nil
}

func getLintersListMarkdown(enabled bool) string {
	var neededLcs []linter.Config
	lcs := lintersdb.GetAllSupportedLinterConfigs()
	for _, lc := range lcs {
		if lc.EnabledByDefault == enabled {
			neededLcs = append(neededLcs, lc)
		}
	}

	var lines []string
	for _, lc := range neededLcs {
		var link string
		if lc.OriginalURL != "" {
			link = fmt.Sprintf("[%s](%s)", lc.Linter.Name(), lc.OriginalURL)
		} else {
			link = lc.Linter.Name()
		}
		line := fmt.Sprintf("- %s - %s", link, lc.Linter.Desc())
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func getThanksList() string {
	var lines []string
	addedAuthors := map[string]bool{}
	for _, lc := range lintersdb.GetAllSupportedLinterConfigs() {
		if lc.OriginalURL == "" {
			continue
		}

		const githubPrefix = "https://github.com/"
		if !strings.HasPrefix(lc.OriginalURL, githubPrefix) {
			continue
		}

		githubSuffix := strings.TrimPrefix(lc.OriginalURL, githubPrefix)
		githubAuthor := strings.Split(githubSuffix, "/")[0]
		if addedAuthors[githubAuthor] {
			continue
		}
		addedAuthors[githubAuthor] = true

		line := fmt.Sprintf("- [%s](https://github.com/%s)",
			githubAuthor, githubAuthor)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
