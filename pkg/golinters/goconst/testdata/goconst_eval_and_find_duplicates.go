//golangcitest:args -Egoconst
//golangcitest:config_path testdata/goconst_eval_and_find_duplicates.yml
package testdata

import "fmt"

const (
	envPrefix   = "FOO_"
	EnvUser     = envPrefix + "USER"
	EnvPassword = envPrefix + "PASSWORD"
)

const EnvUserFull = "FOO_USER" // want "This constant is a duplicate of `EnvUser` at .*goconst_eval_and_find_duplicates.go:9:16"

const KiB = 1 << 10

func _() {
	fmt.Println(envPrefix, EnvUser, EnvPassword, EnvUserFull)

	const kilobytes = 1024 // want "This constant is a duplicate of `KiB` at .*goconst_eval_and_find_duplicates.go:15:13"
	fmt.Println(kilobytes)

	kib := 1024
	fmt.Println(kib)
}
