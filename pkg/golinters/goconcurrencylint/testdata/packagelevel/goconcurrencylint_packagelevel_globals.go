package packagelevel

import "sync"

var sharedPackageMu sync.Mutex
var sharedPackageWG sync.WaitGroup
