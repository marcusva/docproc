package processors

import (
	"encoding/json"
	"fmt"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/rules"
	"io/ioutil"
)

const (
	validatorName = "ContentValidator"
)

func init() {
	Register(validatorName, NewContentValidator)
}

// ContentValidator checks, if a queue.Message validates against a certain set
// of rules.
type ContentValidator struct {
	// rules contains the rules to be executed on processing a message.
	rules []rules.Rule
}

// Name returns the name to be used in configuration files.
func (validator *ContentValidator) Name() string {
	return validatorName
}

// Process processes the message, executing all rules of the ContentValidator
// against its contents. If a rule does not match or throws an error, the check
// will abort with an error.
func (validator *ContentValidator) Process(msg *queue.Message) error {
	for _, rule := range validator.rules {
		ok, err := rule.Test(msg.Content)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("message does not satisfy rule '%v'", rule)
		}
	}
	return nil
}

// NewContentValidator creates a new ContentValidator.
// The parameter map params must contain the following entries:
//
// * "rules": path to a rules file that contains a set of JSON-based rules
//
func NewContentValidator(params map[string]string) (queue.Processor, error) {
	rulefile, ok := params["rules"]
	if !ok {
		return nil, fmt.Errorf("parameter 'rules' missing")
	}
	data, err := ioutil.ReadFile(rulefile)
	if err != nil {
		return nil, err
	}
	var ruleset []rules.Rule
	if err := json.Unmarshal(data, &ruleset); err != nil {
		return nil, err
	}
	if err := rules.Validate(&ruleset); err != nil {
		return nil, err
	}
	return &ContentValidator{
		rules: ruleset,
	}, nil
}
