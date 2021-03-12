//args: -Eerrcheck
//config: linters-settings.errcheck.check-type-assertions=true
package testdata

func ErrorTypeAssertion(filter map[string]interface{}) bool {
	return filter["messages_sent.messageid"].(map[string]interface{})["$ne"] != nil
}
