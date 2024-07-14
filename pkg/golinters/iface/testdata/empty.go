//golangcitest:args -Eiface
//golangcitest:config_path testdata/empty.yml
package testdata

type Entity interface { // want "interface Entity is empty, providing no meaningful behavior"
}
