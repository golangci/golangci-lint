//golangcitest:args -Epathconcat
//golangcitest:config_path pathconcat_ignore_strings.yml
package testdata

func fqnAttr(ns, name string) string {
	return ns + "/attr/" + name // OK: contains ignored string
}

func fqnValue(ns, name, val string) string {
	return ns + "/attr/" + name + "/value/" + val // OK: contains ignored string
}
