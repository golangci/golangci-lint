package printers

import (
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type html struct{}

func NewHTML() Printer {
	return &html{}
}

func (h html) Print(ctx context.Context, issues []result.Issue) error {
	type htmlIssue struct {
		Title  string
		Pos    string
		Linter string
		Code   string
	}
	htmlIssues := []*htmlIssue{}
	for i := range issues {
		pos := fmt.Sprintf("%s:%d", issues[i].FilePath(), issues[i].Line())
		if issues[i].Pos.Column != 0 {
			pos += fmt.Sprintf(":%d", issues[i].Pos.Column)
		}

		htmlIssues = append(htmlIssues, &htmlIssue{
			Title:  strings.TrimSpace(issues[i].Text),
			Pos:    pos,
			Linter: issues[i].FromLinter,
			Code:   strings.Join(issues[i].SourceLines, "\n"),
		})
	}
	t, err := template.New("golangci-lint").Parse(templateContent)
	if err != nil {
		return err
	}

	return t.Execute(logutils.StdOut, struct {
		Issues []*htmlIssue
	}{Issues: htmlIssues})
}

const templateContent = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>golangci-lint</title>
  <link rel="shortcut icon" type="image/png" href="https://golangci-lint.run/favicon-32x32.png">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.2/css/bulma.min.css" ` +
	`integrity="sha512-byErQdWdTqREz6DLAA9pCnLbdoGGhXfU6gm1c8bkf7F51JVmUBlayGe2A31VpXWQP+eiJ3ilTAZHCR3vmMyybA==" crossorigin="anonymous"/>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.2/styles/default.min.css" ` +
	`integrity="sha512-kZqGbhf9JTB4bVJ0G8HCkqmaPcRgo88F0dneK30yku5Y/dep7CZfCnNml2Je/sY4lBoqoksXz4PtVXS4GHSUzQ==" crossorigin="anonymous"/>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.2/highlight.min.js" integrity="sha512-s+tO` +
	`YYcC3Jybgr9mVsdAxsRYlGNq4mlAurOrfNuGMQ/SCofNPu92tjE7YRZCsdEtWL1yGkqk15fU/ark206YTg==" crossorigin="anonymous"></script>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.2/languages/go.min.js" integrity="sha512-+` +
	`UYV2NyyynWEQcZ4sMTKmeppyV331gqvMOGZ61/dqc89Tn1H40lF05ACd03RSD9EWwGutNwKj256mIR8waEJBQ==" crossorigin="anonymous"></script>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react.min.js" ` +
	`integrity="sha256-cLWs9L+cjZg8CjGHMpJqUgKKouPlmoMP/0wIdPtaPGs=" crossorigin="anonymous"></script>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react-dom.min.js" ` +
	`integrity="sha256-JIW8lNqN2EtqC6ggNZYnAdKMJXRQfkPMvdRt+b0/Jxc=" crossorigin="anonymous"></script>
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/babel-standalone/6.17.0/babel.min.js" ` +
	`integrity="sha256-1IWWLlCKFGFj/cjryvC7GDF5wRYnf9tSvNVVEj8Bm+o=" crossorigin="anonymous"></script>
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
    var Highlight = React.createClass({
      componentDidMount: function(){
        var current = ReactDOM.findDOMNode(this);
        hljs.highlightElement(current);
      },
      render: function() { 
        return (
          <pre className="go"><code>{ this.props.code }</code></pre>
        );
      }
    });
    var Issue = React.createClass({
      render: function() {
        return (
          <div className="issue box">
			<div>
              <div className="columns">
                <div className="column is-four-fifths">
                  <h5 className="title is-5 has-text-danger-dark">{ this.props.data.Title }</h5>
                </div>
			    <div className="column is-one-fifth">
                  <h6 className="title is-6">{ this.props.data.Linter }</h6>
                </div>
              </div>
              <strong>{ this.props.data.Pos }</strong>
		  	</div>
            <div className="highlight">
               <Highlight code={ this.props.data.Code }/>
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
                No issues found!
              </div>
            </div>
          );
        }
        var issues = this.props.data.Issues
          .map(function(issue) {
            return (<Issue data={issue} />);
          }.bind(this));
        return (
          <div className="issues">
            { issues }
          </div>
        );
      }
    });
    ReactDOM.render(
      <div className="content">
        <div className="columns is-centered">
          <div className="column is-three-quarters">
            <Issues data={ data } />
          </div>
        </div>
      </div>,
      document.getElementById("content")
    );
  </script>
</body>
</html>`
