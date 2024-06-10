//golangcitest:args -Ettempdir
//golangcitest:config_path testdata/ttempdir_max_recursion_level.yml
package testdata

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRecursive(t *testing.T) {
	t.Log( // recursion level 1
		os.TempDir(), // want "os\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestRecursive"
	)
	t.Log( // recursion level 1
		fmt.Sprintf("%s", // recursion level 2
			os.TempDir(), // want "os\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestRecursive"
		),
	)
	t.Log( // recursion level 1
		filepath.Clean( // recursion level 2
			fmt.Sprintf("%s", // recursion level 3
				os.TempDir(), // want "os\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestRecursive"
			),
		),
	)
	t.Log( // recursion level 1
		filepath.Join( // recursion level 2
			filepath.Clean( // recursion level 3
				fmt.Sprintf("%s", // recursion level 4
					os.TempDir(), // want "os\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestRecursive"
				),
			),
			"test",
		),
	)
	t.Log( // recursion level 1
		fmt.Sprintf("%s/foo-%d", // recursion level 2
			filepath.Join( // recursion level 3
				filepath.Clean( // recursion level 4
					fmt.Sprintf("%s", // recursion level 5
						os.TempDir(), // want "os\\.TempDir\\(\\) should be replaced by `t\\.TempDir\\(\\)` in TestRecursive"
					),
				),
				"test",
			),
			1024,
		),
	)
}
