package gounqvet

import (
	"testing"

	"github.com/golangci/golangci-lint/v2/pkg/config"
)

func TestGounqvetWithSettings(t *testing.T) {
	settings := &config.GounqvetSettings{
		CheckSQLBuilders:    true,
		IgnoredFunctions:    []string{"fmt.Printf"},
		AllowedPatterns:     []string{`SELECT \* FROM information_schema\..*`},
		IgnoredDirectories:  []string{"vendor"},
		IgnoredFilePatterns: []string{"*_test.go"},
	}
	
	linter := New(settings)
	
	if linter == nil {
		t.Fatal("Expected linter to be created")
	}
	
	if linter.Name() != "gounqvet" {
		t.Fatalf("Expected linter name 'gounqvet', got '%s'", linter.Name())
	}
}

func TestGounqvetNilSettings(t *testing.T) {
	linter := New(nil)
	
	if linter == nil {
		t.Fatal("Expected linter to be created with nil settings")
	}
	
	if linter.Name() != "gounqvet" {
		t.Fatalf("Expected linter name 'gounqvet', got '%s'", linter.Name())
	}
}