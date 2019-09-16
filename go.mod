module github.com/golangci/golangci-lint

go 1.11

require (
	github.com/OpenPeeDeeP/depguard v1.0.1
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/fatih/color v1.7.0
	github.com/go-critic/go-critic v0.3.5-0.20190904082202-d79a9f0c64db
	github.com/go-lintpack/lintpack v0.5.2
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/golang/mock v1.3.1
	github.com/golangci/check v0.0.0-20180506172741-cfe4005ccda2
	github.com/golangci/dupl v0.0.0-20180902072040-3e9179ac440a
	github.com/golangci/errcheck v0.0.0-20181223084120-ef45e06d44b6
	github.com/golangci/go-misc v0.0.0-20180628070357-927a3d87b613
	github.com/golangci/go-tools v0.0.0-20190318055746-e32c54105b7c
	github.com/golangci/goconst v0.0.0-20180610141641-041c5f2b40f3
	github.com/golangci/gocyclo v0.0.0-20180528134321-2becd97e67ee
	github.com/golangci/gofmt v0.0.0-20181222123516-0b8337e80d98
	github.com/golangci/ineffassign v0.0.0-20190609212857-42439a7714cc
	github.com/golangci/lint-1 v0.0.0-20190420132249-ee948d087217
	github.com/golangci/maligned v0.0.0-20180506175553-b1d89398deca
	github.com/golangci/misspell v0.0.0-20180809174111-950f5d19e770
	github.com/golangci/prealloc v0.0.0-20180630174525-215b22d4de21
	github.com/golangci/revgrep v0.0.0-20180526074752-d9c87f5ffaf0
	github.com/golangci/unconvert v0.0.0-20180507085042-28b1c447d1f4
	github.com/matoous/godox v0.0.0-20190910121045-032ad8106c86
	github.com/mattn/go-colorable v0.1.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-ps v0.0.0-20190716172923-621e5597135b
	github.com/pkg/errors v0.8.1
	github.com/securego/gosec v0.0.0-20190912120752-140048b2a218
	github.com/shirou/gopsutil v2.18.12+incompatible
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/sourcegraph/go-diff v0.5.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/timakin/bodyclose v0.0.0-20190721030226-87058b9bfcec
	github.com/ultraware/funlen v0.0.2
	github.com/ultraware/whitespace v0.0.2
	github.com/valyala/quicktemplate v1.2.0
	golang.org/x/tools v0.0.0-20190912215617-3720d1ec3678
	gopkg.in/yaml.v2 v2.2.2
	mvdan.cc/interfacer v0.0.0-20180901003855-c20040233aed
	mvdan.cc/lint v0.0.0-20170908181259-adc824a0674b // indirect
	mvdan.cc/unparam v0.0.0-20190720180237-d51796306d8f
)

// https://github.com/golang/tools/pull/160
replace golang.org/x/tools => github.com/golangci/tools v0.0.0-20190914130248-e9260b99c8f1
