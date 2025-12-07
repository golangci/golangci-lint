module github.com/golangci/golangci-lint/v2

go 1.24.0

require (
	4d63.com/gocheckcompilerdirectives v1.3.0
	4d63.com/gochecknoglobals v0.2.2
	dev.gaijin.team/go/exhaustruct/v4 v4.0.0
	github.com/4meepo/tagalign v1.4.3
	github.com/Abirdcfly/dupword v0.1.7
	github.com/AdminBenni/iota-mixing v1.0.0
	github.com/AlwxSin/noinlineerr v1.0.5
	github.com/Antonboom/errname v1.1.1
	github.com/Antonboom/nilnil v1.1.1
	github.com/Antonboom/testifylint v1.6.4
	github.com/BurntSushi/toml v1.5.0
	github.com/Djarvur/go-err113 v0.1.1
	github.com/MirrexOne/unqueryvet v1.3.0
	github.com/OpenPeeDeeP/depguard/v2 v2.2.1
	github.com/alecthomas/chroma/v2 v2.20.0
	github.com/alecthomas/go-check-sumtype v0.3.1
	github.com/alexkohler/nakedret/v2 v2.0.6
	github.com/alexkohler/prealloc v1.0.0
	github.com/alingse/asasalint v0.0.11
	github.com/alingse/nilnesserr v0.2.0
	github.com/ashanbrown/forbidigo/v2 v2.3.0
	github.com/ashanbrown/makezero/v2 v2.1.0
	github.com/bkielbasa/cyclop v1.2.3
	github.com/blizzy78/varnamelen v0.8.0
	github.com/bombsimon/wsl/v4 v4.7.0
	github.com/bombsimon/wsl/v5 v5.3.0
	github.com/breml/bidichk v0.3.3
	github.com/breml/errchkjson v0.4.1
	github.com/butuzov/ireturn v0.4.0
	github.com/butuzov/mirror v1.3.0
	github.com/catenacyber/perfsprint v0.10.1
	github.com/charithe/durationcheck v0.0.11
	github.com/charmbracelet/lipgloss v1.1.0
	github.com/ckaznocha/intrange v0.3.1
	github.com/curioswitch/go-reassign v0.3.0
	github.com/daixiang0/gci v0.13.7
	github.com/denis-tingaikin/go-header v0.5.0
	github.com/fatih/color v1.18.0
	github.com/firefart/nonamedreturns v1.0.6
	github.com/fzipp/gocyclo v0.6.0
	github.com/ghostiam/protogetter v0.3.17
	github.com/go-critic/go-critic v0.14.2
	github.com/go-viper/mapstructure/v2 v2.4.0
	github.com/go-xmlfmt/xmlfmt v1.1.3
	github.com/godoc-lint/godoc-lint v0.10.2
	github.com/gofrs/flock v0.13.0
	github.com/golangci/asciicheck v0.5.0
	github.com/golangci/dupl v0.0.0-20250308024227-f665c8d69b32
	github.com/golangci/go-printf-func-name v0.1.1
	github.com/golangci/gofmt v0.0.0-20250106114630-d62b90e6713d
	github.com/golangci/golines v0.0.0-20250217134842-442fd0091d95
	github.com/golangci/misspell v0.7.0
	github.com/golangci/plugin-module-register v0.1.2
	github.com/golangci/revgrep v0.8.0
	github.com/golangci/swaggoswag v0.0.0-20250504205917-77f2aca3143e
	github.com/golangci/unconvert v0.0.0-20250410112200-a129a6e6413e
	github.com/gordonklaus/ineffassign v0.2.0
	github.com/gostaticanalysis/forcetypeassert v0.2.0
	github.com/gostaticanalysis/nilerr v0.1.2
	github.com/hashicorp/go-version v1.8.0
	github.com/jgautheron/goconst v1.8.2
	github.com/jingyugao/rowserrcheck v1.1.1
	github.com/jjti/go-spancheck v0.6.5
	github.com/julz/importas v0.2.0
	github.com/karamaru-alpha/copyloopvar v1.2.2
	github.com/kisielk/errcheck v1.9.0
	github.com/kkHAIKE/contextcheck v1.1.6
	github.com/kulti/thelper v0.7.1
	github.com/kunwardeep/paralleltest v1.0.15
	github.com/lasiar/canonicalheader v1.1.2
	github.com/ldez/exptostd v0.4.5
	github.com/ldez/gomoddirectives v0.7.1
	github.com/ldez/grignotin v0.10.1
	github.com/ldez/tagliatelle v0.7.2
	github.com/ldez/usetesting v0.5.0
	github.com/leonklingele/grouper v1.1.2
	github.com/macabu/inamedparam v0.2.0
	github.com/manuelarte/embeddedstructfieldcheck v0.4.0
	github.com/manuelarte/funcorder v0.5.0
	github.com/maratori/testableexamples v1.0.1
	github.com/maratori/testpackage v1.1.2
	github.com/matoous/godox v1.1.0
	github.com/mattn/go-colorable v0.1.14
	github.com/mgechev/revive v1.13.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/moricho/tparallel v0.3.2
	github.com/nakabonne/nestif v0.3.1
	github.com/nishanths/exhaustive v0.12.0
	github.com/nishanths/predeclared v0.2.2
	github.com/nunnatsa/ginkgolinter v0.21.2
	github.com/pelletier/go-toml/v2 v2.2.4
	github.com/polyfloyd/go-errorlint v1.8.0
	github.com/quasilyte/go-ruleguard/dsl v0.3.23
	github.com/raeperd/recvcheck v0.2.0
	github.com/rogpeppe/go-internal v1.14.1
	github.com/ryancurrah/gomodguard v1.4.1
	github.com/ryanrolds/sqlclosecheck v0.5.1
	github.com/sanposhiho/wastedassign/v2 v2.1.0
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2
	github.com/sashamelentyev/interfacebloat v1.1.0
	github.com/sashamelentyev/usestdlibvars v1.29.0
	github.com/securego/gosec/v2 v2.22.11-0.20251204091113-daccba6b93d7
	github.com/shirou/gopsutil/v4 v4.25.11
	github.com/sirupsen/logrus v1.9.3
	github.com/sivchari/containedctx v1.0.3
	github.com/sonatard/noctx v0.4.0
	github.com/sourcegraph/go-diff v0.7.0
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10
	github.com/spf13/viper v1.12.0
	github.com/ssgreg/nlreturn/v2 v2.2.1
	github.com/stbenjam/no-sprintf-host-port v0.3.1
	github.com/stretchr/testify v1.11.1
	github.com/tetafro/godot v1.5.4
	github.com/timakin/bodyclose v0.0.0-20241222091800-1db5c5ca4d67
	github.com/timonwong/loggercheck v0.11.0
	github.com/tomarrell/wrapcheck/v2 v2.12.0
	github.com/tommy-muehle/go-mnd/v2 v2.5.1
	github.com/ultraware/funlen v0.2.0
	github.com/ultraware/whitespace v0.2.0
	github.com/uudashr/gocognit v1.2.0
	github.com/uudashr/iface v1.4.1
	github.com/valyala/quicktemplate v1.8.0
	github.com/xen0n/gosmopolitan v1.3.0
	github.com/yagipy/maintidx v1.0.0
	github.com/yeya24/promlinter v0.3.0
	github.com/ykadowak/zerologlint v0.1.5
	gitlab.com/bosi/decorder v0.4.2
	go-simpler.org/musttag v0.14.0
	go-simpler.org/sloglint v0.11.1
	go.augendre.info/arangolint v0.3.1
	go.augendre.info/fatcontext v0.9.0
	go.uber.org/automaxprocs v1.6.0
	go.yaml.in/yaml/v3 v3.0.4
	golang.org/x/mod v0.30.0
	golang.org/x/sync v0.18.0
	golang.org/x/sys v0.38.0
	golang.org/x/tools v0.39.0
	honnef.co/go/tools v0.6.1
	mvdan.cc/gofumpt v0.9.2
	mvdan.cc/unparam v0.0.0-20251027182757-5beb8c8f8f15
)

