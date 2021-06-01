package printers

import (
	"context"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/securego/gosec/v2/cwe"
	"html/template"
	"os"
	"strconv"
)

type ReportInfo struct {
	Issues []*HtmlIssue
}

// HtmlIssue is a subset of the Code Climate spec - https://github.com/codeclimate/spec/blob/master/SPEC.md#data-types
// It is just enough to support GitLab CI Code Quality - https://docs.gitlab.com/ee/user/project/merge_requests/code_quality.html
type HtmlIssue struct {
	Severity   string        `json:"severity"`   // issue severity (how problematic it is)
	Confidence string        `json:"confidence"` // issue confidence (how sure we are we found it)
	Cwe        *cwe.Weakness `json:"cwe"`        // Cwe associated with RuleID
	RuleID     string        `json:"rule_id"`    // Human readable explanation
	What       string        `json:"details"`    // Human readable explanation
	File       string        `json:"file"`       // File name we found it in
	Code       string        `json:"code"`       // Impacted code line
	Line       string        `json:"line"`       // Line number in file
	Col        string        `json:"column"`     // Column number in line
}

type Html struct{}

func NewHtml() *Html {
	return &Html{}
}

func (h Html) Print(ctx context.Context, issues []result.Issue) error {
	htmlIssues := []*HtmlIssue{}
	for i := range issues {
		issue := &issues[i]
		htmlIssue := HtmlIssue{}
		htmlIssue.What = issue.Description()
		htmlIssue.File = issue.Pos.Filename
		htmlIssue.Line = strconv.Itoa(issue.Pos.Line)

		if issue.Severity != "" {
			htmlIssue.Severity = issue.Severity
		}

		htmlIssues = append(htmlIssues, &htmlIssue)
	}

	t, err := template.New("golangci-lint").Parse(templateContent)
	if err != nil {
		return err
	}

	file, err := os.OpenFile("golangci-lint.html", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	return t.Execute(file, ReportInfo{Issues: htmlIssues})
}

const templateContent = `
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Golang Security Checker</title>
  <link rel="shortcut icon" type="image/png" href="https://securego.io/img/favicon.png">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.2/css/bulma.min.css" integrity="sha512-byErQdWdTqREz6DLAA9pCnLbdoGGhXfU6gm1c8bkf7F51JVmUBlayGe2A31VpXWQP+eiJ3ilTAZHCR3vmMyybA==" crossorigin="anonymous"/>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.2/styles/default.min.css" integrity="sha512-kZqGbhf9JTB4bVJ0G8HCkqmaPcRgo88F0dneK30yku5Y/dep7CZfCnNml2Je/sY4lBoqoksXz4PtVXS4GHSUzQ==" crossorigin="anonymous"/>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.2/highlight.min.js" integrity="sha512-s+tOYYcC3Jybgr9mVsdAxsRYlGNq4mlAurOrfNuGMQ/SCofNPu92tjE7YRZCsdEtWL1yGkqk15fU/ark206YTg==" crossorigin="anonymous"></script>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.2/languages/go.min.js" integrity="sha512-+UYV2NyyynWEQcZ4sMTKmeppyV331gqvMOGZ61/dqc89Tn1H40lF05ACd03RSD9EWwGutNwKj256mIR8waEJBQ==" crossorigin="anonymous"></script>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react.min.js" integrity="sha256-cLWs9L+cjZg8CjGHMpJqUgKKouPlmoMP/0wIdPtaPGs=" crossorigin="anonymous"></script>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react-dom.min.js" integrity="sha256-JIW8lNqN2EtqC6ggNZYnAdKMJXRQfkPMvdRt+b0/Jxc=" crossorigin="anonymous"></script>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/babel-standalone/6.17.0/babel.min.js" integrity="sha256-1IWWLlCKFGFj/cjryvC7GDF5wRYnf9tSvNVVEj8Bm+o=" crossorigin="anonymous"></script>
  <style>
  .field-label {
    min-width: 80px;
  }
  .break-word {
    word-wrap: break-word;
  }
  .help {
    white-space: pre-wrap;
  }
  .tag {
    width: 80px;
  }
  </style>
</head>
<body>
  <section class="section">
    <div class="container">
      <div id="content"></div>
    </div>
  </section>
  <script>
    var data = {{ . }};
  </script>
  <script type="text/babel">
    var IssueTag = React.createClass({
      render: function() {
        var level = "tag"
        if (this.props.level === "error") {
          level += " is-danger";
        } else if (this.props.level === "info") {
          level += " is-info";
        }
        level +=" is-rounded";
        return (
          <div className="control">
            <div className="tags has-addons">
              <span className="tag is-dark is-rounded">{ this.props.label }</span>
              <span className={ level }>{ this.props.level }</span>
            </div>
          </div>
        );
      }
    });
    var Highlight = React.createClass({
      componentDidMount: function(){
        var current = ReactDOM.findDOMNode(this);
        hljs.highlightElement(current);
      },
      render: function() { 
        return (
          <pre className="go"><code >{ this.props.code }</code></pre>
        );
      }
    });
    var Issue = React.createClass({
      render: function() {
        return (
          <div className="issue box">
          <div className="columns">
              <div className="column is-three-quarters">
                <strong className="break-word">{ this.props.data.file } (line { this.props.data.line })</strong>
                <p>{ this.props.data.details }</p>
              </div>
              <div className="column is-one-quarter">
                <div className="field is-grouped is-grouped-multiline">
					{console.log(this.props.data.severity)}
                  <IssueTag label="Severity" level={ this.props.data.severity }/>
                </div>
              </div>
            </div>
            <div className="highlight">
              <Highlight code={ this.props.data.what }/>
            </div>
          </div>
        );
      }
    });
    var Issues = React.createClass({
      render: function() {
        if (this.props.data.Issues.length === 0) {
          return (
            <div>
              <div className="notification">
                Awesome! No issues found!
              </div>
            </div>
          );
        }
        var issues = this.props.data.Issues
          .filter(function(issue) {
            return this.props.severity.includes(issue.severity);
          }.bind(this))
          .filter(function(issue) {
            return this.props.confidence.includes(issue.confidence);
          }.bind(this))
          .filter(function(issue) {
            if (this.props.issueType) {
              return issue.details.toLowerCase().startsWith(this.props.issueType.toLowerCase());
            } else {
              return true
            }
          }.bind(this))
          .map(function(issue) {
            return (<Issue data={issue} />);
          }.bind(this));
        if (issues.length === 0) {
          return (
            <div>
              <div className="notification">
                No issues matched given filters
                (of total { this.props.data.Issues.length } issues).
              </div>
            </div>
          );
        }
        return (
          <div className="issues">
            { issues }
          </div>
        );
      }
    });
    var LevelSelector = React.createClass({
      handleChange: function(level) {
        return function(e) {
          var updated = this.props.selected
            .filter(function(item) { return item != level; });
          if (e.target.checked) {
            updated.push(level);
          }
          this.props.onChange(updated);
        }.bind(this);
      },
      render: function() {
        var HIGH = "error", MEDIUM = "MEDIUM", LOW = "info";
        var highDisabled = !this.props.available.includes(HIGH);
        var mediumDisabled = !this.props.available.includes(MEDIUM);
        var lowDisabled = !this.props.available.includes(LOW);
        return (
          <div className="field">
            <div className="control">
              <label className="checkbox" disabled={ highDisabled }>
                <input
                  type="checkbox"
                  checked={ this.props.selected.includes(HIGH) }
                  disabled={ highDisabled }
                  onChange={ this.handleChange(HIGH) }/> High
              </label>
            </div>
            <div className="control">
              <label className="checkbox" disabled={ mediumDisabled }>
                <input
                  type="checkbox"
                  checked={ this.props.selected.includes(MEDIUM) }
                  disabled={ mediumDisabled }
                  onChange={ this.handleChange(MEDIUM) }/> Medium
              </label>
            </div>
            <div className="control">
              <label className="checkbox" disabled={ lowDisabled }>
                <input
                  type="checkbox"
                  checked={ this.props.selected.includes(LOW) }
                  disabled={ lowDisabled }
                  onChange={ this.handleChange(LOW) }/> Low
              </label>
            </div>
          </div>
        );
      }
    });
    var Navigation = React.createClass({
      updateSeverity: function(vals) {
        this.props.onSeverity(vals);
      },
      updateConfidence: function(vals) {
        this.props.onConfidence(vals);
      },
      updateIssueType: function(e) {
        if (e.target.value == "all") {
          this.props.onIssueType(null);
        } else {
          this.props.onIssueType(e.target.value);
        }
      },
      render: function() {
        var issueTypes = this.props.allIssueTypes
          .map(function(it) {
            var matches = this.props.issueType == it
            return (
              <option value={ it } selected={ matches }>
                { it }
              </option>
            );
          }.bind(this));
        return (
          <nav className="panel">
            <div className="panel-heading">Filters</div>
            <div className="panel-block">
              <div className="field is-horizontal">
                <div className="field-label is-normal">
                  <label className="label is-pulled-left">Severity</label>
                </div>
                <div className="field-body">
                  <LevelSelector selected={ this.props.severity } available={ this.props.allSeverities } onChange={ this.updateSeverity } />
                </div>
             </div>
            </div>
            <div className="panel-block">
              <div className="field is-horizontal">
                <div className="field-label is-normal">
                  <label className="label is-pulled-left">Confidence</label>
                </div>
                <div className="field-body">
                  <LevelSelector selected={ this.props.confidence } available={ this.props.allConfidences } onChange={ this.updateConfidence } />
                </div>
              </div>
            </div>
            <div className="panel-block">
              <div className="field is-horizontal">
                <div className="field-label is-normal">
                  <label className="label is-pulled-left">Issue type</label>
                </div>
                <div className="field-body">
                  <div className="field">
                    <div className="control">
                      <div className="select is-fullwidth">
                        <select onChange={ this.updateIssueType }>
                          <option value="all" selected={ !this.props.issueType }>
                            (all)
                          </option>
                          { issueTypes }
                        </select>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </nav>
        );
      }
    });
    var IssueBrowser = React.createClass({
      getInitialState: function() {
        return {};
      },
      componentWillMount: function() {
        this.updateIssues(this.props.data);
      },
      handleSeverity: function(val) {
        this.updateIssueTypes(this.props.data.Issues, val, this.state.confidence);
        this.setState({severity: val});
      },
      handleConfidence: function(val) {
        this.updateIssueTypes(this.props.data.Issues, this.state.severity, val);
        this.setState({confidence: val});
      },
      handleIssueType: function(val) {
        this.setState({issueType: val});
      },
      updateIssues: function(data) {
        if (!data) {
          this.setState({data: data});
          return;
        }
        var allSeverities = data.Issues
          .map(function(issue) {
            return issue.severity
          })
          .sort()
          .filter(function(item, pos, ary) {
            return !pos || item != ary[pos - 1];
          });
        var allConfidences = data.Issues
          .map(function(issue) {
            return issue.confidence
          })
          .sort()
          .filter(function(item, pos, ary) {
            return !pos || item != ary[pos - 1];
          });
        var selectedSeverities = allSeverities;
        var selectedConfidences = allConfidences;
        this.updateIssueTypes(data.Issues, selectedSeverities, selectedConfidences);
        this.setState({
          data: data,
          severity: selectedSeverities,
          allSeverities: allSeverities,
          confidence: selectedConfidences,
          allConfidences: allConfidences,
          issueType: null
        });
      },
      updateIssueTypes: function(issues, severities, confidences) {
        var allTypes = issues
          .filter(function(issue) {
            return severities.includes(issue.severity);
          })
          .filter(function(issue) {
            return confidences.includes(issue.confidence);
          })
          .map(function(issue) {
            return issue.details;
          })
          .sort()
          .filter(function(item, pos, ary) {
            return !pos || item != ary[pos - 1];
          });
        if (this.state.issueType && !allTypes.includes(this.state.issueType)) {
          this.setState({issueType: null});
        }
        this.setState({allIssueTypes: allTypes});
      },
      render: function() {
        return (
          <div className="content">
            <div className="columns">
              <div className="column is-one-quarter">
                <Navigation
                  severity={ this.state.severity } 
                  confidence={ this.state.confidence }
                  issueType={ this.state.issueType }
                  allSeverities={ this.state.allSeverities } 
                  allConfidences={ this.state.allConfidences }
                  allIssueTypes={ this.state.allIssueTypes }
                  onSeverity={ this.handleSeverity } 
                  onConfidence={ this.handleConfidence } 
                  onIssueType={ this.handleIssueType }
                />
              </div>
              <div className="column is-three-quarters">
                <Issues
                  data={ this.props.data }
                  severity={ this.state.severity }
                  confidence={ this.state.confidence }
                  issueType={ this.state.issueType }
                />
              </div>
            </div>
          </div>
        );
      }
    });
    ReactDOM.render(
      <IssueBrowser data={ data } />,
      document.getElementById("content")
    );
  </script>
</body>
</html>`
