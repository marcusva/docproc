package processors

import (
	"encoding/json"
	"fmt"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/rules"
	"io/ioutil"
	"regexp"
	"strings"
)

const (
	valName = "ValueEnricher"
)

func init() {
	Register(valName, NewValueEnricher)
}

var (
	exprSubst = regexp.MustCompile(`\$\{(.*?)\}`)
)

// ValueRule represents a Rule, that carries a key-value pair to be added
// to a queue.Message, if the Rule evaluates to true.
type ValueRule struct {
	rules.Rule
	TargetPath  string      `json:"targetpath"`
	TargetValue interface{} `json:"targetvalue"`
}

// ValueEnricher applies multiple ValueRule items to a queue.Message.
type ValueEnricher struct {
	rules []ValueRule
}

// NewValueEnricher creates a new ValueEnricher
func NewValueEnricher(params map[string]string) (queue.Processor, error) {
	rulefile, ok := params["rules"]
	if !ok {
		return nil, fmt.Errorf("parameter 'rules' missing")
	}
	data, err := ioutil.ReadFile(rulefile)
	if err != nil {
		return nil, err
	}
	var rules []ValueRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}
	if err := validateRules(&rules); err != nil {
		return nil, err
	}
	return &ValueEnricher{rules: rules}, nil
}

// Name returns the name to be used in configuration files.
func (e *ValueEnricher) Name() string {
	return valName
}

func resolveValue(msg *queue.Message, val interface{}) (interface{}, error) {
	switch val.(type) {
	case string:
		arr := exprSubst.FindAllStringSubmatch(val.(string), -1)
		if len(arr) == 0 {
			return val, nil
		}
		ret := val.(string)
		for _, group := range arr {
			// group[0] refers to the outer part of the match, e.g. ${foo}
			// group[1] captures the inner name of ${foo}: foo
			value := rules.ExtractValue(group[1], msg.Content)
			if value == nil {
				return nil, fmt.Errorf("no value for path '%s' found", group[1])
			}
			ret = strings.Replace(ret, group[0], fmt.Sprint(value), -1)
		}
		return ret, nil
	default:
		return val, nil
	}
}

// Process processes the passed message and adds or rewrites values of the
// message's content based on the configured rules.
func (e *ValueEnricher) Process(msg *queue.Message) error {
	for _, rule := range e.rules {
		log.Debugf("testing '%v %v %v' against %v", rule.Path, rule.Operator, rule.Value, msg.Content)
		ok, err := rule.Test(msg.Content)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		newval, err := resolveValue(msg, rule.TargetValue)
		if err != nil {
			return err
		}
		if err := setValue(msg, rule.TargetPath, newval); err != nil {
			return fmt.Errorf("could not set value for path '%s': %v", rule.TargetPath, err)
		}
	}
	return nil
}

func setValue(msg *queue.Message, path string, newval interface{}) error {
	segments := strings.Split(path, ".")
	var parent map[string]interface{}
	last, rest := segments[len(segments)-1], segments[:len(segments)-1]

	parent = msg.Content
	for _, seg := range rest {
		sub, ok := parent[seg]
		if !ok {
			sub = make(map[string]interface{})
			parent[seg] = sub
		} else if _, ismap := sub.(map[string]interface{}); !ismap {
			return fmt.Errorf("path fragment '%s' is not a map", seg)
		}
		parent, _ = sub.(map[string]interface{})
	}
	parent[last] = newval
	return nil
}

func validateRules(ruleset *[]ValueRule) error {
	for _, rule := range *ruleset {
		if rule.Path == "" {
			return fmt.Errorf("empty source path in rule '%v'", rule)
		}
		if !rules.OperatorSupported(rule.Operator) {
			return fmt.Errorf("unsupported operator in rule '%v'", rule)
		}
		if rule.TargetPath == "" {
			return fmt.Errorf("empty target path in rule '%v'", rule)
		}
		if rule.TargetValue == "" {
			return fmt.Errorf("empty target value in rule '%v'", rule)
		}
	}
	return nil
}
