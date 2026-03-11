package junit

// NewTestsuite instantiate a Testsuite
func NewTestsuite(name string) *Testsuite {
	return &Testsuite{
		Name: name,
	}
}

// NewFailure instantiate a Failure
func NewFailure(message string, text string) *Failure {
	return &Failure{
		Message: message,
		Text:    text,
	}
}

// NewTestcase instantiate a Testcase
func NewTestcase(name string, failure *Failure) *Testcase {
	return &Testcase{
		Name:    name,
		Failure: failure,
	}
}
