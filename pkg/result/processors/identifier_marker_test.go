package processors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestIdentifierMarker(t *testing.T) {
	cases := []struct{ in, out string }{
		{"unknown field Address in struct literal", "unknown field `Address` in struct literal"},
		{
			"invalid operation: res (variable of type github.com/iotexproject/iotex-core/explorer/idl/explorer.GetBlkOrActResponse) has no field or method Address",
			"invalid operation: `res` (variable of type `github.com/iotexproject/iotex-core/explorer/idl/explorer.GetBlkOrActResponse`) has no field or method `Address`",
		},
		{
			"should use a simple channel send/receive instead of select with a single case",
			"should use a simple channel send/receive instead of `select` with a single case",
		},
		{"var testInputs is unused", "var `testInputs` is unused"},
		{"undeclared name: stateIDLabel", "undeclared name: `stateIDLabel`"},
		{
			"exported type Metrics should have comment or be unexported",
			"exported type `Metrics` should have comment or be unexported",
		},
		{
			`comment on exported function NewMetrics should be of the form "NewMetrics ..."`,
			"comment on exported function `NewMetrics` should be of the form `NewMetrics ...`",
		},
		{
			"cannot use addr (variable of type string) as github.com/iotexproject/iotex-core/pkg/keypair.PublicKey value in argument to action.FakeSeal",
			"cannot use addr (variable of type `string`) as `github.com/iotexproject/iotex-core/pkg/keypair.PublicKey` value in argument to `action.FakeSeal`",
		},
		{"other declaration of out", "other declaration of `out`"},
		{"should check returned error before deferring response.Close()", "should check returned error before deferring `response.Close()`"},
		{"should use time.Since instead of time.Now().Sub", "should use `time.Since` instead of `time.Now().Sub`"},
		{"TestFibZeroCount redeclared in this block", "`TestFibZeroCount` redeclared in this block"},
		{"should replace i += 1 with i++", "should replace `i += 1` with `i++`"},
		{"createEntry - result err is always nil", "`createEntry` - result `err` is always `nil`"},
		{
			"should omit comparison to bool constant, can be simplified to !projectIntegration.Model.Storage",
			"should omit comparison to bool constant, can be simplified to `!projectIntegration.Model.Storage`",
		},
		{
			"if block ends with a return statement, so drop this else and outdent its block",
			"`if` block ends with a `return` statement, so drop this `else` and outdent its block",
		},
		{
			"should write pupData := ms.m[pupID] instead of pupData, _ := ms.m[pupID]",
			"should write `pupData := ms.m[pupID]` instead of `pupData, _ := ms.m[pupID]`",
		},
		{"no value of type uint is less than 0", "no value of type `uint` is less than `0`"},
		{"redundant return statement", "redundant `return` statement"},
		{"struct field Id should be ID", "struct field `Id` should be `ID`"},
		{
			"don't use underscores in Go names; var Go_lint should be GoLint",
			"don't use underscores in Go names; var `Go_lint` should be `GoLint`",
		},
		{
			"G501: Blacklisted import crypto/md5: weak cryptographic primitive",
			"G501: Blacklisted import `crypto/md5`: weak cryptographic primitive",
		},
		{
			"S1017: should replace this if statement with an unconditional strings.TrimPrefix",
			"S1017: should replace this `if` statement with an unconditional `strings.TrimPrefix`",
		},
	}
	p := NewIdentifierMarker()

	for _, c := range cases {
		out, err := p.Process([]result.Issue{{Text: c.in}})
		require.NoError(t, err)
		assert.Equal(t, []result.Issue{{Text: c.out}}, out)
	}
}
