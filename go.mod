module github.com/golangci/golangci-lint

go 1.11

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/OpenPeeDeeP/depguard v1.0.0
	github.com/StackExchange/wmi v0.0.0-20180116203802-5d049714c4a6 // indirect
	github.com/fatih/color v1.6.0
	github.com/go-critic/go-critic v0.3.5-0.20190526074819-1df300866540
	github.com/go-lintpack/lintpack v0.5.2
	github.com/go-ole/go-ole v1.2.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/mock v1.0.0
	github.com/golangci/check v0.0.0-20180506172741-cfe4005ccda2
	github.com/golangci/dupl v0.0.0-20180902072040-3e9179ac440a
	github.com/golangci/errcheck v0.0.0-20181003203344-ef45e06d44b6
	github.com/golangci/go-misc v0.0.0-20180628070357-927a3d87b613
	github.com/golangci/go-tools v0.0.0-20190318055746-e32c54105b7c
	github.com/golangci/goconst v0.0.0-20180610141641-041c5f2b40f3
	github.com/golangci/gocyclo v0.0.0-20180528134321-2becd97e67ee
	github.com/golangci/gofmt v0.0.0-20181222123516-0b8337e80d98
	github.com/golangci/gosec v0.0.0-20180901114220-66fb7fc33547
	github.com/golangci/ineffassign v0.0.0-20190609212857-42439a7714cc
	github.com/golangci/lint-1 v0.0.0-20180610141402-ee948d087217
	github.com/golangci/maligned v0.0.0-20180506175553-b1d89398deca
	github.com/golangci/misspell v0.0.0-20180809174111-950f5d19e770
	github.com/golangci/prealloc v0.0.0-20180630174525-215b22d4de21
	github.com/golangci/revgrep v0.0.0-20180526074752-d9c87f5ffaf0
	github.com/golangci/unconvert v0.0.0-20180507085042-28b1c447d1f4
	github.com/hashicorp/hcl v0.0.0-20180404174102-ef8a98b0bbce // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/magiconair/properties v1.7.6 // indirect
	github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-isatty v0.0.3 // indirect
	github.com/mitchellh/go-homedir v1.0.0
	github.com/mitchellh/go-ps v0.0.0-20170309133038-4fdf99ab2936
	github.com/mitchellh/mapstructure v0.0.0-20180220230111-00c29f56e238 // indirect
	github.com/nbutton23/zxcvbn-go v0.0.0-20171102151520-eafdab6b0663 // indirect
	github.com/onsi/gomega v1.4.2 // indirect
	github.com/pelletier/go-toml v1.1.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/shirou/gopsutil v0.0.0-20180427012116-c95755e4bcd7
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/sirupsen/logrus v1.0.5
	github.com/sourcegraph/go-diff v0.5.1
	github.com/spf13/afero v1.1.0 // indirect
	github.com/spf13/cast v1.2.0 // indirect
	github.com/spf13/cobra v0.0.2
	github.com/spf13/jwalterweatherman v0.0.0-20180109140146-7c0cea34c8ec // indirect
	github.com/spf13/pflag v1.0.1
	github.com/spf13/viper v1.0.2
	github.com/stretchr/testify v1.2.2
	github.com/timakin/bodyclose v0.0.0-00010101000000-87058b9bfcec
	github.com/ultraware/funlen v0.0.1
	github.com/valyala/quicktemplate v1.1.1
	golang.org/x/crypto v0.0.0-20190313024323-a1f597ede03a // indirect
	golang.org/x/sys v0.0.0-20190312061237-fead79001313 // indirect
	golang.org/x/tools v0.0.0-20190521203540-521d6ed310dd
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/yaml.v2 v2.2.1
	mvdan.cc/interfacer v0.0.0-20180901003855-c20040233aed
	mvdan.cc/lint v0.0.0-20170908181259-adc824a0674b // indirect
	mvdan.cc/unparam v0.0.0-20190124213536-fbb59629db34
)

// https://github.com/golang/tools/pull/139
replace golang.org/x/tools => github.com/golangci/tools v0.0.0-20190713050349-979bdb7f8cc8
