package gs

import (
	"github.com/viant/toolbox"
	"strings"
)

type predicate struct {
	Filter
}

func (p *predicate) Apply(source interface{}) bool {
	event, ok := source.(*Event)
	if !ok {
		return false
	}
	if len(p.Prefix) > 0 {
		hasPrefix := false
		for _, prefix := range p.Prefix {
			if strings.HasPrefix(event.Name, prefix) {
				hasPrefix = true
				break
			}
		}
		if !hasPrefix {
			return false
		}
	}
	if len(p.Suffix) > 0 {
		hasSuffix := false
		for _, suffix := range p.Suffix {
			if strings.HasSuffix(event.Name, suffix) {
				hasSuffix = true
				break
			}
		}
		if !hasSuffix {
			return false
		}
	}
	return true
}

//NewPredicate creates a new predicate
func NewPredicate(filter interface{}) (toolbox.Predicate, error) {
	result := &predicate{}
	err := toolbox.DefaultConverter.AssignConverted(&result.Filter, filter)
	return result, err
}
