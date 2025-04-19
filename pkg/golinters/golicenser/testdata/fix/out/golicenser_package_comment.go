// Copyright (c) 2025 golangci-lint <test@example.com>.
// This file is a part of golangci-lint.

// Package testdata is used for testing whether golicenser works correctly.
// golicenser should ignore this package comment and add a license header above it.
//
//golangcitest:args -Egolicenser
//golangcitest:expected_exitcode 0
//golangcitest:config_path testdata/golicenser-fix.yml
package testdata
