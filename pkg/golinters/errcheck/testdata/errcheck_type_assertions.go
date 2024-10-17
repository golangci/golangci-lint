//golangcitest:args -Eerrcheck
//golangcitest:config_path testdata/errcheck_type_assertions.yml
//golangcitest:expected_exitcode 1
package testdata

func ErrorTypeAssertion(filter map[string]interface{}) bool {
	return filter["messages_sent.messageid"].(map[string]interface{})["$ne"] != nil // want "Error return value is not checked"
}
