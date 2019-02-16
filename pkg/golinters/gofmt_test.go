package golinters

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/golangci/golangci-lint/pkg/result"

	"github.com/stretchr/testify/assert"
	diffpkg "sourcegraph.com/sourcegraph/go-diff/diff"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

func testDiffProducesChanges(t *testing.T, log logutils.Log, diff string, expectedChanges ...Change) {
	diffs, err := diffpkg.ParseMultiFileDiff([]byte(diff))
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Len(t, diffs, 1)
	hunks := diffs[0].Hunks
	assert.NotEmpty(t, hunks)

	var changes []Change
	for _, hunk := range hunks {
		p := hunkChangesParser{
			log: log,
		}
		changes = append(changes, p.parse(hunk)...)
	}

	assert.Equal(t, expectedChanges, changes)
}

func TestExtractChangesFromHunkAddOnly(t *testing.T) {
	const diff = `diff --git a/internal/shared/logutil/log.go b/internal/shared/logutil/log.go
index 258b340..43d04bf 100644
--- a/internal/shared/logutil/log.go
+++ b/internal/shared/logutil/log.go
@@ -1,5 +1,6 @@
 package logutil

+// added line
 type Func func(format string, args ...interface{})

 type Log interface {
`

	testDiffProducesChanges(t, nil, diff, Change{
		LineRange: result.Range{
			From: 2,
			To:   2,
		},
		Replacement: result.Replacement{
			NewLines: []string{
				"",
				"// added line",
			},
		},
	})
}

func TestExtractChangesFromHunkAddOnlyOnFirstLine(t *testing.T) {
	const diff = `diff --git a/internal/shared/logutil/log.go b/internal/shared/logutil/log.go
index 258b340..97e6660 100644
--- a/internal/shared/logutil/log.go
+++ b/internal/shared/logutil/log.go
@@ -1,3 +1,4 @@
+// added line
 package logutil

 type Func func(format string, args ...interface{})
`

	testDiffProducesChanges(t, nil, diff, Change{
		LineRange: result.Range{
			From: 1,
			To:   1,
		},
		Replacement: result.Replacement{
			NewLines: []string{
				"// added line",
				"package logutil",
			},
		},
	})
}

func TestExtractChangesFromHunkAddOnlyOnFirstLineWithSharedOriginalLine(t *testing.T) {
	const diff = `diff --git a/internal/shared/logutil/log.go b/internal/shared/logutil/log.go
index 258b340..7ff80c9 100644
--- a/internal/shared/logutil/log.go
+++ b/internal/shared/logutil/log.go
@@ -1,4 +1,7 @@
+// added line 1
 package logutil
+// added line 2
+// added line 3

 type Func func(format string, args ...interface{})
`
	testDiffProducesChanges(t, nil, diff, Change{
		LineRange: result.Range{
			From: 1,
			To:   1,
		},
		Replacement: result.Replacement{
			NewLines: []string{
				"// added line 1",
				"package logutil",
				"// added line 2",
				"// added line 3",
			},
		},
	})
}

func TestExtractChangesFromHunkAddOnlyInAllDiff(t *testing.T) {
	const diff = `diff --git a/test.go b/test.go
new file mode 100644
index 0000000..6399915
--- /dev/null
+++ b/test.go
@@ -0,0 +1,3 @@
+package test
+
+// line
`

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log := logutils.NewMockLog(ctrl)
	log.EXPECT().Infof("The diff contains only additions: no original or deleted lines: %#v", gomock.Any())
	var noChanges []Change
	testDiffProducesChanges(t, log, diff, noChanges...)
}

func TestExtractChangesFromHunkAddOnlyMultipleLines(t *testing.T) {
	const diff = `diff --git a/internal/shared/logutil/log.go b/internal/shared/logutil/log.go
index 258b340..3b83a94 100644
--- a/internal/shared/logutil/log.go
+++ b/internal/shared/logutil/log.go
@@ -2,6 +2,9 @@ package logutil

 type Func func(format string, args ...interface{})

+// add line 1
+// add line 2
+
 type Log interface {
        Fatalf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
`

	testDiffProducesChanges(t, nil, diff, Change{
		LineRange: result.Range{
			From: 4,
			To:   4,
		},
		Replacement: result.Replacement{
			NewLines: []string{
				"",
				"// add line 1",
				"// add line 2",
				"",
			},
		},
	})
}

