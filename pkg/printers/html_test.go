package printers

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

const expectedHTML = `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>golangci-lint</title>
    <link rel="shortcut icon" type="image/png" href="https://golangci-lint.run/favicon-32x32.png">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.2/css/bulma.min.css"
          integrity="sha512-byErQdWdTqREz6DLAA9pCnLbdoGGhXfU6gm1c8bkf7F51JVmUBlayGe2A31VpXWQP+eiJ3ilTAZHCR3vmMyybA=="
          crossorigin="anonymous" referrerpolicy="no-referrer"/>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.2/styles/default.min.css"
          integrity="sha512-kZqGbhf9JTB4bVJ0G8HCkqmaPcRgo88F0dneK30yku5Y/dep7CZfCnNml2Je/sY4lBoqoksXz4PtVXS4GHSUzQ=="
          crossorigin="anonymous" referrerpolicy="no-referrer"/>
    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.2/highlight.min.js"
            integrity="sha512-s+tOYYcC3Jybgr9mVsdAxsRYlGNq4mlAurOrfNuGMQ/SCofNPu92tjE7YRZCsdEtWL1yGkqk15fU/ark206YTg=="
            crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.2/languages/go.min.js"
            integrity="sha512-+UYV2NyyynWEQcZ4sMTKmeppyV331gqvMOGZ61/dqc89Tn1H40lF05ACd03RSD9EWwGutNwKj256mIR8waEJBQ=="
            crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/react/17.0.2/umd/react.production.min.js"
            integrity="sha512-qlzIeUtTg7eBpmEaS12NZgxz52YYZVF5myj89mjJEesBd/oE9UPsYOX2QAXzvOAZYEvQohKdcY8zKE02ifXDmA=="
            crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <script type="text/javascript"
            src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/17.0.2/umd/react-dom.production.min.js"
            integrity="sha512-9jGNr5Piwe8nzLLYTk8QrEMPfjGU0px80GYzKZUxi7lmCfrBjtyCc1V5kkS5vxVwwIB7Qpzc7UxLiQxfAN30dw=="
            crossorigin="anonymous" referrerpolicy="no-referrer"></script>
    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/babel-standalone/6.26.0/babel.min.js"
            integrity="sha512-kp7YHLxuJDJcOzStgd6vtpxr4ZU9kjn77e6dBsivSz+pUuAuMlE2UTdKB7jjsWT84qbS8kdCWHPETnP/ctrFsA=="
            crossorigin="anonymous" referrerpolicy="no-referrer"></script>
</head>
<body>
<section class="section">
    <div class="container">
        <div id="content"></div>
    </div>
</section>
<script>
    const data = {"Issues":[{"Title":"some issue","Pos":"path/to/filea.go:10:4","Linter":"linter-a","Code":""},{"Title":"another issue","Pos":"path/to/fileb.go:300:9","Linter":"linter-b","Code":"func foo() {\n\tfmt.Println(\"bar\")\n}"}]};
</script>
<script type="text/babel">
  class Highlight extends React.Component {
    componentDidMount() {
      hljs.highlightElement(ReactDOM.findDOMNode(this));
    }

    render() {
      return <pre className="go"><code>{this.props.code}</code></pre>;
    }
  }

  class Issue extends React.Component {
    render() {
      return (
        <div className="issue box">
          <div>
            <div className="columns">
              <div className="column is-four-fifths">
                <h5 className="title is-5 has-text-danger-dark">{this.props.data.Title}</h5>
              </div>
              <div className="column is-one-fifth">
                <h6 className="title is-6">{this.props.data.Linter}</h6>
              </div>
            </div>
            <strong>{this.props.data.Pos}</strong>
          </div>
          <div className="highlight">
            <Highlight code={this.props.data.Code}/>
          </div>
        </div>
      );
    }
  }

  class Issues extends React.Component {
    render() {
      if (!this.props.data.Issues || this.props.data.Issues.length === 0) {
        return (
          <div>
            <div className="notification">
              No issues found!
            </div>
          </div>
        );
      }

      return (
        <div className="issues">
          {this.props.data.Issues.map(issue => (<Issue data={issue}/>))}
        </div>
      );
    }
  }

  ReactDOM.render(
    <div className="content">
      <div className="columns is-centered">
        <div className="column is-three-quarters">
          <Issues data={data}/>
        </div>
      </div>
    </div>,
    document.getElementById("content")
  );
</script>
</body>
</html>`

func TestHTML_Print(t *testing.T) {
	issues := []result.Issue{
		{
			FromLinter: "linter-a",
			Severity:   "warning",
			Text:       "some issue",
			Pos: token.Position{
				Filename: "path/to/filea.go",
				Offset:   2,
				Line:     10,
				Column:   4,
			},
		},
		{
			FromLinter: "linter-b",
			Severity:   "error",
			Text:       "another issue",
			SourceLines: []string{
				"func foo() {",
				"\tfmt.Println(\"bar\")",
				"}",
			},
			Pos: token.Position{
				Filename: "path/to/fileb.go",
				Offset:   5,
				Line:     300,
				Column:   9,
			},
		},
	}

	buf := new(bytes.Buffer)
	printer := NewHTML(buf)

	err := printer.Print(issues)
	require.NoError(t, err)

	assert.Equal(t, expectedHTML, buf.String())
}
