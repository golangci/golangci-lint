package functions

import (
	"go/types"
	"sync"

	"github.com/golangci/go-tools/callgraph"
	"github.com/golangci/go-tools/callgraph/static"
	"github.com/golangci/go-tools/ssa"
	"github.com/golangci/go-tools/staticcheck/vrp"
)

var stdlibDescs = map[string]Description{
	"errors.New": Description{Pure: true},

	"fmt.Errorf":  Description{Pure: true},
	"fmt.Sprintf": Description{Pure: true},
	"fmt.Sprint":  Description{Pure: true},

	"sort.Reverse": Description{Pure: true},

	"strings.Map":            Description{Pure: true},
	"strings.Repeat":         Description{Pure: true},
	"strings.Replace":        Description{Pure: true},
	"strings.Title":          Description{Pure: true},
	"strings.ToLower":        Description{Pure: true},
	"strings.ToLowerSpecial": Description{Pure: true},
	"strings.ToTitle":        Description{Pure: true},
	"strings.ToTitleSpecial": Description{Pure: true},
	"strings.ToUpper":        Description{Pure: true},
	"strings.ToUpperSpecial": Description{Pure: true},
	"strings.Trim":           Description{Pure: true},
	"strings.TrimFunc":       Description{Pure: true},
	"strings.TrimLeft":       Description{Pure: true},
	"strings.TrimLeftFunc":   Description{Pure: true},
	"strings.TrimPrefix":     Description{Pure: true},
	"strings.TrimRight":      Description{Pure: true},
	"strings.TrimRightFunc":  Description{Pure: true},
	"strings.TrimSpace":      Description{Pure: true},
	"strings.TrimSuffix":     Description{Pure: true},

	"(*net/http.Request).WithContext": Description{Pure: true},

	"math/rand.Read":         Description{NilError: true},
	"(*math/rand.Rand).Read": Description{NilError: true},
}

type Description struct {
	// The function is known to be pure
	Pure bool
	// The function is known to be a stub
	Stub bool
	// The function is known to never return (panics notwithstanding)
	Infinite bool
	// Variable ranges
	Ranges vrp.Ranges
	Loops  []Loop
	// Function returns an error as its last argument, but it is
	// always nil
	NilError            bool
	ConcreteReturnTypes []*types.Tuple
}

type descriptionEntry struct {
	ready  chan struct{}
	result Description
}

type Descriptions struct {
	CallGraph *callgraph.Graph
	mu        sync.Mutex
	cache     map[*ssa.Function]*descriptionEntry
}

func NewDescriptions(prog *ssa.Program) *Descriptions {
	return &Descriptions{
		CallGraph: static.CallGraph(prog),
		cache:     map[*ssa.Function]*descriptionEntry{},
	}
}

func (d *Descriptions) Get(fn *ssa.Function) Description {
	d.mu.Lock()
	fd := d.cache[fn]
	if fd == nil {
		fd = &descriptionEntry{
			ready: make(chan struct{}),
		}
		d.cache[fn] = fd
		d.mu.Unlock()

		{
			fd.result = stdlibDescs[fn.RelString(nil)]
			fd.result.Pure = fd.result.Pure || d.IsPure(fn)
			fd.result.Stub = fd.result.Stub || d.IsStub(fn)
			fd.result.Infinite = fd.result.Infinite || !terminates(fn)
			fd.result.Ranges = vrp.BuildGraph(fn).Solve()
			fd.result.Loops = findLoops(fn)
			fd.result.NilError = fd.result.NilError || IsNilError(fn)
			fd.result.ConcreteReturnTypes = concreteReturnTypes(fn)
		}

		close(fd.ready)
	} else {
		d.mu.Unlock()
		<-fd.ready
	}
	return fd.result
}

func IsNilError(fn *ssa.Function) bool {
	// TODO(dh): This is very simplistic, as we only look for constant
	// nil returns. A more advanced approach would work transitively.
	// An even more advanced approach would be context-aware and
	// determine nil errors based on inputs (e.g. io.WriteString to a
	// bytes.Buffer will always return nil, but an io.WriteString to
	// an os.File might not). Similarly, an os.File opened for reading
	// won't error on Close, but other files will.
	res := fn.Signature.Results()
	if res.Len() == 0 {
		return false
	}
	last := res.At(res.Len() - 1)
	if types.TypeString(last.Type(), nil) != "error" {
		return false
	}

	if fn.Blocks == nil {
		return false
	}
	for _, block := range fn.Blocks {
		if len(block.Instrs) == 0 {
			continue
		}
		ins := block.Instrs[len(block.Instrs)-1]
		ret, ok := ins.(*ssa.Return)
		if !ok {
			continue
		}
		v := ret.Results[len(ret.Results)-1]
		c, ok := v.(*ssa.Const)
		if !ok {
			return false
		}
		if !c.IsNil() {
			return false
		}
	}
	return true
}
