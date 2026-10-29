package packagelevel

func badSplitPackageMutexCase() {
	sharedPackageMu.Lock()
}

func badSplitPackageWaitGroupCase() {
	sharedPackageWG.Add(1)
	sharedPackageWG.Wait()
}
