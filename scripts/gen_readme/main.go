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
		tmplPath = "README.tmpl.md"
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

func getLatestVersion() (string, error) {
	out, err := exec.Command("git", "tag", "-l", "--sort=-v:refname").Output()
	if err != nil {
		return "", fmt.Errorf("failed to run git tag: %s", err)
	}

	lines := bytes.Split(out, []byte("\n"))
	return string(lines[0]), nil
}

func buildTemplateContext() (map[string]interface{}, error) {
	golangciYaml, err := ioutil.ReadFile(".golangci.yml")
	if err != nil {
		return nil, fmt.Errorf("can't read .golangci.yml: %s", err)
	}

	golangciYamlExample, err := ioutil.ReadFile(".golangci.example.yml")
	if err != nil {
		return nil, fmt.Errorf("can't read .golangci.example.yml: %s", err)
	}

	if err = exec.Command("make", "build").Run(); err != nil {
		return nil, fmt.Errorf("can't run go install: %s", err)
	}

	lintersOut, err := exec.Command("./golangci-lint", "help", "linters").Output()
	if err != nil {
		return nil, fmt.Errorf("can't run linters cmd: %s", err)
	}

	lintersOutParts := bytes.Split(lintersOut, []byte("\n\n"))

	helpCmd := exec.Command("./golangci-lint", "run", "-h")
	helpCmd.Env = append(helpCmd.Env, os.Environ()...)
	helpCmd.Env = append(helpCmd.Env, "HELP_RUN=1") // make default concurrency stable: don't depend on machine CPU number
	help, err := helpCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("can't run help cmd: %s", err)
	}

	helpLines := bytes.Split(help, []byte("\n"))
	shortHelp := bytes.Join(helpLines[2:], []byte("\n"))
	changeLog, err := ioutil.ReadFile("CHANGELOG.md")
	if err != nil {
		return nil, err
	}

	latestVersion, err := getLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %s", err)
	}

	return map[string]interface{}{
		"GolangciYaml":                     strings.TrimSpace(string(golangciYaml)),
		"GolangciYamlExample":              strings.TrimSpace(string(golangciYamlExample)),
		"LintersCommandOutputEnabledOnly":  string(lintersOutParts[0]),
		"LintersCommandOutputDisabledOnly": string(lintersOutParts[1]),
		"EnabledByDefaultLinters":          getLintersListMarkdown(true),
		"DisabledByDefaultLinters":         getLintersListMarkdown(false),
		"ThanksList":                       getThanksList(),
		"RunHelpText":                      string(shortHelp),
		"ChangeLog":                        string(changeLog),
		"LatestVersion":                    latestVersion,
	}, nil
}

func getLintersListMarkdown(enabled bool) string {
	var neededLcs []*linter.Config
	lcs := lintersdb.NewManager(nil).GetAllSupportedLinterConfigs()
	for _, lc := range lcs {
		if lc.EnabledByDefault == enabled {
			neededLcs = append(neededLcs, lc)
		}
	}

	var lines []string
	for _, lc := range neededLcs {
		var link string
		if lc.OriginalURL != "" {
			link = fmt.Sprintf("[%s](%s)", lc.Name(), lc.OriginalURL)
		} else {
			link = lc.Name()
		}
		line := fmt.Sprintf("- %s - %s", link, lc.Linter.Desc())
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func getThanksList() string {
	var lines []string
	addedAuthors := map[string]bool{}
	for _, lc := range lintersdb.NewManager(nil).GetAllSupportedLinterConfigs() {
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
