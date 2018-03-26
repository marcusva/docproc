package rules

import (
	"bytes"
	"encoding/json"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

const (
	dmap = `
{
	"textVal":  "some text to test",
	"uintVal":  1234,
	"intVal":   -985,
	"floatVal": 1234.56,
	"nested1": {
		"nestedVal": "a nested value",
		"nested2": {
			"another": "Another nested value"
		},
		"array": [10, 20, 30, 40, 50, 60, 70]
	},
	"arrayMap": [
		{ "name": "map1", "value": "value1" },
		{ "name": "map2", "value": "value2" },
		{ "name": "map3", "value": "value3" },
		{ "name": "map4", "value": "value4" }
	]
}
`
)

var (
	validRules = []Rule{
		// contains
		{Path: "textVal", Operator: "contains", Value: "some"},
		// not contains
		{Path: "textVal", Operator: "not contains", Value: "apple pie"},
		// contains
		{Path: "textVal", Operator: "in", Value: "this is some text to test!"},
		// not contains
		{Path: "textVal", Operator: "not in", Value: "apple pie"},
		// equals
		{Path: "textVal", Operator: "eq", Value: "some text to test"},
		{Path: "uintVal", Operator: "=", Value: 1234},
		{Path: "intVal", Operator: "==", Value: -985},
		{Path: "floatVal", Operator: "equals", Value: 1234.56},
		// not equals
		{Path: "textVal", Operator: "neq", Value: "different"},
		{Path: "uintVal", Operator: "<>", Value: 654},
		{Path: "intVal", Operator: "!=", Value: 2222},
		{Path: "floatVal", Operator: "not equals", Value: 556},
		{Path: "floatVal", Operator: "neq", Value: "test"},
		// greater than
		{Path: "textVal", Operator: "gt", Value: "short"},
		{Path: "textVal", Operator: "gte", Value: "short"},
		{Path: "textVal", Operator: "gte", Value: "some text to test"},
		{Path: "uintVal", Operator: ">", Value: 4},
		{Path: "uintVal", Operator: ">=", Value: 4},
		{Path: "uintVal", Operator: ">=", Value: 1234},
		{Path: "intVal", Operator: "greater than", Value: -7474},
		{Path: "intVal", Operator: "greater than or equals", Value: -7474},
		{Path: "intVal", Operator: "greater than or equals", Value: -985},
		{Path: "floatVal", Operator: "gt", Value: 200},
		{Path: "floatVal", Operator: "gte", Value: 200},
		{Path: "floatVal", Operator: "gte", Value: 1234.56},
		// less than
		{Path: "textVal", Operator: "lt", Value: "very long text here!"},
		{Path: "textVal", Operator: "lte", Value: "very long text here!"},
		{Path: "textVal", Operator: "lte", Value: "some text to test"},
		{Path: "uintVal", Operator: "<", Value: 120000},
		{Path: "uintVal", Operator: "<=", Value: 120000},
		{Path: "uintVal", Operator: "<=", Value: 1234},
		{Path: "intVal", Operator: "less than", Value: 0},
		{Path: "intVal", Operator: "less than or equals", Value: 0},
		{Path: "intVal", Operator: "less than or equals", Value: -985},
		{Path: "floatVal", Operator: "lt", Value: 9929484},
		{Path: "floatVal", Operator: "lte", Value: 9929484},
		{Path: "floatVal", Operator: "lte", Value: 1234.56},
		// exists, not exists
		{Path: "textVal", Operator: "exists", Value: nil},
		{Path: "nonexisting", Operator: "not exists", Value: nil},
	}
	erroneousRules = []Rule{
		{Path: "textVal", Operator: "=", Value: 1234},
		{Path: "floatVal", Operator: "=", Value: true},
	}
	invalidComparators = []Rule{
		// Strings
		{Path: "textVal", Operator: "is the same as", Value: "some text to test"},
		{Path: "textVal", Operator: "", Value: nil},
		// Numerics
		{Path: "uintVal", Operator: "wants to be", Value: 1234},
		{Path: "uintVal", Operator: "", Value: 1234},
	}
	pathRules = []Rule{
		{Path: "textVal", Operator: "exists", Value: nil},
		{Path: "nested1.nestedVal", Operator: "exists", Value: nil},
		{Path: "nested1.invalid", Operator: "not exists", Value: nil},
		{Path: "nested1.nested2.another", Operator: "exists", Value: nil},
		{Path: "nested1.array[0]", Operator: "eq", Value: 10},
		{Path: "arrayMap[2].name", Operator: "eq", Value: "map3"},
	}
	validSubRules = []Rule{
		{Path: "uintVal", Operator: "=", Value: 1234, SubRules: []Rule{
			{Path: "intVal", Operator: "<", Value: 0},
		}},
	}
	invalidSubRules = []Rule{
		{Path: "uintVal", Operator: "=", Value: 1234, SubRules: []Rule{
			{Path: "nonexisting", Operator: ">", Value: 0},
		}},
	}
	invalidRules = []Rule{
		{Path: "", Operator: "=", Value: 1},
		{Path: "path", Operator: "", Value: 1},
		{Path: "path", Operator: ">", Value: 1, SubRules: []Rule{
			{Path: "", Operator: ">", Value: 0},
		}},
	}
)

func TestRules(t *testing.T) {
	var ct map[string]interface{}
	buf := bytes.NewBufferString(dmap).Bytes()
	assert.FailOnErr(t, json.Unmarshal(buf, &ct))

	for _, r := range validRules {
		ok, err := r.Test(ct)
		assert.FailOnErr(t, err)
		assert.Equal(t, ok, true)
		assert.FailOnErr(t, r.Validate())
	}
	assert.FailOnErr(t, Validate(&validRules))

	for _, r := range erroneousRules {
		ok, err := r.Test(ct)
		assert.Err(t, err, "%v", r)
		assert.Equal(t, ok, false)
		assert.FailOnErr(t, r.Validate())
	}

	for _, r := range invalidComparators {
		_, err := r.Test(ct)
		assert.Err(t, err)
		assert.Err(t, r.Validate())
	}
	assert.Err(t, Validate(&invalidComparators))

	for _, r := range validSubRules {
		ok, err := r.Test(ct)
		assert.FailOnErr(t, err)
		assert.Equal(t, ok, true)
		assert.FailOnErr(t, r.Validate())
	}
	assert.FailOnErr(t, Validate(&validSubRules))

	for _, r := range invalidSubRules {
		ok, err := r.Test(ct)
		assert.FailOnErr(t, err)
		assert.Equal(t, ok, false)
	}

	for _, r := range invalidRules {
		assert.Err(t, r.Validate(), "rule %v must not validate", r)
	}
	assert.Err(t, Validate(&invalidRules))
}

func TestPaths(t *testing.T) {
	var ct map[string]interface{}
	buf := bytes.NewBufferString(dmap).Bytes()
	assert.FailOnErr(t, json.Unmarshal(buf, &ct))

	for _, r := range pathRules {
		ok, err := r.Test(ct)
		assert.FailOnErr(t, err)
		assert.Equal(t, ok, true)
	}
}

func TestOperatorSupported(t *testing.T) {
	operators := []string{
		"=", "==", "eq", "equals",
		"!=", "<>", "neq", "not equals",
		">", "gt", "greater than",
		">=", "gte", "greater than or equals",
		"<", "lt", "less than",
		"<=", "lte", "less than or equals",
		"exists", "not exists",
		"contains", "not contains",
		"in", "not in",
	}
	for _, op := range operators {
		assert.FailIfNot(t, OperatorSupported(op), "operator '%s' failed unexpectedly", op)
	}
	assert.FailIf(t, OperatorSupported(""))
	assert.FailIf(t, OperatorSupported("something wrong"))
}
