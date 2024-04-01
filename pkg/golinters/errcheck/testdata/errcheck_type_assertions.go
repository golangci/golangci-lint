//golangcitest:args -Eerrcheck
//golangcitest:config_path testdata/errcheck_type_assertions.yml
//golangcitest:expected_exitcode 0
package testdata

func ErrorTypeAssertion(filter map[string]interface{}) bool {
	return filter["messages_sent.messageid"].(map[string]interface{})["$ne"] != nil
}
