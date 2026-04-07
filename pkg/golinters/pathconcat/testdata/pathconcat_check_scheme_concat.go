//golangcitest:args -Epathconcat
//golangcitest:config_path pathconcat_check_scheme_concat.yml
package testdata

func schemeConcat(host string) string {
	return "https://" + host // want `use url\.JoinPath\(\) instead of string concatenation with "/"`
}
