//golangcitest:args -Eerrcheck
//golangcitest:config_path testdata/configs/errcheck_type_assertions.yml
package testdata

func ErrorTypeAssertion(filter map[string]interface{}) bool {
	return filter["messages_sent.messageid"].(map[string]interface{})["$ne"] != nil
}
