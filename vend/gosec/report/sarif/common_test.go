package sarif_test

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	"github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/securego/gosec/v2/report/sarif"
)

var (
	sarifSchemaOnce sync.Once
	sarifSchema     *jsonschema.Schema
	sarifSchemaErr  error
)

//go:embed testdata/sarif-schema-2.1.0.json
var sarifSchemaJSON []byte

func validateSarifSchema(report *sarif.Report) error {
	GinkgoHelper()
	sarifSchemaOnce.Do(func() {
		schema, err := jsonschema.UnmarshalJSON(bytes.NewReader(sarifSchemaJSON))
		if err != nil {
			sarifSchemaErr = fmt.Errorf("unmarshal local sarif schema: %w", err)
			return
		}

		compiler := jsonschema.NewCompiler()
		if err := compiler.AddResource(sarif.Schema, schema); err != nil {
			sarifSchemaErr = fmt.Errorf("compile sarif schema: %w", err)
			return
		}

		sarifSchema, sarifSchemaErr = compiler.Compile(sarif.Schema)
	})

	if sarifSchemaErr != nil {
		return sarifSchemaErr
	}

	// Marshal the report to JSON
	v, err := json.MarshalIndent(report, "", "\t")
	if err != nil {
		return err
	}

	// Unmarshal into any for schema validation
	data, err := jsonschema.UnmarshalJSON(bufio.NewReader(bytes.NewReader(v)))
	if err != nil {
		return err
	}

	return sarifSchema.Validate(data)
}
