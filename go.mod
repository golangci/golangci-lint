module github.com/golangci/golangci-lint

go 1.20

require (
	4d63.com/gocheckcompilerdirectives v1.2.1
	4d63.com/gochecknoglobals v0.2.1
	github.com/4meepo/tagalign v1.3.3
	github.com/Abirdcfly/dupword v0.0.13
	github.com/Antonboom/errname v0.1.12
	github.com/Antonboom/nilnil v0.1.7
	github.com/Antonboom/testifylint v1.0.2
	github.com/BurntSushi/toml v1.3.2
	github.com/Djarvur/go-err113 v0.0.0-20210108212216-aea10b59be24
	github.com/GaijinEntertainment/go-exhaustruct/v3 v3.1.0
	github.com/OpenPeeDeeP/depguard/v2 v2.2.0
	github.com/alecthomas/go-check-sumtype v0.1.3
	github.com/alexkohler/nakedret/v2 v2.0.2
	github.com/alexkohler/prealloc v1.0.0
	github.com/alingse/asasalint v0.0.11
	github.com/ashanbrown/forbidigo v1.6.0
	github.com/ashanbrown/makezero v1.1.1
	github.com/bkielbasa/cyclop v1.2.1
	github.com/blizzy78/varnamelen v0.8.0
	github.com/bombsimon/wsl/v4 v4.2.0
	github.com/breml/bidichk v0.2.7
	github.com/breml/errchkjson v0.3.6
	github.com/butuzov/ireturn v0.2.2
	github.com/butuzov/mirror v1.1.0
	github.com/catenacyber/perfsprint v0.5.0
	github.com/charithe/durationcheck v0.0.10
	github.com/curioswitch/go-reassign v0.2.0
	github.com/daixiang0/gci v0.12.1
	github.com/denis-tingaikin/go-header v0.4.3
	github.com/esimonov/ifshort v1.0.4
	github.com/fatih/color v1.16.0
	github.com/firefart/nonamedreturns v1.0.4
	github.com/fzipp/gocyclo v0.6.0
	github.com/ghostiam/protogetter v0.3.3
	github.com/go-critic/go-critic v0.11.0
	github.com/go-xmlfmt/xmlfmt v1.1.2
	github.com/gofrs/flock v0.8.1
	github.com/golangci/check v0.0.0-20180506172741-cfe4005ccda2
	github.com/golangci/dupl v0.0.0-20180902072040-3e9179ac440a
	github.com/golangci/go-misc v0.0.0-20220329215616-d24fe342adfe
	github.com/golangci/gofmt v0.0.0-20231018234816-f50ced29576e
	github.com/golangci/lint-1 v0.0.0-20191013205115-297bf364a8e0
	github.com/golangci/maligned v0.0.0-20180506175553-b1d89398deca
	github.com/golangci/misspell v0.4.1
	github.com/golangci/revgrep v0.5.2
	github.com/golangci/unconvert v0.0.0-20180507085042-28b1c447d1f4
	github.com/gordonklaus/ineffassign v0.1.0
	github.com/gostaticanalysis/forcetypeassert v0.1.0
	github.com/gostaticanalysis/nilerr v0.1.1
	github.com/hashicorp/go-version v1.6.0
	github.com/hexops/gotextdiff v1.0.3
	github.com/jgautheron/goconst v1.7.0
	github.com/jingyugao/rowserrcheck v1.1.1
	github.com/jirfag/go-printf-func-name v0.0.0-20200119135958-7558a9eaa5af
	github.com/jjti/go-spancheck v0.5.2
	github.com/julz/importas v0.1.0
	github.com/kisielk/errcheck v1.6.3
	github.com/kkHAIKE/contextcheck v1.1.4
	github.com/kulti/thelper v0.6.3
	github.com/kunwardeep/paralleltest v1.0.9
	github.com/kyoh86/exportloopref v0.1.11
	github.com/ldez/gomoddirectives v0.2.3
	github.com/ldez/tagliatelle v0.5.0
	github.com/leonklingele/grouper v1.1.1
	github.com/lufeee/execinquery v1.2.1
	github.com/macabu/inamedparam v0.1.3
	github.com/maratori/testableexamples v1.0.0
	github.com/maratori/testpackage v1.1.1
	github.com/matoous/godox v0.0.0-20230222163458-006bad1f9d26
	github.com/mattn/go-colorable v0.1.13
	github.com/mbilski/exhaustivestruct v1.2.0
	github.com/mgechev/revive v1.3.5
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-ps v1.0.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/moricho/tparallel v0.3.1
	github.com/nakabonne/nestif v0.3.1
	github.com/nishanths/exhaustive v0.12.0
	github.com/nishanths/predeclared v0.2.2
	github.com/nunnatsa/ginkgolinter v0.15.2
	github.com/polyfloyd/go-errorlint v1.4.7
	github.com/quasilyte/go-ruleguard/dsl v0.3.22
	github.com/ryancurrah/gomodguard v1.3.0
	github.com/ryanrolds/sqlclosecheck v0.5.1
	github.com/sanposhiho/wastedassign/v2 v2.0.7
	github.com/sashamelentyev/interfacebloat v1.1.0
	github.com/sashamelentyev/usestdlibvars v1.24.0
	github.com/securego/gosec/v2 v2.18.2
	github.com/shazow/go-diff v0.0.0-20160112020656-b6b7b6733b8c
	github.com/shirou/gopsutil/v3 v3.23.12
	github.com/sirupsen/logrus v1.9.3
	github.com/sivchari/containedctx v1.0.3
	github.com/sivchari/nosnakecase v1.7.0
	github.com/sivchari/tenv v1.7.1
	github.com/sonatard/noctx v0.0.2
	github.com/sourcegraph/go-diff v0.7.0
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.12.0
	github.com/ssgreg/nlreturn/v2 v2.2.1
	github.com/stbenjam/no-sprintf-host-port v0.1.1
	github.com/stretchr/testify v1.8.4
	github.com/tdakkota/asciicheck v0.2.0
	github.com/tetafro/godot v1.4.16
	github.com/timakin/bodyclose v0.0.0-20230421092635-574207250966
	github.com/timonwong/loggercheck v0.9.4
	github.com/tomarrell/wrapcheck/v2 v2.8.1
	github.com/tommy-muehle/go-mnd/v2 v2.5.1
	github.com/ultraware/funlen v0.1.0
	github.com/ultraware/whitespace v0.1.0
	github.com/uudashr/gocognit v1.1.2
	github.com/valyala/quicktemplate v1.7.0
	github.com/xen0n/gosmopolitan v1.2.2
	github.com/yagipy/maintidx v1.0.0
	github.com/yeya24/promlinter v0.2.0
	github.com/ykadowak/zerologlint v0.1.5
	gitlab.com/bosi/decorder v0.4.1
	go-simpler.org/musttag v0.8.0
	go-simpler.org/sloglint v0.4.0
	golang.org/x/exp v0.0.0-20240103183307-be819d1f06fc
	golang.org/x/tools v0.17.0
	gopkg.in/yaml.v3 v3.0.1
	honnef.co/go/tools v0.4.6
	mvdan.cc/gofumpt v0.5.0
	mvdan.cc/interfacer v0.0.0-20180901003855-c20040233aed
	mvdan.cc/unparam v0.0.0-20221223090309-7455f1af531d
)

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/ccojocar/zxcvbn-go v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/chavacava/garif v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/ettle/strcase v0.2.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-toolsmith/astcast v1.1.0 // indirect
	github.com/go-toolsmith/astcopy v1.1.0 // indirect
	github.com/go-toolsmith/astequal v1.1.0 // indirect
	github.com/go-toolsmith/astfmt v1.1.0 // indirect
	github.com/go-toolsmith/astp v1.1.0 // indirect
	github.com/go-toolsmith/strparse v1.1.0 // indirect
	github.com/go-toolsmith/typep v1.1.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/gostaticanalysis/analysisutil v0.7.1 // indirect
	github.com/gostaticanalysis/comment v1.4.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kisielk/gotool v1.0.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/quasilyte/go-ruleguard v0.4.0 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/regex/syntax v0.0.0-20210819130434-b3f0c404a727 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/t-yuki/gocover-cobertura v0.0.0-20180217150009-aaee18c8195c // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20231219180239-dc181d75b848 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	mvdan.cc/lint v0.0.0-20170908181259-adc824a0674b // indirect
)
