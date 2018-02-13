// Package rules provides a simple rule evaluation system for map-based data.
package rules

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	exprIndex = regexp.MustCompile(`(.*)\[(\d+)\]$`)
)

// asFloat converts a value to a float64. This applies to basic types, such as
// int, uint, float and string and their variations. Other types (e.g. struct,
// array, ...) will return an error.
func asFloat(v interface{}) (float64, error) {
	switch t := v.(type) {
	case float64:
		return v.(float64), nil
	case float32:
		return float64(t), nil
	case int:
		return float64(t), nil
	case int8:
		return float64(t), nil
	case int16:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case uint:
		return float64(t), nil
	case uint8:
		return float64(t), nil
	case uint16:
		return float64(t), nil
	case uint32:
		return float64(t), nil
	case uint64:
		return float64(t), nil
	case string:
		return strconv.ParseFloat(v.(string), 64)
	default:
		return 0.0, fmt.Errorf("unsupported value '%v' for numeric float", v)
	}
}

// cmpNumeric compares two float64 using the specified operator. If the operator
// is unknown, an error will be returned.
func cmpNumeric(lval, rval float64, op string) (bool, error) {
	switch op {
	case "=", "==", "eq", "equals":
		return lval == rval, nil
	case "<>", "!=", "neq", "not equals":
		return lval != rval, nil
	case ">", "gt", "greater than":
		return lval > rval, nil
	case "<", "lt", "less than":
		return lval < rval, nil
	case ">=", "gte", "greater than or equals":
		return lval >= rval, nil
	case "<=", "lte", "less than or equals":
		return lval <= rval, nil
	default:
		return false, fmt.Errorf("unsupported operator '%s' for numerics", op)
	}
}

// cmpString compares two strings using the specified operator. If the operator
// is unknown, an error will be returned.
func cmpString(lval, rval string, op string) (bool, error) {
	switch op {
	case "=", "==", "eq", "equals":
		return lval == rval, nil
	case "<>", "!=", "neq", "not equals":
		return lval != rval, nil
	case ">", "gt", "greater than":
		return lval > rval, nil
	case "<", "lt", "less than":
		return lval < rval, nil
	case ">=", "gte", "greater than or equals":
		return lval >= rval, nil
	case "<=", "lte", "less than or equals":
		return lval <= rval, nil
	case "contains":
		return strings.Contains(lval, rval), nil
	case "not contains":
		return !strings.Contains(lval, rval), nil
	case "in":
		return strings.Contains(rval, lval), nil
	case "not in":
		return !strings.Contains(rval, lval), nil
	default:
		return false, fmt.Errorf("unsupported operator '%s' for strings", op)
	}
}

func getArrayValue(subpath string, m map[string]interface{}) interface{} {
	arr := exprIndex.FindStringSubmatch(subpath)
	if len(arr) == 3 {
		idx, err := strconv.Atoi(arr[2])
		if err == nil {
			if marray, ok := m[arr[1]].([]interface{}); ok {
				if len(marray) > idx {
					return marray[idx]
				}
			}
		}
	}
	return nil
}

// ExtractValue retrieves the value for the passed path from the map. The path
// can contain dots ('.') to indicate access to nested maps as well as array
// indices ('[index]').
func ExtractValue(path string, m map[string]interface{}) interface{} {
	segments := strings.Split(path, ".")
	last, rest := segments[len(segments)-1], segments[:len(segments)-1]
	var ok bool
	smap := m

	for _, part := range rest {
		// Check for array access "foo.bar.baz[index]"
		if part[len(part)-1] == ']' {
			arrayval := getArrayValue(part, smap)
			if smap, ok = arrayval.(map[string]interface{}); ok {
				continue
			}
			return nil // Invalid array access
		}
		if smap, ok = smap[part].(map[string]interface{}); !ok {
			return nil
		}
	}
	if last[len(last)-1] == ']' {
		return getArrayValue(last, smap)
	}
	return smap[last]
}

// OperatorSupported checks, if the passed op is a valid operator for a Rule.
func OperatorSupported(op string) bool {
	switch op {
	case "=", "==", "eq", "equals":
		return true
	case "<>", "!=", "neq", "not equals":
		return true
	case ">", "gt", "greater than":
		return true
	case "<", "lt", "less than":
		return true
	case ">=", "gte", "greater than or equals":
		return true
	case "<=", "lte", "less than or equals":
		return true
	case "contains", "not contains":
		return true
	case "in", "not in":
		return true
	case "exists", "not exists":
		return true
	default:
		return false
	}
}

// Rule
type Rule struct {
	Name     string      `json:"name,omitempty"`
	Path     string      `json:"path"`
	Operator string      `json:"op"`
	Value    interface{} `json:"value"`
	SubRules []Rule      `json:"subrules,omitempty"`
}

// Test tests, if the the passed map matches the rule criteria.
func (r *Rule) Test(m map[string]interface{}) (bool, error) {
	for _, subrule := range r.SubRules {
		rv, e := subrule.Test(m)
		if rv == false || e != nil {
			return false, e
		}
	}

	left := ExtractValue(r.Path, m)

	if r.Operator == "exists" {
		return left != nil, nil
	}
	if r.Operator == "not exists" {
		return left == nil, nil
	}
	if left == nil {
		// Value does not exist
		return false, nil
	}

	right := r.Value
	switch right.(type) {
	case string:
		return cmpString(fmt.Sprint(left), right.(string), r.Operator)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		fval, e := asFloat(left)
		if e != nil {
			return false, e
		}
		rval, e := asFloat(right)
		if e != nil {
			return false, e
		}
		return cmpNumeric(fval, rval, r.Operator)
	default:
		return false, fmt.Errorf("unsupported value '%v'", left)
	}
}
