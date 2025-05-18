//golangcitest:args -Eembeddedstructfieldcheck
//golangcitest:config_path testdata/embeddedstructfieldcheck_mutex.yml
package testdata

import "sync"

type MutextEmbedded struct {
	sync.Mutex // want `sync.Mutex should not be embedded`
}

type MutextNotEmbedded struct {
	mu sync.Mutex
}

type PointerMutextEmbedded struct {
	*sync.Mutex // want `sync.Mutex should not be embedded`
}

type RWMutextEmbedded struct {
	sync.RWMutex // want `sync.RWMutex should not be embedded`
}

type RWMutextNotEmbedded struct {
	mu sync.RWMutex
}

type PointerRWMutextEmbedded struct {
	*sync.RWMutex // want `sync.RWMutex should not be embedded`
}
