package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/golangci/golangci-lint/internal/renameio"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
)

const listItemPrefix = "list-item-"

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
	log.Print("Successfully expanded templates")
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
	contentBytes, err := os.ReadFile(path)
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
	req, err := http.NewRequest( //nolint:noctx
		http.MethodGet,
		"https://api.github.com/repos/golangci/golangci-lint/releases/latest",
		http.NoBody,
	)
	if err != nil {
		return "", fmt.Errorf("failed to prepare a http request: %w", err)
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get http response for the latest tag: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read a body for the latest tag: %w", err)
	}
	release := latestRelease{}
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal the body for the latest tag: %w", err)
	}
	return release.TagName, nil
}

func buildTemplateContext() (map[string]string, error) {
	golangciYamlExample, err := os.ReadFile(".golangci.reference.yml")
	if err != nil {
		return nil, fmt.Errorf("can't read .golangci.reference.yml: %w", err)
	}

	snippets, err := extractExampleSnippets(golangciYamlExample)
	if err != nil {
		return nil, fmt.Errorf("can't read .golangci.reference.yml: %w", err)
	}

	if err = exec.Command("make", "build").Run(); err != nil {
		return nil, fmt.Errorf("can't run go install: %w", err)
	}

	lintersOut, err := exec.Command("./golangci-lint", "help", "linters").Output()
	if err != nil {
		return nil, fmt.Errorf("can't run linters cmd: %w", err)
	}

	lintersOutParts := bytes.Split(lintersOut, []byte("\n\n"))

	helpCmd := exec.Command("./golangci-lint", "run", "-h")
	helpCmd.Env = append(helpCmd.Env, os.Environ()...)
	helpCmd.Env = append(helpCmd.Env, "HELP_RUN=1") // make default concurrency stable: don't depend on machine CPU number
	help, err := helpCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("can't run help cmd: %w", err)
	}

	helpLines := bytes.Split(help, []byte("\n"))
	shortHelp := bytes.Join(helpLines[2:], []byte("\n"))
	changeLog, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		return nil, err
	}

	latestVersion, err := getLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest version: %w", err)
	}

	return map[string]string{
		"LintersExample":                   snippets.LintersSettings,
		"ConfigurationExample":             snippets.ConfigurationFile,
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
		"|Name|Description|Presets|AutoFix|Since|",
		"|---|---|---|---|---|---|",
	}

	for _, lc := range neededLcs {
		line := fmt.Sprintf("|%s|%s|%s|%v|%s|",
			getName(lc),
			getDesc(lc),
			strings.Join(lc.InPresets, ", "),
			check(lc.CanAutoFix, "Auto fix supported"),
			lc.Since,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func getName(lc *linter.Config) string {
	name := lc.Name()

	if lc.OriginalURL != "" {
		name = fmt.Sprintf("[%s](%s)", name, lc.OriginalURL)
	}

	if hasSettings(lc.Name()) {
		name = fmt.Sprintf("%s&nbsp;[%s](#%s)", name, spanWithID(listItemPrefix+lc.Name(), "Configuration", "⚙️"), lc.Name())
	}

	if !lc.IsDeprecated() {
		return name
	}

	title := "deprecated"
	if lc.Deprecation.Replacement != "" {
		title += fmt.Sprintf(" since %s", lc.Deprecation.Since)
	}

	return name + "&nbsp;" + span(title, "⚠")
}

func getDesc(lc *linter.Config) string {
	desc := lc.Linter.Desc()
	if lc.IsDeprecated() {
		desc = lc.Deprecation.Message
		if lc.Deprecation.Replacement != "" {
			desc += fmt.Sprintf(" Replaced by %s.", lc.Deprecation.Replacement)
		}
	}

	return strings.ReplaceAll(desc, "\n", "<br/>")
}

func check(b bool, title string) string {
	if b {
		return span(title, "✔")
	}
	return ""
}

func hasSettings(name string) bool {
	tp := reflect.TypeOf(config.LintersSettings{})

	for i := 0; i < tp.NumField(); i++ {
		if strings.EqualFold(name, tp.Field(i).Name) {
			return true
		}
	}

	return false
}

func span(title, icon string) string {
	return fmt.Sprintf(`<span title=%q>%s</span>`, title, icon)
}

func spanWithID(id, title, icon string) string {
	return fmt.Sprintf(`<span id=%q title=%q>%s</span>`, id, title, icon)
}

type authorDetails struct {
	Linters []string
	Profile string
	Avatar  string
}

func getThanksList() string {
	addedAuthors := map[string]*authorDetails{}

	for _, lc := range lintersdb.NewManager(nil, nil).GetAllSupportedLinterConfigs() {
		if lc.OriginalURL == "" {
			continue
		}

		linterURL := lc.OriginalURL
		if lc.Name() == "staticcheck" {
			linterURL = "https://github.com/dominikh/go-tools"
		}

		if author := extractAuthor(linterURL, "https://github.com/"); author != "" && author != "golangci" {
			if _, ok := addedAuthors[author]; ok {
				addedAuthors[author].Linters = append(addedAuthors[author].Linters, lc.Name())
			} else {
				addedAuthors[author] = &authorDetails{
					Linters: []string{lc.Name()},
					Profile: fmt.Sprintf("[%[1]s](https://github.com/sponsors/%[1]s)", author),
					Avatar:  fmt.Sprintf(`<img src="https://github.com/%[1]s.png" alt="%[1]s" style="max-width: 100%%;" width="20px;" />`, author),
				}
			}
		} else if author := extractAuthor(linterURL, "https://gitlab.com/"); author != "" {
			if _, ok := addedAuthors[author]; ok {
				addedAuthors[author].Linters = append(addedAuthors[author].Linters, lc.Name())
			} else {
				addedAuthors[author] = &authorDetails{
					Linters: []string{lc.Name()},
					Profile: fmt.Sprintf("[%[1]s](https://gitlab.com/%[1]s)", author),
				}
			}
		} else {
			continue
		}
	}

	var authors []string
	for author := range addedAuthors {
		authors = append(authors, author)
	}

	sort.Slice(authors, func(i, j int) bool {
		return strings.ToLower(authors[i]) < strings.ToLower(authors[j])
	})

	lines := []string{
		"|Author|Linter(s)|",
		"|---|---|",
	}

	for _, author := range authors {
		lines = append(lines, fmt.Sprintf("|%s %s|%s|",
			addedAuthors[author].Avatar, addedAuthors[author].Profile, strings.Join(addedAuthors[author].Linters, ", ")))
	}

	return strings.Join(lines, "\n")
}

func extractAuthor(originalURL, prefix string) string {
	if !strings.HasPrefix(originalURL, prefix) {
		return ""
	}

	return strings.SplitN(strings.TrimPrefix(originalURL, prefix), "/", 2)[0]
}

type SettingSnippets struct {
	ConfigurationFile string
	LintersSettings   string
}

func extractExampleSnippets(example []byte) (*SettingSnippets, error) {
	var data yaml.Node
	err := yaml.Unmarshal(example, &data)
	if err != nil {
		return nil, err
	}

	root := data.Content[0]

	globalNode := &yaml.Node{
		Kind:        root.Kind,
		Style:       root.Style,
		Tag:         root.Tag,
		Value:       root.Value,
		Anchor:      root.Anchor,
		Alias:       root.Alias,
		HeadComment: root.HeadComment,
		LineComment: root.LineComment,
		FootComment: root.FootComment,
		Line:        root.Line,
		Column:      root.Column,
	}

	snippets := SettingSnippets{}

	builder := strings.Builder{}

	for j, node := range root.Content {
		switch node.Value {
		case "run", "output", "linters", "linters-settings", "issues", "severity":
		default:
			continue
		}

		nextNode := root.Content[j+1]

		newNode := &yaml.Node{
			Kind: nextNode.Kind,
			Content: []*yaml.Node{
				{
					HeadComment: fmt.Sprintf("See the dedicated %q documentation section.", node.Value),
					Kind:        node.Kind,
					Style:       node.Style,
					Tag:         node.Tag,
					Value:       "option",
				},
				{
					Kind:  node.Kind,
					Style: node.Style,
					Tag:   node.Tag,
					Value: "value",
				},
			},
		}

		globalNode.Content = append(globalNode.Content, node, newNode)

		if node.Value == "linters-settings" {
			snippets.LintersSettings, err = getLintersSettingSnippets(node, nextNode)
			if err != nil {
				return nil, err
			}

			_, _ = builder.WriteString(
				fmt.Sprintf(
					"### `%s` configuration\n\nSee the dedicated [linters-settings](/usage/linters) documentation section.\n\n",
					node.Value,
				),
			)
			continue
		}

		nodeSection := &yaml.Node{
			Kind:    root.Kind,
			Style:   root.Style,
			Tag:     root.Tag,
			Value:   root.Value,
			Content: []*yaml.Node{node, nextNode},
		}

		snippet, errSnip := marshallSnippet(nodeSection)
		if errSnip != nil {
			return nil, errSnip
		}

		_, _ = builder.WriteString(fmt.Sprintf("### `%s` configuration\n\n%s", node.Value, snippet))
	}

	overview, err := marshallSnippet(globalNode)
	if err != nil {
		return nil, err
	}

	snippets.ConfigurationFile = overview + builder.String()

	return &snippets, nil
}

func getLintersSettingSnippets(node, nextNode *yaml.Node) (string, error) {
	builder := &strings.Builder{}

	for i := 0; i < len(nextNode.Content); i += 2 {
		r := &yaml.Node{
			Kind:  nextNode.Kind,
			Style: nextNode.Style,
			Tag:   nextNode.Tag,
			Value: node.Value,
			Content: []*yaml.Node{
				{
					Kind:  node.Kind,
					Value: node.Value,
				},
				{
					Kind:    nextNode.Kind,
					Content: []*yaml.Node{nextNode.Content[i], nextNode.Content[i+1]},
				},
			},
		}

		_, _ = fmt.Fprintf(builder, "### %s\n\n", nextNode.Content[i].Value)
		_, _ = fmt.Fprintln(builder, "```yaml")

		encoder := yaml.NewEncoder(builder)
		encoder.SetIndent(2)

		err := encoder.Encode(r)
		if err != nil {
			return "", err
		}

		_, _ = fmt.Fprintln(builder, "```")
		_, _ = fmt.Fprintln(builder)
		_, _ = fmt.Fprintf(builder, "[%s](#%s)\n\n", span("Back to the top", "🔼"), listItemPrefix+nextNode.Content[i].Value)
		_, _ = fmt.Fprintln(builder)
	}

	return builder.String(), nil
}

func marshallSnippet(node *yaml.Node) (string, error) {
	builder := &strings.Builder{}

	if node.Value != "" {
		_, _ = fmt.Fprintf(builder, "### %s\n\n", node.Value)
	}
	_, _ = fmt.Fprintln(builder, "```yaml")

	encoder := yaml.NewEncoder(builder)
	encoder.SetIndent(2)

	err := encoder.Encode(node)
	if err != nil {
		return "", err
	}

	_, _ = fmt.Fprintln(builder, "```")
	_, _ = fmt.Fprintln(builder)

	return builder.String(), nil
}
