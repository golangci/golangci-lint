module github.com/golangci/golangci-lint

go 1.12

require (
	github.com/OpenPeeDeeP/depguard v1.0.1
	github.com/bombsimon/wsl/v2 v2.0.0
	github.com/fatih/color v1.7.0
	github.com/go-critic/go-critic v0.4.1
	github.com/go-lintpack/lintpack v0.5.2
	github.com/gofrs/flock v0.0.0-20190320160742-5135e617513b
	github.com/golangci/check v0.0.0-20180506172741-cfe4005ccda2
	github.com/golangci/dupl v0.0.0-20180902072040-3e9179ac440a
	github.com/golangci/errcheck v0.0.0-20181223084120-ef45e06d44b6
	github.com/golangci/go-misc v0.0.0-20180628070357-927a3d87b613
	github.com/golangci/goconst v0.0.0-20180610141641-041c5f2b40f3
	github.com/golangci/gocyclo v0.0.0-20180528134321-2becd97e67ee
	github.com/golangci/gofmt v0.0.0-20190930125516-244bba706f1a
	github.com/golangci/ineffassign v0.0.0-20190609212857-42439a7714cc
	github.com/golangci/lint-1 v0.0.0-20191013205115-297bf364a8e0
	github.com/golangci/maligned v0.0.0-20180506175553-b1d89398deca
	github.com/golangci/misspell v0.0.0-20180809174111-950f5d19e770
	github.com/golangci/prealloc v0.0.0-20180630174525-215b22d4de21
	github.com/golangci/revgrep v0.0.0-20180526074752-d9c87f5ffaf0
	github.com/golangci/unconvert v0.0.0-20180507085042-28b1c447d1f4
	github.com/jirfag/go-printf-func-name v0.0.0-20191110105641-45db9963cdd3
	github.com/jingyugao/rowserrcheck v0.0.0-20191204022205-72ab7603b68a
	github.com/matoous/godox v0.0.0-20190911065817-5d6d842e92eb // v1.0
	github.com/mattn/go-colorable v0.1.4
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-ps v0.0.0-20190716172923-621e5597135b
	github.com/pkg/errors v0.8.1
	github.com/securego/gosec v0.0.0-20200103095621-79fbf3af8d83
	github.com/shirou/gopsutil v0.0.0-20190901111213-e4ec7b275ada // v2.19.8
	github.com/sirupsen/logrus v1.4.2
	github.com/sourcegraph/go-diff v0.5.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0
	github.com/timakin/bodyclose v0.0.0-20190930140734-f7f2e9bca95e
	github.com/tommy-muehle/go-mnd v1.1.1
	github.com/ultraware/funlen v0.0.2
	github.com/ultraware/whitespace v0.0.4
	github.com/uudashr/gocognit v1.0.1
	github.com/valyala/quicktemplate v1.2.0
	golang.org/x/tools v0.0.0-20200102140908-9497f49d5709
	gopkg.in/yaml.v2 v2.2.7
	honnef.co/go/tools v0.0.1-2019.2.3
	mvdan.cc/interfacer v0.0.0-20180901003855-c20040233aed
	mvdan.cc/lint v0.0.0-20170908181259-adc824a0674b // indirect
	mvdan.cc/unparam v0.0.0-20190720180237-d51796306d8f
)

// https://github.com/golang/tools/pull/156
// https://github.com/golang/tools/pull/160
replace golang.org/x/tools => github.com/golangci/tools v0.0.0-20190915081525-6aa350649b1c
