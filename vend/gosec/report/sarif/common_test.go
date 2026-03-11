package sarif_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	"github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/securego/gosec/v2/report/sarif"
)

var (
	sarifSchemaOnce   sync.Once
	sarifSchema       *jsonschema.Schema
	sarifSchemaErr    error
	sarifSchemaClient = &http.Client{Timeout: 30 * time.Second}
)

func validateSarifSchema(report *sarif.Report) error {
	GinkgoHelper()
	sarifSchemaOnce.Do(func() {
		resp, err := sarifSchemaClient.Get(sarif.Schema)
		if err != nil {
			sarifSchemaErr = fmt.Errorf("fetch sarif schema: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			sarifSchemaErr = fmt.Errorf("fetch sarif schema: unexpected status %s", resp.Status)
			return
		}

		schema, err := jsonschema.UnmarshalJSON(resp.Body)
		if err != nil {
			sarifSchemaErr = fmt.Errorf("error unmarshaling schema: %w", err)
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