require (
	codeberg.org/chavacava/garif v0.2.0 // indirect
	dev.gaijin.team/go/golib v0.6.0 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/alfatraining/structtag v1.0.0 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/ccojocar/zxcvbn-go v1.0.4 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/charmbracelet/colorprofile v0.2.3-0.20250311203215-f60798e515dc // indirect
	github.com/charmbracelet/x/ansi v0.8.0 // indirect
	github.com/charmbracelet/x/cellbuf v0.0.13-0.20250311204145-2c3ea96c31dd // indirect
	github.com/charmbracelet/x/term v0.2.1 // indirect
	github.com/dave/dst v0.27.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/ebitengine/purego v0.9.1 // indirect
	github.com/ettle/strcase v0.2.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-toolsmith/astcast v1.1.0 // indirect
	github.com/go-toolsmith/astcopy v1.1.0 // indirect
	github.com/go-toolsmith/astequal v1.2.0 // indirect
	github.com/go-toolsmith/astfmt v1.1.0 // indirect
	github.com/go-toolsmith/astp v1.1.0 // indirect
	github.com/go-toolsmith/strparse v1.1.0 // indirect
	github.com/go-toolsmith/typep v1.1.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/gostaticanalysis/analysisutil v0.7.1 // indirect
	github.com/gostaticanalysis/comment v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix/v2 v2.1.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hexops/gotextdiff v1.0.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/muesli/termenv v0.16.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/quasilyte/go-ruleguard v0.4.5 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/regex/syntax v0.0.0-20210819130434-b3f0c404a727 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/tklauser/go-sysconf v0.3.16 // indirect
	github.com/tklauser/numcpus v0.11.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20251023183803-a4bb9ffd2546 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/tools/go/expect v0.1.1-deprecated // indirect
	golang.org/x/tools/go/packages/packagestest v0.1.1-deprecated // indirect
	google.golang.org/protobuf v1.36.8 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
