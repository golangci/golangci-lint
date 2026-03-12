package sarif

// Level SARIF level
// From https://docs.oasis-open.org/sarif/sarif/v2.0/csprd02/sarif-v2.0-csprd02.html#_Toc10127839
type Level string

const (
	// None : The concept of “severity” does not apply to this result because the kind
	// property (§3.27.9) has a value other than "fail".
	None = Level("none")
	// Note : The rule specified by ruleId was evaluated and a minor problem or an opportunity
	// to improve the code was found.
	Note = Level("note")
	// Warning : The rule specified by ruleId was evaluated and a problem was found.
	Warning = Level("warning")
	// Error : The rule specified by ruleId was evaluated and a serious problem was found.
	Error = Level("error")
	// Version : SARIF Schema version
	Version = "2.1.0"
	// Schema : SARIF Schema URL
	Schema = "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/main/sarif-2.1/schema/sarif-schema-2.1.0.json"
)
