package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/caarlos0/ctrlc"
	"github.com/fatih/color"
	"github.com/goreleaser/goreleaser/internal/middleware"
	"github.com/goreleaser/goreleaser/internal/pipeline"
	"github.com/goreleaser/goreleaser/internal/static"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/goreleaser/goreleaser/pkg/defaults"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

type releaseOptions struct {
	Config       string
	ReleaseNotes string
	Snapshot     bool
	SkipPublish  bool
	SkipSign     bool
	SkipValidate bool
	RmDist       bool
	Parallelism  int
	Timeout      time.Duration
}

func main() {
	// enable colored output on travis
	if os.Getenv("CI") != "" {
		color.NoColor = false
	}
	log.SetHandler(cli.Default)

	fmt.Println()
	defer fmt.Println()

	var app = kingpin.New("goreleaser", "Deliver Go binaries as fast and easily as possible")
	var debug = app.Flag("debug", "Enable debug mode").Bool()
	var config = app.Flag("config", "Load configuration from file").Short('c').Short('f').PlaceHolder(".goreleaser.yml").String()
	var initCmd = app.Command("init", "Generates a .goreleaser.yml file").Alias("i")
	var checkCmd = app.Command("check", "Checks if configuration is valid").Alias("c")
	var releaseCmd = app.Command("release", "Releases the current project").Alias("r").Default()
	var releaseNotes = releaseCmd.Flag("release-notes", "Load custom release notes from a markdown file").PlaceHolder("notes.md").String()
	var snapshot = releaseCmd.Flag("snapshot", "Generate an unversioned snapshot release, skipping all validations and without publishing any artifacts").Bool()
	var skipPublish = releaseCmd.Flag("skip-publish", "Skips publishing artifacts").Bool()
	var skipSign = releaseCmd.Flag("skip-sign", "Skips signing the artifacts").Bool()
	var skipValidate = releaseCmd.Flag("skip-validate", "Skips several sanity checks").Bool()
	var rmDist = releaseCmd.Flag("rm-dist", "Remove the dist folder before building").Bool()
	var parallelism = releaseCmd.Flag("parallelism", "Amount tasks to run concurrently").Short('p').Default("4").Int()
	var timeout = releaseCmd.Flag("timeout", "Timeout to the entire release process").Default("30m").Duration()

	app.Version(buildVersion(version, commit, date, builtBy))
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	app.UsageTemplate(static.UsageTemplate)

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	switch cmd {
	case initCmd.FullCommand():
		var filename = *config
		if filename == "" {
			filename = ".goreleaser.yml"
		}
		if err := initProject(filename); err != nil {
			log.WithError(err).Error("failed to init project")
			os.Exit(1)
			return
		}
		log.WithField("file", filename).Info("config created; please edit accordingly to your needs")
	case checkCmd.FullCommand():
		if err := checkConfig(*config); err != nil {
			log.WithError(err).Errorf(color.New(color.Bold).Sprintf("config is invalid"))
			os.Exit(1)
			return
		}
		log.Infof(color.New(color.Bold).Sprintf("config is valid"))
	case releaseCmd.FullCommand():
		start := time.Now()
		log.Infof(color.New(color.Bold).Sprintf("releasing using goreleaser %s...", version))
		var options = releaseOptions{
			Config:       *config,
			ReleaseNotes: *releaseNotes,
			Snapshot:     *snapshot,
			SkipPublish:  *skipPublish,
			SkipValidate: *skipValidate,
			SkipSign:     *skipSign,
			RmDist:       *rmDist,
			Parallelism:  *parallelism,
			Timeout:      *timeout,
		}
		if err := releaseProject(options); err != nil {
			log.WithError(err).Errorf(color.New(color.Bold).Sprintf("release failed after %0.2fs", time.Since(start).Seconds()))
			os.Exit(1)
			return
		}
		log.Infof(color.New(color.Bold).Sprintf("release succeeded after %0.2fs", time.Since(start).Seconds()))
	}
}

func checkConfig(filename string) error {
	cfg, err := loadConfig(filename)
	if err != nil {
		return err
	}
	var ctx = context.New(cfg)
	return ctrlc.Default.Run(ctx, func() error {
		for _, pipe := range defaults.Defaulters {
			if err := middleware.ErrHandler(pipe.Default)(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func releaseProject(options releaseOptions) error {
	cfg, err := loadConfig(options.Config)
	if err != nil {
		return err
	}
	ctx, cancel := context.NewWithTimeout(cfg, options.Timeout)
	defer cancel()
	ctx.Parallelism = options.Parallelism
	log.Debugf("parallelism: %v", ctx.Parallelism)
	ctx.ReleaseNotes = options.ReleaseNotes
	ctx.Snapshot = options.Snapshot
	ctx.SkipPublish = ctx.Snapshot || options.SkipPublish
	ctx.SkipValidate = ctx.Snapshot || options.SkipValidate
	ctx.SkipSign = options.SkipSign
	ctx.RmDist = options.RmDist
	return ctrlc.Default.Run(ctx, func() error {
		for _, pipe := range pipeline.Pipeline {
			if err := middleware.Logging(
				pipe.String(),
				middleware.ErrHandler(pipe.Run),
				middleware.DefaultInitialPadding,
			)(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

// InitProject creates an example goreleaser.yml in the current directory
func initProject(filename string) error {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		if err != nil {
			return err
		}
		return fmt.Errorf("%s already exists", filename)
	}
	log.Infof(color.New(color.Bold).Sprintf("Generating %s file", filename))
	return ioutil.WriteFile(filename, []byte(static.ExampleConfig), 0644)
}

func loadConfig(path string) (config.Project, error) {
	if path != "" {
		return config.Load(path)
	}
	for _, f := range [4]string{
		".goreleaser.yml",
		".goreleaser.yaml",
		"goreleaser.yml",
		"goreleaser.yaml",
	} {
		proj, err := config.Load(f)
		if err != nil && os.IsNotExist(err) {
			continue
		}
		return proj, err
	}
	// the user didn't specify a config file and the known possible file names
	// don't exist, so, return an empty config and a nil err.
	log.Warn("could not find a config file, using defaults...")
	return config.Project{}, nil
}

func buildVersion(version, commit, date, builtBy string) string {
	var result = fmt.Sprintf("version: %s", version)
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}
