package executors

import "context"

type Executor interface {
	Run(ctx context.Context, name string, args ...string) (string, error)

	WithEnv(k, v string) Executor
	SetEnv(k, v string)

	WorkDir() string
	WithWorkDir(wd string) Executor

	Clean()
}
