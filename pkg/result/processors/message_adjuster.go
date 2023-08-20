package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/result"
)

type MessageRule struct {
	Linter          string
	ExistingMessage string
	NewMessage      string
}

type messageRule struct {
	linter          string
	existingMessage *regexp.Regexp
	newMessage      string
}

type MessageAdjuster struct {
	messageRules []messageRule
}

func NewMessageAdjuster(msgRules []MessageRule) *MessageAdjuster {
	parsedRules := make([]messageRule, 0, len(msgRules))
	for _, rule := range msgRules {
		parsedRules = append(parsedRules, messageRule{
			linter:          rule.Linter,
			existingMessage: regexp.MustCompile(rule.ExistingMessage),
			newMessage:      rule.NewMessage,
		})
	}

	return &MessageAdjuster{
		messageRules: parsedRules,
	}
}

func (ma MessageAdjuster) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(ma.messageRules) == 0 {
		return issues, nil
	}

	return transformIssues(issues, func(issue *result.Issue) *result.Issue {
		for _, msgRule := range ma.messageRules {
			if msgRule.linter != "" && msgRule.linter != issue.FromLinter {
				return issue
			}

			existingText := issue.Text
			newText := msgRule.existingMessage.ReplaceAllString(existingText, msgRule.newMessage)
			if existingText != newText {
				issue.Text = newText
				return issue
			}
		}

		return issue
	}), nil
}

func (ma MessageAdjuster) Name() string {
	return "message_adjuster"
}
func (ma MessageAdjuster) Finish() {}
