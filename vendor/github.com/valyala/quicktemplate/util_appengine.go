// +build appengine appenginevm

package quicktemplate

func unsafeStrToBytes(s string) []byte {
	return []byte(s)
}

func unsafeBytesToStr(z []byte) string {
	return string(z)
}