func TestExtractChangesFromHunkAddOnlyDifferentLines(t *testing.T) {
	const diff = `diff --git a/internal/shared/logutil/log.go b/internal/shared/logutil/log.go
index 258b340..e5ed2ad 100644
--- a/internal/shared/logutil/log.go
+++ b/internal/shared/logutil/log.go
@@ -2,9 +2,12 @@ package logutil

 type Func func(format string, args ...interface{})

+// add line 1
+
 type Log interface {
        Fatalf(format string, args ...interface{})
        Errorf(format string, args ...interface{})
+       // add line 2
        Warnf(format string, args ...interface{})
        Infof(format string, args ...interface{})
		Debugf(key string, format string, args ...interface{})
`

	expectedChanges := []Change{
		{
			LineRange: result.Range{
				From: 4,
				To:   4,
			},
			Replacement: result.Replacement{
				NewLines: []string{
					"",
					"// add line 1",
					"",
				},
			},
		},
		{
			LineRange: result.Range{
				From: 7,
				To:   7,
			},
			Replacement: result.Replacement{
				NewLines: []string{
					"       Errorf(format string, args ...interface{})",
					"       // add line 2",
				},
			},
		},
	}

	testDiffProducesChanges(t, nil, diff, expectedChanges...)
}

func TestExtractChangesDeleteOnlyFirstLines(t *testing.T) {
	const diff = `diff --git a/internal/shared/logutil/log.go b/internal/shared/logutil/log.go
index 258b340..0fb554e 100644
--- a/internal/shared/logutil/log.go
+++ b/internal/shared/logutil/log.go
@@ -1,5 +1,3 @@
-package logutil
-
 type Func func(format string, args ...interface{})

 type Log interface {
`

	testDiffProducesChanges(t, nil, diff, Change{
		LineRange: result.Range{
			From: 1,
			To:   2,
		},
		Replacement: result.Replacement{
			NeedOnlyDelete: true,
		},
	})
}

func TestExtractChangesReplaceLine(t *testing.T) {
	const diff = `diff --git a/internal/shared/logutil/log.go b/internal/shared/logutil/log.go
index 258b340..c2a8516 100644
--- a/internal/shared/logutil/log.go
+++ b/internal/shared/logutil/log.go
@@ -1,4 +1,4 @@
-package logutil
+package test2

 type Func func(format string, args ...interface{})
`

	testDiffProducesChanges(t, nil, diff, Change{
		LineRange: result.Range{
			From: 1,
			To:   1,
		},
		Replacement: result.Replacement{
			NewLines: []string{"package test2"},
		},
	})
}

func TestExtractChangesReplaceLineAfterFirstLineAdding(t *testing.T) {
	const diff = `diff --git a/internal/shared/logutil/log.go b/internal/shared/logutil/log.go
index 258b340..43fc0de 100644
--- a/internal/shared/logutil/log.go
+++ b/internal/shared/logutil/log.go
@@ -1,6 +1,7 @@
+// added line
 package logutil

-type Func func(format string, args ...interface{})
+// changed line

 type Log interface {
        Fatalf(format string, args ...interface{})`

	testDiffProducesChanges(t, nil, diff, Change{
		LineRange: result.Range{
			From: 1,
			To:   1,
		},
		Replacement: result.Replacement{
			NewLines: []string{
				"// added line",
				"package logutil",
			},
		},
	}, Change{
		LineRange: result.Range{
			From: 3,
			To:   3,
		},
		Replacement: result.Replacement{
			NewLines: []string{
				"// changed line",
			},
		},
	})
}

func TestGofmtDiff(t *testing.T) {
	const diff = `diff --git a/gofmt.go b/gofmt.go
index 2c9f78d..c0d5791 100644
--- a/gofmt.go
+++ b/gofmt.go
@@ -1,9 +1,9 @@
 //args: -Egofmt
 package p

- func gofmt(a, b int) int {
-         if a != b {
-                 return 1
+func gofmt(a, b int) int {
+       if a != b {
+               return 1
        }
-         return 2
+       return 2
 }
`
	testDiffProducesChanges(t, nil, diff, Change{
		LineRange: result.Range{
			From: 4,
			To:   6,
		},
		Replacement: result.Replacement{
			NewLines: []string{
				"func gofmt(a, b int) int {",
				"       if a != b {",
				"               return 1",
			},
		},
	}, Change{
		LineRange: result.Range{
			From: 8,
			To:   8,
		},
		Replacement: result.Replacement{
			NewLines: []string{
				"       return 2",
			},
		},
	})
}
