package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golangci/golangci-lint/internal/renameio"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
)

var stateFilePath = filepath.Join("docs", "template_data.state")

func main() {
	var onlyWriteState bool
	flag.BoolVar(&onlyWriteState, "only-state", false, fmt.Sprintf("Only write hash of state to %s and exit", stateFilePath))
	flag.Parse()

	replacements, err := buildTemplateContext()
	if err != nil {
		log.Fatalf("Failed to build template context: %s", err)
	}

	if err = updateStateFile(replacements); err != nil {
		log.Fatalf("Failed to update state file: %s", err)
	}

	if onlyWriteState {
		return
	}

	if err := rewriteDocs(replacements); err != nil {
		log.Fatalf("Failed to rewrite docs: %s", err)
	}
	log.Printf("Successfully expanded templates")
}

func updateStateFile(replacements map[string]string) error {
	replBytes, err := json.Marshal(replacements)
	if err != nil {
		return fmt.Errorf("failed to json marshal replacements: %w", err)
	}

	h := sha256.New()
	if _, err := h.Write(replBytes); err != nil {
		return err
	}

	var contentBuf bytes.Buffer
	contentBuf.WriteString("This file stores hash of website templates to trigger " +
		"Netlify rebuild when something changes, e.g. new linter is added.\n")
	contentBuf.WriteString(hex.EncodeToString(h.Sum(nil)))

	return renameio.WriteFile(stateFilePath, contentBuf.Bytes(), os.ModePerm)
}

func rewriteDocs(replacements map[string]string) error {
	madeReplacements := map[string]bool{}
	err := filepath.Walk(filepath.Join("docs", "src", "docs"),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			return processDoc(path, replacements, madeReplacements)
		})
	if err != nil {
		return fmt.Errorf("failed to walk dir: %w", err)
	}

	if len(madeReplacements) != len(replacements) {
		for key := range replacements {
			if !madeReplacements[key] {
				log.Printf("Replacement %q wasn't performed", key)
			}
		}
		return fmt.Errorf("%d replacements weren't performed", len(replacements)-len(madeReplacements))
	}
	return nil
}

func processDoc(path string, replacements map[string]string, madeReplacements map[string]bool) error {
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	content := string(contentBytes)
	hasReplacements := false
	for key, replacement := range replacements {
		nextContent := content
		nextContent = strings.ReplaceAll(nextContent, fmt.Sprintf("{.%s}", key), replacement)

		// Yaml formatter in mdx code section makes extra spaces, need to match them too.
		nextContent = strings.ReplaceAll(nextContent, fmt.Sprintf("{ .%s }", key), replacement)

		if nextContent != content {
			hasReplacements = true
			madeReplacements[key] = true
			content = nextContent
		}
	}
	if !hasReplacements {
		return nil
	}

	log.Printf("Expanded template in %s, saving it", path)
	if err = renameio.WriteFile(path, []byte(content), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write changes to file %s: %w", path, err)
	}

	return nil
}

type latestRelease struct {
	TagName string `json:"tag_name"`
}

func getLatestVersion() (string, error) {
	req, err := http.NewRequest( // nolint:noctx
		http.MethodGet,
		"https://api.github.com/repos/golangci/golangci-lint/releases/latest",
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to prepare a http request: %s", err)
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get http response for the latest tag: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read a body for the latest tag: %s", err)
	}
	release := latestRelease{}
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal the body for the latest tag: %s", err)
	}
	return release.TagName, nil
}

func buildTemplateContext() (map[string]string, error) {
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

	return map[string]string{
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
	lcs := lintersdb.NewManager(nil, nil).GetAllSupportedLinterConfigs()
	for _, lc := range lcs {
		if lc.EnabledByDefault == enabled {
			neededLcs = append(neededLcs, lc)
		}
	}

	sort.Slice(neededLcs, func(i, j int) bool {
		return neededLcs[i].Name() < neededLcs[j].Name()
	})

	lines := []string{
		"|Name|Description|Presets|AutoFix|Deprecated|",
		"|---|---|---|---|---|",
	}

	for _, lc := range neededLcs {
		line := fmt.Sprintf("|%s|%s|%s|%v|%s|",
			getName(lc),
			getDesc(lc),
			strings.Join(lc.InPresets, ", "),
			check(lc.CanAutoFix, "Auto fix supported"),
			check(lc.DeprecatedMessage != "", "Deprecated"),
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func getName(lc *linter.Config) string {
	name := lc.Name()

	if lc.OriginalURL != "" {
		name = fmt.Sprintf("[%s](%s)", lc.Name(), lc.OriginalURL)
	}

	if lc.DeprecatedMessage != "" {
		name += ` <span title="deprecated">⚠</span>`
	}

	return name
}

func getDesc(lc *linter.Config) string {
	desc := lc.Linter.Desc()
	if lc.DeprecatedMessage != "" {
		desc = lc.DeprecatedMessage
	}

	return strings.ReplaceAll(desc, "\n", "<br/>")
}

func check(b bool, title string) string {
	if b {
		return `<span title="` + title + `">✔</span>`
	}
	return ""
}

func getThanksList() string {
	var lines []string
	addedAuthors := map[string]bool{}
	for _, lc := range lintersdb.NewManager(nil, nil).GetAllSupportedLinterConfigs() {
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
