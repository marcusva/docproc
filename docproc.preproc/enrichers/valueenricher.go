package enrichers

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

var (
	exprSubst = regexp.MustCompile(`\$\{(.*?)\}`)
)

func init() {
	Register("ValueEnricher", NewValueEnricher)
}

// ValueRule represents a Rule, that carries a key-value pair to be added
// to a queue.Message, if the Rule evaluates to true.
type ValueRule struct {
	*rules.Rule
	TargetPath  string
	TargetValue interface{}
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
	return &ValueEnricher{rules: rules}, nil
}

// Name returns the name of the ValueEnricher to be used in configuration
// files.
func (e *ValueEnricher) Name() string {
	return "ValueEnricher"
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

func setValue(values map[string]interface{}, segs []string, newval interface{}) (map[string]interface{}, error) {
	if len(segs) == 1 {
		values[segs[0]] = newval
		return values, nil
	}
	tmp, ok := values[segs[0]]
	if !ok {
		tmp = make(map[string]interface{})
	}
	if _, ok := tmp.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("subpath '%s' is not a map", segs[0])
	}
	ret, err := setValue(tmp.(map[string]interface{}), segs[1:], newval)
	if err != nil {
		return nil, err
	}
	values[segs[0]] = ret
	return values, nil
}

func (e *ValueEnricher) Process(msg *queue.Message) error {
	for _, rule := range e.rules {
		log.Debugf("testing '%v %v %v' against %v", rule.Path, rule.Operator, rule.Value, msg.Content)
		ok, err := rule.Test(msg.Content)
		if err != nil {
			return err
		}
		if ok {
			segments := strings.Split(rule.TargetPath, ".")
			newval, err := resolveValue(msg, rule.TargetValue)
			if err != nil {
				return err
			}
			content, err := setValue(msg.Content, segments, newval)
			if err != nil {
				return fmt.Errorf("could not set value for path '%s': %v", rule.TargetPath, err)
			}
			msg.Content = content
		}
	}
	return nil
}
